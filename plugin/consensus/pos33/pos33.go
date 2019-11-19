package pos33

import (
	"bytes"
	"errors"
	"math/rand"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/ed25519"
	"github.com/33cn/chain33/common/merkle"
	"github.com/33cn/chain33/queue"
	drivers "github.com/33cn/chain33/system/consensus"
	ced25519 "github.com/33cn/chain33/system/crypto/ed25519"
	ct "github.com/33cn/chain33/system/dapp/coins/types"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/pos33/types"
	// "github.com/33cn/chain33/util"
)

const rootSeed = "YCC-ROOT"
const pos33MinFee = 1e7

// RootPrivKey is the root private key for ycc
var RootPrivKey crypto.PrivKey

// RootAddr is the root account address for ycc
var RootAddr string

var myCrypto crypto.Crypto

func genKeyFromSeed(seed []byte) (crypto.PrivKey, error) {
	_, priv, err := ed25519.GenerateKey(bytes.NewReader(crypto.Sha256(seed)))
	if err != nil {
		return nil, err
	}
	return ced25519.PrivKeyEd25519(*priv), nil
}

func init() {
	drivers.Reg("pos33", New)
	cpt, err := crypto.New("ed25519")
	if err != nil {
		panic(err)
	}
	myCrypto = cpt
	priv, err := genKeyFromSeed([]byte(rootSeed))
	if err != nil {
		panic(err)
	}
	RootPrivKey = priv
	RootAddr = address.PubKeyToAddress(RootPrivKey.PubKey().Bytes()).String()
	rand.Seed(time.Now().Unix())
}

// Client is the pos33 consensus client
type Client struct {
	*drivers.BaseClient
	conf *subConfig
	n    *node
	priv crypto.PrivKey
}

type subConfig struct {
	Pos33SecretSeed    string `json:"Pos33SecretSeed,omitempty"`
	Pos33ListenAddr    string `json:"Pos33ListenAddr,omitempty"`
	Pos33AdvertiseAddr string `json:"Pos33AdvertiseAddr,omitempty"`
	Pos33BootPeerAddr  string `json:"Pos33BootPeerAddr,omitempty"`
	Pos33MaxTxs        int64  `json:"Pos33MaxTxs,omitempty"`
	Pos33BlockTime     int64  `json:"Pos33BlockTime,omitempty"`
	Pos33BlockTimeout  int64  `json:"Pos33BlockTimeout,omitempty"`
	Pos33MinFee        int64  `json:"Pos33MinFee,omitempty"`
}

// New create pos33 consensus client
func New(cfg *types.Consensus, sub []byte) queue.Module {
	c := drivers.NewBaseClient(cfg)
	var subcfg subConfig
	if sub != nil {
		types.MustDecode(sub, &subcfg)
	}

	n := newNode(&subcfg)
	client := &Client{BaseClient: c, n: n, conf: &subcfg}
	client.n.Client = client
	client.Cfg.Genesis = RootAddr
	c.SetChild(client)
	return client
}

// Close is close the client
func (client *Client) Close() {}

// ProcEvent do nothing?
func (client *Client) ProcEvent(msg *queue.Message) bool {
	return false
}

func (client *Client) newBlock(lastBlock *types.Block, txs []*types.Transaction, height int64) (*types.Block, error) {
	if lastBlock.Height+1 != height {
		plog.Error("newBlock height error", "lastHeight", lastBlock.Height, "height", height)
		return nil, errors.New("the last block too low")
	}

	ch := make(chan []*Tx, 1)
	go func() { ch <- client.RequestTx(int(client.conf.Pos33MaxTxs), nil) }()
	select {
	case <-time.After(time.Millisecond * 300):
	case ts := <-ch:
		txs = append(txs, ts...)
	}

	bt := time.Now().Unix()
	return &types.Block{
		ParentHash: lastBlock.Hash(),
		Height:     lastBlock.Height + 1,
		Txs:        txs,
		TxHash:     merkle.CalcMerkleRoot(txs),
		BlockTime:  bt,
	}, nil
}

// CheckBlock check block callback
func (client *Client) CheckBlock(parent *types.Block, current *types.BlockDetail) error {
	return client.n.checkBlock(current.Block, parent)
}

func (client *Client) allWeight(height int64) int {
	k := []byte(pt.KeyPos33AllWeight)
	v, err := client.Get(k)
	if err != nil {
		plog.Error(err.Error())
		return 0
	}
	var allw pt.Pos33AllWeight
	err = types.Decode(v, &allw)
	if err != nil {
		plog.Error(err.Error())
		return 0
	}
	w := allw.AllWeight
	for h, nw := range allw.NewWeight {
		if h+pt.Pos33SortitionSize*2-h%pt.Pos33SortitionSize < height {
			w += nw
		}
	}

	return int(w)
}

func (client *Client) getWeight(addr string, height int64) int {
	v, err := client.Get([]byte(pt.KeyPos33WeightPrefix + addr))
	if err != nil {
		plog.Info("getWeight error", "error", err.Error())
		return 0
	}
	var w pt.Pos33Weight
	err = types.Decode(v, &w)
	if err != nil {
		plog.Error("getWeight error", "error", err.Error())
		return 0
	}
	allw := int64(0)
	for h, nw := range w.Weights {
		if h < pt.Pos33SortitionSize || h+pt.Pos33SortitionSize*2-h%pt.Pos33SortitionSize <= height {
			allw += nw
		}
	}
	return int(allw)
}

// AddBlock notice driver a new block incoming
func (client *Client) AddBlock(b *types.Block) error {
	client.n.addBlock(b)
	return nil
}

// CreateBlock will start run
func (client *Client) CreateBlock() {
	for {
		if !client.IsMining() || !(client.IsCaughtUp() || client.Cfg.ForceMining) {
			plog.Info("createblock.ismining is disable or client is caughtup is false")
			time.Sleep(time.Second)
			continue
		}
		if client.getWeight(client.n.addr, client.GetCurrentHeight()) == 0 {
			plog.Info("if do consensus, must deposit 1,000,000 YCC")
			time.Sleep(time.Second)
			continue
		}
		break
	}
	client.n.runLoop()
}

// CreateGenesisTx used generate the first txs
func (client *Client) CreateGenesisTx() (ret []*types.Transaction) {
	// the 1st tx for issue 10,000,000,000 YCC
	act0 := &ct.CoinsAction_Genesis{Genesis: &types.AssetsGenesis{Amount: types.Coin * 100 * 1e8}}
	tx := &types.Transaction{
		Execer:  []byte("coins"),
		To:      RootAddr,
		Payload: types.Encode(&ct.CoinsAction{Value: act0, Ty: ct.CoinsActionGenesis}),
	}
	ret = append(ret, tx)

	// the 2th tx for the genesis accout frozon margin,
	// so the second block must created by the genesis accout.
	tx = &types.Transaction{}
	tx.Execer = []byte("pos33")
	tx.To = address.GetExecAddress("pos33").String()
	dact := &pt.Pos33DepositAction{W: 1000}
	tx.Payload = types.Encode(&pt.Pos33Action{Value: &pt.Pos33Action_Deposit{Deposit: dact}, Ty: int32(pt.Pos33ActionDeposit)})
	tx.Sign(types.ED25519, RootPrivKey)
	ret = append(ret, tx)
	return
}

// write block to chain
func (client *Client) setBlock(b *types.Block) error {
	if b == nil {
		plog.Crit("block is nil")
		return nil
	}

	plog.Info("setBlock", "height", b.Height, "txCount", len(b.Txs))
	lastBlock, err := client.RequestBlock(b.Height - 1)
	if err != nil {
		return err
	}
	err = client.WriteBlock(lastBlock.StateHash, b)
	if err != nil {
		return err
	}
	return nil
}

func getBlockReword(b *types.Block) (*pt.Pos33MinerAction, error) {
	if b == nil {
		return nil, errors.New("b is nil")
	}
	tx := b.Txs[0]
	var pact pt.Pos33Action
	err := types.Decode(tx.Payload, &pact)
	if err != nil {
		return nil, err
	}
	act := pact.GetMiner()
	return act, nil
}

// Get used search block store db
func (client *Client) Get(key []byte) ([]byte, error) {
	query := &types.LocalDBGet{Keys: [][]byte{key}}
	msg := client.GetQueueClient().NewMessage("blockchain", types.EventLocalGet, query)
	client.GetQueueClient().Send(msg, true)
	resp, err := client.GetQueueClient().Wait(msg)

	if err != nil {
		plog.Error(err.Error()) //no happen for ever
		return nil, err
	}
	value := resp.GetData().(*types.LocalReplyValue).Values[0]
	if value == nil {
		return nil, types.ErrNotFound
	}
	return value, nil
}

func (client *Client) sendTx(tx *types.Transaction) error {
	qcli := client.GetQueueClient()
	if qcli == nil {
		panic("client not bind message queue.")
	}
	msg := qcli.NewMessage("mempool", types.EventTx, tx)
	err := qcli.Send(msg, true)
	if err != nil {
		return err
	}
	resp, err := qcli.Wait(msg)
	if err != nil {
		return err
	}
	r := resp.GetData().(*types.Reply)
	if r.IsOk {
		return nil
	}
	plog.Info("sendTx error:", "error", string(r.Msg))
	return errors.New(string(r.Msg))
}
