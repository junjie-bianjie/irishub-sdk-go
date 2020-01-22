package client

import (
	"errors"
	"fmt"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"strconv"

	"github.com/irisnet/irishub-sdk-go/modules/stake"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/event"
	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Client struct {
	bank.Bank
	event.Event
	stake.Stake
}

func NewClient(cfg types.SDKConfig) Client {
	cdc := types.NewAmino()
	rpc := net.NewRPCClient(cfg.NodeURI)

	ctx := &types.TxContext{
		Codec:   cdc,
		ChainID: cfg.ChainID,
		Online:  cfg.Online,
		KeyDAO:  cfg.KeyDAO,
		Network: cfg.Network,
		Mode:    cfg.Mode,
		RPC:     rpc,
	}

	baseClient := baseClient{ctx}
	client := Client{
		Bank:  bank.NewBankClient(baseClient),
		Event: event.NewEvent(baseClient),
		Stake: stake.NewStakeClient(baseClient),
	}
	client.setNetwork(ctx.Network)
	return client
}

func (c Client) setNetwork(network types.Network) {
	types.SetNetwork(network)
}

type baseClient struct {
	*types.TxContext
}

func (bm baseClient) Broadcast(baseTx types.BaseTx, msg []types.Msg) (types.Result, error) {
	err := bm.prepareTxContext(baseTx)
	if err != nil {
		return nil, err
	}
	tx, err := bm.BuildAndSign(baseTx.From, msg)
	if err != nil {
		return nil, err
	}
	return bm.broadcastTx(tx)
}

func (bm baseClient) Query(path string, data interface{}, result interface{}) error {
	bz, err := bm.Codec.MarshalJSON(data)
	if err != nil {
		return err
	}

	res, err := bm.RPC.Query(path, bz)
	if err != nil {
		return err
	}
	if err = bm.Codec.UnmarshalJSON(res, result); err != nil {
		return err
	}
	return nil
}

func (bm baseClient) QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, "subspace")
	opts := rpcclient.ABCIQueryOptions{
		//Height: cliCtx.Height,
		Prove: false,
	}

	result, err := bm.RPC.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.New(resp.Log)
	}
	return resp.Value, nil
}

func (bm baseClient) QueryAccount(address string) (baseAccount types.BaseAccount, err error) {
	addr, err := types.AccAddressFromBech32(address)
	if err != nil {
		return baseAccount, err
	}
	param := bank.QueryAccountParams{
		Address: addr,
	}
	if err = bm.Query("custom/acc/account", param, &baseAccount); err != nil {
		return baseAccount, err
	}
	return
}

func (bm baseClient) GetSender(name string) types.AccAddress {
	keyDAO := (*bm.TxContext).KeyDAO
	keystore := keyDAO.Read(name)
	return types.MustAccAddressFromBech32(keystore.GetAddress())
}

func (bm baseClient) GetRPC() net.RPCClient {
	return (*bm.TxContext).RPC
}

func (bm baseClient) GetCodec() types.Codec {
	return (*bm.TxContext).Codec
}

func (bm baseClient) prepareTxContext(baseTx types.BaseTx) error {
	ctx := bm.TxContext
	if ctx.Online {
		keyStore := ctx.KeyDAO.Read(baseTx.From)
		account, err := bm.QueryAccount(keyStore.GetAddress())
		if err != nil {
			return err
		}
		ctx = ctx.WithAccountNumber(account.AccountNumber).
			WithSequence(account.Sequence)
	}
	if len(baseTx.Gas) > 0 {
		gas, err := strconv.ParseUint(baseTx.Gas, 10, 64)
		if err != nil {
			return errors.New("gas must be either integer")
		}
		ctx = ctx.WithGas(gas)
	}

	if len(baseTx.Fee) > 0 {
		ctx = ctx.WithFee(baseTx.Fee)
	}

	if len(baseTx.Mode) > 0 {
		ctx = ctx.WithMode(baseTx.Mode)
	}

	if baseTx.Simulate {
		ctx = ctx.WithSimulate(baseTx.Simulate)
	}

	ctx = ctx.WithMemo(baseTx.Memo)
	return nil
}
func (bm baseClient) broadcastTx(txBytes []byte) (types.Result, error) {
	switch bm.Mode {
	case types.Commit:
		return bm.broadcastTxCommit(txBytes)
	case types.Async:
		return bm.broadcastTxAsync(txBytes)
	case types.Sync:
		return bm.broadcastTxSync(txBytes)

	}
	panic("invalid broadcast mode")
}

// broadcastTxCommit broadcasts transaction bytes to a Tendermint node
// and waits for a commit.
func (bm baseClient) broadcastTxCommit(tx []byte) (result types.ResultBroadcastTxCommit, err error) {
	res, err := bm.RPC.BroadcastTxCommit(tx)
	if err != nil {
		return result, err
	}

	if !res.CheckTx.IsOK() {
		return result, errors.New(res.CheckTx.Log)
	}

	if !res.DeliverTx.IsOK() {
		return result, errors.New(res.DeliverTx.Log)
	}
	return types.ResultBroadcastTxCommit{
		CheckTx:   res.CheckTx,
		DeliverTx: res.DeliverTx,
		Hash:      res.Hash,
		Height:    res.Height,
	}, err
}

// BroadcastTxSync broadcasts transaction bytes to a Tendermint node
// synchronously.
func (bm baseClient) broadcastTxSync(tx []byte) (result types.ResultBroadcastTxCommit, err error) {
	res, err := bm.RPC.BroadcastTxSync(tx)
	if err != nil {
		return result, err
	}

	return types.ResultBroadcastTxCommit{
		Hash: res.Hash,
		CheckTx: abci.ResponseCheckTx{
			Code: res.Code,
			Data: res.Data,
			Log:  res.Log,
		},
	}, nil
}

// BroadcastTxAsync broadcasts transaction bytes to a Tendermint node
// asynchronously.
func (bm baseClient) broadcastTxAsync(tx []byte) (result types.ResultBroadcastTx, err error) {
	res, err := bm.RPC.BroadcastTxAsync(tx)
	if err != nil {
		return result, err
	}

	return types.ResultBroadcastTx{
		Code: res.Code,
		Data: res.Data,
		Log:  res.Log,
		Hash: res.Hash,
	}, nil
}