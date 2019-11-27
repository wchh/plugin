package pos33

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/difficulty"
	"github.com/33cn/chain33/common/log/log15"
	driver "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pt "github.com/33cn/plugin/plugin/dapp/pos33/types"
	"github.com/golang/protobuf/proto"
)

var plog = log15.New("module", "pos33")

type node struct {
	*Client
	gss *gossip

	// I'm candidate proposer in these blocks
	ips map[int64]*pt.Pos33SortitionMsg
	// I'm candidate verifer in these blocks
	ivs map[int64][]*pt.Pos33SortitionMsg
	// receive candidate proposers
	cps map[int64]map[int]map[string]*pt.Pos33SortitionMsg
	// receive candidate verifers
	cvs map[int64]map[int]map[string][]*pt.Pos33VoteMsg
	// receive candidate blocks
	cbs map[int64]map[string]*types.Block

	voteOkHeight int64
	lastBlock    *types.Block
	bch          chan *types.Block
	etm          *time.Timer
}

// New create pos33 consensus client
func newNode(conf *subConfig) *node {
	n := &node{
		ips: make(map[int64]*pt.Pos33SortitionMsg),
		ivs: make(map[int64][]*pt.Pos33SortitionMsg),
		cps: make(map[int64]map[int]map[string]*pt.Pos33SortitionMsg),
		cvs: make(map[int64]map[int]map[string][]*pt.Pos33VoteMsg),
		cbs: make(map[int64]map[string]*types.Block),
		bch: make(chan *types.Block, 1),
	}

	plog.Info("@@@@@@@ node start:", "addr", addr, "conf", conf)
	return n
}

func unmarshal(b []byte) (*pt.Pos33Msg, error) {
	var pm pt.Pos33Msg
	err := proto.Unmarshal(b, &pm)
	if err != nil {
		return nil, err
	}
	return &pm, nil
}

func (n *node) getNotNullBlock(height int64) (*types.Block, error) {
	for i := height; i >= height-10; i-- {
		b, err := n.RequestBlock(i)
		if err != nil {
			return nil, err
		}
		if len(b.Txs) > 0 {
			return b, nil
		}
	}
	return nil, nil
}

func (n *node) genMinerTx(height int64, round int, strHash string) (*types.Transaction, int, error) {
	var vs []*pt.Pos33VoteMsg
	mv := n.cvs[height][round]
	for hash, vs := range mv {
		if hash == strHash {
			vs = append(vs, vs...)
		}
	}

	plog.Info("genRewordTx", "height", height, "vsw", len(vs))

	sort := n.ips[height+1]

	data, err := proto.Marshal(&pt.TicketAction{
		Value: &pt.TicketAction_Pminer{
			Pminer: &pt.Pos33Miner{
				Votes: vs,
				Sort:  sort,
			},
		},
		Ty: pt.TicketActionMiner,
	})

	if err != nil {
		panic(err)
	}

	tx := &types.Transaction{
		Execer:  []byte(pt.TicketX),
		To:      driver.ExecAddress(pt.TicketX),
		Payload: data,
		Nonce:   rand.Int63(),
		Expire:  time.Now().Unix() + 10,
	}
	priv := n.privmap[n.ticketsMap[sort.Input.TicketID].MinerAddress]
	tx.Sign(types.SECP256K1, priv)
	return tx, len(vs), nil
}

func (n *node) diff(w int) uint32 {
	// return types.GetP(0).PowLimitBits
	cfg := n.GetAPI().GetConfig()
	oldDiff := difficulty.CompactToBig(cfg.GetP(0).PowLimitBits)
	newDiff := new(big.Int).Add(oldDiff, big.NewInt(int64(w+1)))
	return difficulty.BigToCompact(newDiff)
}

func (n *node) makeBlock(round int) (*types.Block, error) {
	vh := n.lastBlock.Height
	height := vh + 1
	cfg := n.GetAPI().GetConfig()
	strHash := string(n.lastBlock.Hash(cfg))
	tx, w, err := n.genMinerTx(vh, round, strHash)
	if err != nil {
		plog.Error("genRewordTx error", "err", err.Error(), "height", height)
		return nil, err
	}

	nb, err := n.newBlock(n.lastBlock, []*Tx{tx}, height)
	if err != nil {
		plog.Error("makeBlock error", "height", height, "error", err.Error())
		return nil, err
	}

	nb.Difficulty = n.diff(w)

	plog.Info("@@@@@@@ I make a block: ", "height", height, "hash", hexs(nb.Hash(n.GetAPI().GetConfig())), "txCount", len(nb.Txs), "diff", nb.Difficulty)
	return nb, nil
}

func (n *node) addBlock(b *types.Block) {
	if !n.IsCaughtUp() {
		return
	}
	plog.Info("node.addBlock", "height", b.Height, "hash", hexs(b.Hash(n.GetAPI().GetConfig())))

	if n.lastBlock != nil && b.Height <= n.lastBlock.Height {
		plog.Info("addBlock nil", "height", b.Height)
		return
	}

	select {
	case n.bch <- b:
	case <-n.bch:
		n.bch <- b
	}
}

func (n *node) clear(height int64) {
	// clear the caches
	for h := range n.cbs {
		if h+10 <= height {
			delete(n.cbs, h)
		}
	}
	for h := range n.cvs {
		if h+10 <= height {
			delete(n.cvs, h)
		}
	}
	for h := range n.cps {
		if h+10 <= height {
			delete(n.cps, h)
		}
	}
	for h := range n.ips {
		if h+10 <= height {
			delete(n.ips, h)
		}
	}
	for h := range n.ivs {
		if h+10 <= height {
			delete(n.ivs, h)
		}
	}
}

func addr(sig *types.Signature) string {
	if sig == nil {
		return ""
	}
	return address.PubKeyToAddress(sig.Pubkey).String()
}

func (n *node) handleBlock(b *types.Block) {
	strHash := string(b.Hash(n.GetAPI().GetConfig()))
	plog.Info("node.handleBlock", "height", b.Height, "bp", addr(b.Txs[0].Signature), "hash", hexs(b.Hash(n.GetAPI().GetConfig())))

	if n.cbs[b.Height] == nil {
		n.cbs[b.Height] = make(map[string]*types.Block)
	}

	n.cbs[b.Height][strHash] = b
}

func (n *node) checkBlock(b, pb *types.Block) error {
	plog.Info("node.checkBlock", "height", b.Height)
	if b.Height < 2 {
		return nil
	}
	if !n.IsCaughtUp() {
		return nil
	}
	if len(b.Txs) == 0 {
		return errors.New("nil block error")
	}

	act, err := getMiner(b)
	if err != nil {
		return err
	}

	height := b.Height - 1
	seed := act.GetSort().GetHash()

	round := -1
	// check votes
	for _, v := range act.Votes {
		if !v.Verify() {
			return errors.New("vote signature error")
		}
		m := v.Sort
		// check height, height of vote is the pre-block height
		if m.Input.Height != height {
			return errors.New("height error")
		}

		// check round, all round of vote must be same
		if round == -1 {
			round = int(m.Input.Round)
		} else if round != int(m.Input.Round) {
			return errors.New("round error")
		}

		// check vrf
		err = n.verifySort(height, 1, seed, m)
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *node) getMinerSeed(height int64) ([]byte, error) {
	startHeight := height - height%pt.Pos33SortitionSize
	b, err := n.RequestBlock(startHeight)
	if err != nil {
		plog.Info("should't go here. do nothing")
		return nil, err
	}
	seed := zeroHash[:]
	act, err := getMiner(b)
	if err != nil {
		plog.Error("getBlockReword err", "error", err, "height", b.Height)
		return nil, err
	}
	seed = act.Sort.GetHash()
	return seed, nil
}

var zeroHash [32]byte

func (n *node) reSortition(height int64, round int) {
	seed, err := n.getMinerSeed(height)
	if err != nil {
		return
	}
	const staps = 2
	for s := 0; s < staps; s++ {
		sms := n.sort(seed, height, round, s)
		if sms == nil {
			plog.Info("node.sortition nil")
			continue
		}
		plog.Info("node.sortition", "height", height, "weight", len(sms))
		if s == 0 {
			n.ips[height] = sms[0]
		} else {
			n.ivs[height] = sms
		}
	}
}

func (n *node) sortition(b *types.Block, round int) {
	seed := zeroHash[:]
	startHeight := int64(0)
	if b != nil {
		act, err := getMiner(b)
		if err != nil {
			plog.Error("getBlockReword err", "error", err, "height", b.Height)
		} else {
			seed = act.Sort.Hash
			startHeight = b.Height
		}
	}
	startHeight++
	const steps = 2
	for s := 0; s < steps; s++ {
		for i := 0; i < pt.Pos33SortitionSize; i++ {
			height := startHeight + int64(i)
			sms := n.sort(seed, height, round, s)
			if sms == nil {
				plog.Info("node.sortition nil")
				continue
			}
			plog.Info("node.sortition", "height", height, "weight", len(sms))
			if s == 0 {
				n.ips[height] = sms[0]
			} else {
				n.ivs[height] = sms
			}
		}
	}
}

func (n *node) handleVote(vm *pt.Pos33VoteMsg) error {
	if !vm.Verify() {
		return errors.New("votemsg verify false")
	}
	height := vm.Sort.Input.Height
	if vm == nil || height <= n.lastBlock.Height {
		return fmt.Errorf("votemsg too late height = %d", height)
	}

	a := addr(vm.Sig)
	plog.Info("handleVoteMsg", "height", height, "voter", a, "bhash", hexs(vm.BlockHash))

	seed, err := n.getMinerSeed(height)
	if err != nil {
		return err
	}
	err = n.verifySort(height, 2, seed, vm.Sort)
	if err != nil {
		return err
	}

	if n.cvs[height] == nil {
		n.cvs[height] = make(map[int]map[string][]*pt.Pos33VoteMsg)
	}

	round := int(vm.Sort.Input.Round)
	if n.cvs[height][round] == nil {
		mp := make(map[string][]*pt.Pos33VoteMsg)
		n.cvs[height][round] = mp
	}

	vs := n.cvs[height][round][string(vm.BlockHash)]
	for _, v := range vs {
		if v.Equal(vm) {
			return errors.New("votemsg repeated")
		}
	}
	vs = append(vs, vm)
	n.cvs[height][round][string(vm.BlockHash)] = vs

	if n.voteOkHeight >= height {
		return nil
	}
	strHash := string(vm.BlockHash)

	if len(n.cvs[height][round][strHash])*3 > pt.Pos33VoterSize*2 {
		b, ok := n.cbs[height][strHash]
		if !ok {
			return nil
		}
		plog.Info("@@@ set block 2f+1 @@@", "height", height, "bp", addr(b.Txs[0].Signature), "hash", hexs(vm.BlockHash))
		n.voteOkHeight = b.Height

		n.setBlock(b)
	}
	return nil
}

func (n *node) handleVotesMsg(vms *pt.Pos33Votes) {
	if n.lastBlock == nil {
		return
	}
	if len(vms.Vs) == 0 {
		plog.Error("votemsg sortition is 0")
		return
	}
	for _, vm := range vms.Vs {
		err := n.handleVote(vm)
		if err != nil {
			plog.Error("handleVote error", "err", err)
			return
		}
	}
}

func (n *node) handleSortitionMsg(m *pt.Pos33SortitionMsg) {
	height := m.Input.Height
	plog.Info("handleElectMsg", "height", height)
	seed, err := n.getMinerSeed(height)
	n.verifySort(height, 0, seed, m)
	if err != nil {
		plog.Info("check rand error:", "error", err.Error())
		return
	}
	if n.cps[height] == nil {
		n.cps[height] = make(map[int]map[string]*pt.Pos33SortitionMsg)
	}
	round := int(m.Input.Round)
	if n.cps[height][round] == nil {
		n.cps[height][round] = make(map[string]*pt.Pos33SortitionMsg)
	}
	n.cps[height][round][m.Input.TicketID] = m
}

func (n *node) handlePos33Msg(pm *pt.Pos33Msg) bool {
	if pm == nil {
		return false
	}
	switch pm.Ty {
	case pt.Pos33Msg_B:
		var m pt.Pos33BlockMsg
		err := types.Decode(pm.Data, &m)
		if err != nil {
			plog.Error(err.Error())
			return false
		}
		n.handleBlock(m.B)
	case pt.Pos33Msg_S:
		var m pt.Pos33SortitionMsg
		err := types.Decode(pm.Data, &m)
		if err != nil {
			plog.Error(err.Error())
			return false
		}
		n.handleSortitionMsg(&m)
	case pt.Pos33Msg_V:
		var vt pt.Pos33Votes
		err := types.Decode(pm.Data, &vt)
		if err != nil {
			plog.Error(err.Error())
			return false
		}
		n.handleVotesMsg(&vt)
	default:
		panic("not support this message type")
	}

	return true
}

// doGossipMsg multi-goroutine verify pos33 message
func (n *node) doGossipMsg() chan *pt.Pos33Msg {
	num := 4
	ch := make(chan *pt.Pos33Msg, num*16)
	for i := 0; i < num; i++ {
		go func() {
			for {
				pm, err := unmarshal(<-n.gss.C)
				if err != nil {
					plog.Error(err.Error())
					continue
				}
				ch <- pm
			}
		}()
	}
	return ch
}

func reseTm(tm *time.Timer, d time.Duration) {
	if !tm.Stop() {
		select {
		case <-tm.C:
		default:
		}
	}
	tm.Reset(d)
}

func (n *node) firstSortition(firtstBlock *types.Block) {
	n.sortition(nil, 0)
}

func (n *node) runLoop() {
	lb, err := n.RequestLastBlock()
	if err != nil {
		plog.Error(err.Error())
		return
	}

	plog.Info("pos33 node runing.......", "last block height", lb.Height)

	n.gss = newGossip(n.nodeID(), n.conf.ListenAddr, n.conf.AdvertiseAddr, n.conf.BootPeerAddr)
	go n.gss.runBroadcast()
	msgch := n.doGossipMsg()

	n.etm = time.NewTimer(time.Hour)
	ch := make(chan int64, 1)

	if lb.Height == 0 {
		n.firstSortition(lb)
		time.AfterFunc(time.Second, func() { n.addBlock(lb) })
	}
	round := 0

	for {
		if !n.miningOK() {
			time.Sleep(time.Second)
			continue
		}
		select {
		case msg := <-msgch:
			n.handlePos33Msg(msg)
		case height := <-ch:
			if height == n.lastBlock.Height+1 {
				round++
				plog.Info("vote timeout: ", "height", height, "round", round)
				n.reSortition(height, round)
				n.elect(height, round)
				/*
					time.AfterFunc(time.Second*1, func() {
						ch <- height
					})
				*/
			}
		case <-n.etm.C:
			height := n.lastBlock.Height + 1
			plog.Info("elect timeout: ", "height", height)
			n.vote(height, round)

			to := (round + 1) * 2
			if to > 16 {
				to = 16
			}

			time.AfterFunc(time.Second*time.Duration(to), func() {
				ch <- height
			})
		case b := <-n.bch: // new block add to chain
			if b.Height < n.voteOkHeight {
				break
			}
			if b.Height%pt.Pos33SortitionSize == 0 {
				n.sortition(b, 0)
				n.flushTicket()
			}
			round = 0
			n.lastBlock = b
			plog.Info("set last block", "height", b.Height)
			n.elect(b.Height+1, 0)
			n.clear(b.Height - 1)
		}
	}
}

func hexs(b []byte) string {
	s := hex.EncodeToString(b)
	if len(s) <= 16 {
		return s
	}
	return s[:16]
}

func (n *node) vote(height int64, round int) {
	ivs, ok := n.ivs[height]
	if !ok {
		plog.Info("I'm not verifer", "height", height)
		return
	}

	var pss []*pt.Pos33SortitionMsg
	for _, s := range n.cps[height][round] {
		pss = append(pss, s)
	}
	if pss == nil {
		plog.Info("no block maker", "height", height)
		return
	}
	plog.Info("vote len(pss)", "height", height, "len(pss)", len(pss))

	sort.Sort(pt.Sorts(pss))
	bp := address.PubKeyToAddress(pss[0].Pubkey).String()

	plog.Info("vote bp", "height", height, "bp", bp, "len(cbs)", len(n.cbs[height]))
	var vb *types.Block
	for _, b := range n.cbs[height] {
		if addr(b.Txs[0].Signature) == bp {
			vb = b
			break
		}
	}
	if vb == nil {
		plog.Info("NO block vote out")
		return
	}
	var vs []*pt.Pos33VoteMsg
	for _, s := range ivs {
		v := &pt.Pos33VoteMsg{Sort: s, BlockHash: vb.Hash(n.GetAPI().GetConfig())}
		t := n.ticketsMap[s.Input.TicketID]
		priv := n.privmap[t.MinerAddress]
		v.Sign(priv)
		vs = append(vs, v)
	}
	v := &pt.Pos33Votes{Vs: vs}
	n.handleVotesMsg(v)
	n.gss.gossip(marshalVoteMsg(v))
}

func (n *node) elect(height int64, round int) {
	n.etm = time.NewTimer(time.Millisecond * 1000)
	pm, ok := n.ips[height]
	if !ok {
		plog.Info("elect: I'm not Proposer", "height", height, "round", round)
		return
	}
	nb, err := n.makeBlock(round)
	if err != nil {
		plog.Error(err.Error(), "height", height)
		return
	}
	n.handleSortitionMsg(pm)
	n.handleBlock(nb)
	n.gss.gossip(marshalElectMsg(pm))
	n.gss.gossip(marshalBlockMsg(&pt.Pos33BlockMsg{B: nb}))
}

func marshalBlockMsg(m *pt.Pos33BlockMsg) []byte {
	pm := &pt.Pos33Msg{
		Data: types.Encode(m),
		Ty:   pt.Pos33Msg_B,
	}
	return types.Encode(pm)
}

func marshalElectMsg(m *pt.Pos33SortitionMsg) []byte {
	pm := &pt.Pos33Msg{
		Data: types.Encode(m),
		Ty:   pt.Pos33Msg_S,
	}
	return types.Encode(pm)
}

func marshalVoteMsg(v *pt.Pos33Votes) []byte {
	pm := &pt.Pos33Msg{
		Data: types.Encode(v),
		Ty:   pt.Pos33Msg_V,
	}
	return types.Encode(pm)
}
