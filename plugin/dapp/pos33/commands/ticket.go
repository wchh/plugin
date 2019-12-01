// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package commands

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	ty "github.com/33cn/plugin/plugin/dapp/pos33/types"
	"github.com/spf13/cobra"
)

// Pos33TicketCmd ticket command type
func Pos33TicketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pos33.ticket",
		Short: "Pos33Ticket management",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		BindMinerCmd(),
		CountPos33TicketCmd(),
		ClosePos33TicketCmd(),
		GetColdAddrByMinerCmd(),
	)

	return cmd
}

// BindMinerCmd bind miner
func BindMinerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind_miner",
		Short: "Bind private key to miner address",
		Run:   bindMiner,
	}
	addBindMinerFlags(cmd)
	return cmd
}

func addBindMinerFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("bind_addr", "b", "", "miner address")
	cmd.MarkFlagRequired("bind_addr")

	cmd.Flags().StringP("origin_addr", "o", "", "origin address")
	cmd.MarkFlagRequired("origin_addr")
}

func bindMiner(cmd *cobra.Command, args []string) {
	title, _ := cmd.Flags().GetString("title")
	cfg := types.GetCliSysParam(title)

	bindAddr, _ := cmd.Flags().GetString("bind_addr")
	originAddr, _ := cmd.Flags().GetString("origin_addr")
	//c, _ := crypto.New(types.GetSignName(wallet.SignType))
	//a, _ := common.FromHex(key)
	//privKey, _ := c.PrivKeyFromBytes(a)
	//originAddr := account.PubKeyToAddress(privKey.PubKey().Bytes()).String()
	ta := &ty.Pos33TicketAction{}
	tBind := &ty.Pos33TicketBind{
		MinerAddress:  bindAddr,
		ReturnAddress: originAddr,
	}
	ta.Value = &ty.Pos33TicketAction_Tbind{Tbind: tBind}
	ta.Ty = ty.Pos33TicketActionBind

	tx, err := types.CreateFormatTx(cfg, ty.Pos33TicketX, types.Encode(ta))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	txHex := types.Encode(tx)
	fmt.Println(hex.EncodeToString(txHex))
}

// CountPos33TicketCmd get ticket count
func CountPos33TicketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Get ticket count",
		Run:   countPos33Ticket,
	}
	return cmd
}

func countPos33Ticket(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	var res int64
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "ticket.GetPos33TicketCount", nil, &res)
	ctx.Run()
}

// ClosePos33TicketCmd close all accessible tickets
func ClosePos33TicketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "close",
		Short: "Close tickets",
		Run:   closePos33Ticket,
	}
	addCloseBindAddr(cmd)
	return cmd
}

func addCloseBindAddr(cmd *cobra.Command) {
	cmd.Flags().StringP("miner_addr", "m", "", "miner address (optional)")
}

func closePos33Ticket(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	bindAddr, _ := cmd.Flags().GetString("miner_addr")
	status, err := getWalletStatus(rpcLaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	isAutoMining := status.(rpctypes.WalletStatus).IsAutoMining
	if isAutoMining {
		fmt.Fprintln(os.Stderr, types.ErrMinerNotClosed)
		return
	}

	tClose := &ty.Pos33TicketClose{
		MinerAddress: bindAddr,
	}

	var res rpctypes.ReplyHashes
	rpc, err := jsonclient.NewJSONClient(rpcLaddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	err = rpc.Call("pos33.ClosePos33Tickets", tClose, &res)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if len(res.Hashes) == 0 {
		fmt.Println("no ticket to be close")
		return
	}

	data, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	fmt.Println(string(data))
}

func getWalletStatus(rpcAddr string) (interface{}, error) {
	rpc, err := jsonclient.NewJSONClient(rpcAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	var res rpctypes.WalletStatus
	err = rpc.Call("Chain33.GetWalletStatus", nil, &res)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	return res, nil
}

// GetColdAddrByMinerCmd get cold address by miner
func GetColdAddrByMinerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cold",
		Short: "Get cold wallet address of miner",
		Run:   coldAddressOfMiner,
	}
	addColdAddressOfMinerFlags(cmd)
	return cmd
}

func addColdAddressOfMinerFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("miner", "m", "", "miner address")
	cmd.MarkFlagRequired("miner")
}

func coldAddressOfMiner(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	addr, _ := cmd.Flags().GetString("miner")
	reqaddr := &types.ReqString{
		Data: addr,
	}
	var params rpctypes.Query4Jrpc
	params.Execer = "pos33"
	params.FuncName = "MinerSourceList"
	params.Payload = types.MustPBToJSON(reqaddr)

	var res types.ReplyStrings
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, &res)
	ctx.Run()
}
