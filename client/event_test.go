package client_test

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"time"
)

func (c *ClientTestSuite) TestSubscribeNewBlock() {
	err := c.SubscribeNewBlock(func(sub types.Subscription) {
		bz, _ := json.Marshal(sub.GetData())
		fmt.Println(string(bz))
		sub.Unsubscribe()
	})
	require.NoError(c.T(), err)
	time.Sleep(20 * time.Second)
}

func (c *ClientTestSuite) TestSubscribeTx() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	to := "iaa120v5ev44cwft687l0jcr5ec3vh2626vsschv7e"
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Async,
	}
	result, err := c.Send(to, coins, baseTx)
	require.NoError(c.T(), err)
	require.True(c.T(), result.IsSuccess())
	ch := make(chan int)
	query := types.EventQueryTxFor(result.GetHash())
	err = c.SubscribeTx(query, func(sub types.Subscription) {
		tx := sub.GetData().(types.EventDataTx)
		require.Equal(c.T(), result.GetHash(), tx.Hash)
		sub.Unsubscribe()
		ch <- 1
	})
	require.NoError(c.T(), err)
	<-ch
}