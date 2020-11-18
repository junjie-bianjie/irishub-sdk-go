package random

import sdk "github.com/irisnet/irishub-sdk-go/types"

// expose Random module api for user
type RandomI interface {
	sdk.Module

	QueryRandom(ReqId string) (QueryRandomResp, sdk.Error)
	QueryRandomRequestQueue(height int64) ([]QueryRandomRequestQueueResp, sdk.Error)
}

type QueryRandomResp struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"`
	Height        int64  `json:"height" yaml:"height"`
	Value         string `json:"value" yaml:"value"`
}

type QueryRandomRequestQueueResp struct {
	Height           int64     `json:"height" yaml:"height"`
	Consumer         string    `json:"consumer" yaml:"consumer"`
	TxHash           string    `json:"tx_hash" yaml:"tx_hash"`
	Oracle           bool      `json:"oracle" yaml:"oracle"`
	ServiceFeeCap    sdk.Coins `json:"service_fee_cap" yaml:"service_fee_cap"`
	ServiceContextID string    `json:"service_context_id" yaml:"service_context_id"`
}
