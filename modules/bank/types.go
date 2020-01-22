package bank

import "github.com/irisnet/irishub-sdk-go/types"

type Bank interface {
	GetAccount(address string) (types.BaseAccount, error)
	GetTokenStats(tokenID string) (types.TokenStats, error)
	Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error)
	Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error)
	SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error)
}

type bankClient struct {
	types.TxCtxManager
}

// defines the params for query: "custom/acc/account"
type QueryAccountParams struct {
	Address types.AccAddress
}

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{id}'
type QueryTokenParams struct {
	TokenId string
}