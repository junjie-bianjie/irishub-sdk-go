package oracle

import (
	"errors"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/tendermint/tendermint/libs/bytes"
	"strings"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName = "oracle"
	RUNNING    = "running"
	PAUSED     = "paused"
	COMPLETED  = "completed"

	tagFeedName = "feed-name"
)

var (
	_ original.Msg = MsgCreateFeed{}
	_ original.Msg = MsgStartFeed{}
	_ original.Msg = MsgPauseFeed{}
	_ original.Msg = MsgEditFeed{}

	cdc = original.NewAminoCodec()

	tagFeedValue = func(feedName string) string {
		return fmt.Sprintf("%s.%s", tagFeedName, feedName)
	}
)

func init() {
	registerCodec(cdc)
}

//______________________________________________________________________

// MsgCreateFeed - struct for create a feed
type MsgCreateFeed struct {
	FeedName          string                `json:"feed_name"`
	LatestHistory     uint64                `json:"latest_history"`
	Description       string                `json:"description"`
	Creator           original.AccAddress   `json:"creator"`
	ServiceName       string                `json:"service_name"`
	Providers         []original.AccAddress `json:"providers"`
	Input             string                `json:"input"`
	Timeout           int64                 `json:"timeout"`
	ServiceFeeCap     original.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64                `json:"repeated_frequency"`
	AggregateFunc     string                `json:"aggregate_func"`
	ValueJsonPath     string                `json:"value_json_path"`
	ResponseThreshold uint16                `json:"response_threshold"`
}

func (msg MsgCreateFeed) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgCreateFeed) Type() string {
	return "create_feed"
}

// ValidateBasic implements Msg.
func (msg MsgCreateFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}

	serviceName := strings.TrimSpace(msg.ServiceName)
	if len(serviceName) == 0 {
		return errors.New("serviceName missed")
	}

	if len(msg.Providers) == 0 {
		return errors.New("providers missed")
	}

	aggregateFunc := strings.TrimSpace(msg.AggregateFunc)
	if len(aggregateFunc) == 0 {
		return errors.New("aggregateFunc missed")
	}

	valueJsonPath := strings.TrimSpace(msg.ValueJsonPath)
	if len(valueJsonPath) == 0 {
		return errors.New("valueJsonPath missed")
	}

	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreateFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgCreateFeed) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgStartFeed - struct for start a feed
type MsgStartFeed struct {
	FeedName string              `json:"feed_name"`
	Creator  original.AccAddress `json:"creator"`
}

func (msg MsgStartFeed) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgStartFeed) Type() string {
	return "start_feed"
}

// ValidateBasic implements Msg.
func (msg MsgStartFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}
	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgStartFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgStartFeed) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgPauseFeed - struct for stop a started feed
type MsgPauseFeed struct {
	FeedName string              `json:"feed_name"`
	Creator  original.AccAddress `json:"creator"`
}

func (msg MsgPauseFeed) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgPauseFeed) Type() string {
	return "pause_feed"
}

// ValidateBasic implements Msg.
func (msg MsgPauseFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}
	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgPauseFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgPauseFeed) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Creator}
}

//______________________________________________________________________

// MsgEditFeed - struct for edit a existed feed
type MsgEditFeed struct {
	FeedName          string                `json:"feed_name"`
	Description       string                `json:"description"`
	LatestHistory     uint64                `json:"latest_history"`
	Providers         []original.AccAddress `json:"providers"`
	Timeout           int64                 `json:"timeout"`
	ServiceFeeCap     original.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64                `json:"repeated_frequency"`
	ResponseThreshold uint16                `json:"response_threshold"`
	Creator           original.AccAddress   `json:"creator"`
}

func (msg MsgEditFeed) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgEditFeed) Type() string {
	return "edit_feed"
}

// ValidateBasic implements Msg.
func (msg MsgEditFeed) ValidateBasic() error {
	feedName := strings.TrimSpace(msg.FeedName)
	if len(feedName) == 0 {
		return errors.New("feedName missed")
	}

	if len(msg.Creator) == 0 {
		return errors.New("creator missed")
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgEditFeed) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg.
func (msg MsgEditFeed) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Creator}
}

//-------------------------------for query--------------------------
type feed struct {
	FeedName         string              `json:"feed_name"`
	Description      string              `json:"description"`
	AggregateFunc    string              `json:"aggregate_func"`
	ValueJsonPath    string              `json:"value_json_path"`
	LatestHistory    uint64              `json:"latest_history"`
	RequestContextID bytes.HexBytes      `json:"request_context_id"`
	Creator          original.AccAddress `json:"creator"`
}

type feedContext struct {
	Feed              feed                  `json:"feed"`
	ServiceName       string                `json:"service_name"`
	Providers         []original.AccAddress `json:"providers"`
	Input             string                `json:"input"`
	Timeout           int64                 `json:"timeout"`
	ServiceFeeCap     original.Coins        `json:"service_fee_cap"`
	RepeatedFrequency uint64                `json:"repeated_frequency"`
	ResponseThreshold uint16                `json:"response_threshold"`
	State             string                `json:"state"`
}

func (fc feedContext) Convert() interface{} {
	var providers []string
	for _, provider := range fc.Providers {
		providers = append(providers, provider.String())
	}
	return rpc.FeedContext{
		Feed: rpc.Feed{
			FeedName:         fc.Feed.FeedName,
			Description:      fc.Feed.Description,
			AggregateFunc:    fc.Feed.AggregateFunc,
			ValueJsonPath:    fc.Feed.ValueJsonPath,
			LatestHistory:    fc.Feed.LatestHistory,
			RequestContextID: fc.Feed.RequestContextID.String(),
			Creator:          fc.Feed.Creator.String(),
		},
		ServiceName:       fc.ServiceName,
		Providers:         providers,
		Input:             fc.Input,
		Timeout:           fc.Timeout,
		ServiceFeeCap:     fc.ServiceFeeCap,
		RepeatedFrequency: fc.RepeatedFrequency,
		ResponseThreshold: fc.ResponseThreshold,
		State:             fc.State,
	}
}

type feedContexts []feedContext

func (fcs feedContexts) Convert() interface{} {
	result := make([]rpc.FeedContext, len(fcs))
	for i, fc := range fcs {
		result[i] = fc.Convert().(rpc.FeedContext)
	}
	return result
}

type feedValue struct {
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func (fv feedValue) Convert() interface{} {
	return rpc.FeedValue{
		Data:      fv.Data,
		Timestamp: fv.Timestamp,
	}
}

type feedValues []feedValue

func (fvs feedValues) Convert() interface{} {
	result := make([]rpc.FeedValue, len(fvs))
	for i, fv := range fvs {
		result[i] = rpc.FeedValue{
			Data:      fv.Data,
			Timestamp: fv.Timestamp,
		}
	}
	return result
}

func registerCodec(cdc original.Codec) {
	cdc.RegisterConcrete(MsgCreateFeed{}, "irishub/oracle/MsgCreateFeed")
	cdc.RegisterConcrete(MsgStartFeed{}, "irishub/oracle/MsgStartFeed")
	cdc.RegisterConcrete(MsgPauseFeed{}, "irishub/oracle/MsgPauseFeed")
	cdc.RegisterConcrete(MsgEditFeed{}, "irishub/oracle/MsgEditFeed")

	cdc.RegisterConcrete(feed{}, "irishub/oracle/Feed")
	cdc.RegisterConcrete(feedContext{}, "irishub/oracle/FeedContext")
}