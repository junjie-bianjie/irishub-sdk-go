package oracle

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type oracleClient struct {
	original.BaseClient
	*log.Logger
}

func (o oracleClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (o oracleClient) Name() string {
	return ModuleName
}

func Create(ac original.BaseClient) rpc.Oracle {
	return oracleClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

//CreateFeed create a stopped feed
func (o oracleClient) CreateFeed(request rpc.FeedCreateRequest) (original.ResultTx, original.Error) {
	creator, err := o.QueryAddress(request.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	var providers []original.AccAddress
	for _, provider := range request.Providers {
		p, err := original.AccAddressFromBech32(provider)
		if err != nil {
			return original.ResultTx{}, original.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	//amt, err := o.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgCreateFeed{
		FeedName:      request.FeedName,
		LatestHistory: request.LatestHistory,
		Description:   request.Description,
		Creator:       creator,
		ServiceName:   request.ServiceName,
		Providers:     providers,
		Input:         request.Input,
		Timeout:       request.Timeout,
		//ServiceFeeCap:     amt,
		RepeatedFrequency: request.RepeatedFrequency,
		AggregateFunc:     request.AggregateFunc,
		ValueJsonPath:     request.ValueJsonPath,
		ResponseThreshold: request.ResponseThreshold,
	}
	return o.BuildAndSend([]original.Msg{msg}, request.BaseTx)
}

//StartFeed start a stopped feed
func (o oracleClient) StartFeed(feedName string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	creator, err := o.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgStartFeed{
		FeedName: feedName,
		Creator:  creator,
	}
	return o.BuildAndSend([]original.Msg{msg}, baseTx)
}

//CreateAndStartFeed create and start a stopped feed
func (o oracleClient) CreateAndStartFeed(request rpc.FeedCreateRequest) (original.ResultTx, original.Error) {
	creator, err := o.QueryAddress(request.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	var providers []original.AccAddress
	for _, provider := range request.Providers {
		p, err := original.AccAddressFromBech32(provider)
		if err != nil {
			return original.ResultTx{}, original.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	//amt, err := o.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msgCreateFeed := MsgCreateFeed{
		FeedName:      request.FeedName,
		LatestHistory: request.LatestHistory,
		Description:   request.Description,
		Creator:       creator,
		ServiceName:   request.ServiceName,
		Providers:     providers,
		Input:         request.Input,
		Timeout:       request.Timeout,
		//ServiceFeeCap:     amt,
		RepeatedFrequency: request.RepeatedFrequency,
		AggregateFunc:     request.AggregateFunc,
		ValueJsonPath:     request.ValueJsonPath,
		ResponseThreshold: request.ResponseThreshold,
	}

	msgStartFeed := MsgStartFeed{
		FeedName: request.FeedName,
		Creator:  creator,
	}
	return o.BuildAndSend([]original.Msg{msgCreateFeed, msgStartFeed}, request.BaseTx)
}

//PauseFeed pause a running feed
func (o oracleClient) PauseFeed(feedName string, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	creator, err := o.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgPauseFeed{
		FeedName: feedName,
		Creator:  creator,
	}
	return o.BuildAndSend([]original.Msg{msg}, baseTx)
}

//EditFeed edit a feed
func (o oracleClient) EditFeed(request rpc.FeedEditRequest) (original.ResultTx, original.Error) {
	creator, err := o.QueryAddress(request.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	var providers []original.AccAddress
	for _, provider := range request.Providers {
		p, err := original.AccAddressFromBech32(provider)
		if err != nil {
			return original.ResultTx{}, original.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	//amt, err := o.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgEditFeed{
		FeedName:      request.FeedName,
		LatestHistory: request.LatestHistory,
		Description:   request.Description,
		Creator:       creator,
		Providers:     providers,
		Timeout:       request.Timeout,
		//ServiceFeeCap:     amt,
		RepeatedFrequency: request.RepeatedFrequency,
		ResponseThreshold: request.ResponseThreshold,
	}
	return o.BuildAndSend([]original.Msg{msg}, request.BaseTx)
}

//QueryFeed return the feed by feedName
func (o oracleClient) QueryFeed(feedName string) (rpc.FeedContext, original.Error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}

	var ctx feedContext
	if err := o.QueryWithResponse("custom/oracle/feed", param, &ctx); err != nil {
		return rpc.FeedContext{}, original.Wrap(err)
	}
	return ctx.Convert().(rpc.FeedContext), nil
}

//QueryFeeds return all feeds by state
func (o oracleClient) QueryFeeds(state string) ([]rpc.FeedContext, original.Error) {
	param := struct {
		State string
	}{
		State: state,
	}

	var fcs feedContexts
	if err := o.QueryWithResponse("custom/oracle/feeds", param, &fcs); err != nil {
		return nil, original.Wrap(err)
	}
	return fcs.Convert().([]rpc.FeedContext), nil
}

//QueryFeedValue return all feed values by feedName
func (o oracleClient) QueryFeedValue(feedName string) ([]rpc.FeedValue, original.Error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}

	var fvs feedValues
	if err := o.QueryWithResponse("custom/oracle/feedValue", param, &fvs); err != nil {
		return nil, original.Wrap(err)
	}
	return fvs.Convert().([]rpc.FeedValue), nil
}

func (o oracleClient) SubscribeFeedValue(feedName string, handler func(value rpc.FeedValue)) original.Error {
	feed, err := o.QueryFeed(feedName)
	if err != nil {
		return err
	}

	isInValidState := func(state string) bool {
		if state == COMPLETED || state == PAUSED || state == "" {
			return true
		}
		return false
	}

	if isInValidState(feed.State) {
		return original.Wrapf("feed:%s state is invalid:%s", feedName, feed.State)
	}

	handleResult := func(value string, sub1, sub2 original.Subscription) {
		o.Info().Str("feed-value", value).
			Msg("received feed value")
		var fv feedValue
		if err := cdc.UnmarshalJSON([]byte(value), &fv); err == nil {
			handler(fv.Convert().(rpc.FeedValue))
			f, err := o.QueryFeed(feedName)
			if err != nil || isInValidState(f.State) {
				_ = o.Unsubscribe(sub1)
				_ = o.Unsubscribe(sub2)
			}
		}
	}

	var sub1, sub2 original.Subscription

	blockBuilder := original.NewEventQueryBuilder().
		AddCondition(original.Cond(tagFeedName).Contains(original.EventValue(feedName)))
	sub1, err = o.SubscribeNewBlock(blockBuilder, func(block original.EventDataNewBlock) {
		tagValue := tagFeedValue(feedName)
		result := block.ResultEndBlock.Events.GetValues(tagValue, "")

		handleResult(result[0], sub1, sub2)
	})

	txBuilder := original.NewEventQueryBuilder().
		AddCondition(original.Cond(tagFeedName).Contains(original.EventValue(feedName))).
		AddCondition(original.Cond(original.ActionKey).EQ("respond_service"))
	sub2, err = o.SubscribeTx(txBuilder, func(tx original.EventDataTx) {
		tagValue := tagFeedValue(feedName)
		result := tx.Result.Events.GetValues(tagValue, "")

		handleResult(result[0], sub1, sub2)
	})
	return err
}