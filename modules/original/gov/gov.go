package gov

import (
	"encoding/json"
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type govClient struct {
	original.BaseClient
	*log.Logger
}

func Create(ac original.BaseClient) rpc.Gov {
	return govClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

//Deposit is responsible for depositing some tokens for proposal
func (g govClient) Deposit(proposalID uint64, amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	depositor, err := g.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	//amt, err := g.ToMinCoin(amount...)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		//Amount:     amt,
	}
	g.Info().
		Uint64("proposalID", proposalID).
		Str("depositor", depositor.String()).
		//Str("amount", amt.String()).
		Msg("execute gov deposit")
	return g.BuildAndSend([]original.Msg{msg}, baseTx)
}

//Vote is responsible for voting for proposal
func (g govClient) Vote(proposalID uint64, option rpc.VoteOption, baseTx original.BaseTx) (original.ResultTx, original.Error) {
	voter, err := g.QueryAddress(baseTx.From)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	op, err := VoteOptionFromString(option)
	if err != nil {
		return original.ResultTx{}, original.Wrap(err)
	}

	msg := MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     op,
	}
	g.Info().
		Uint64("proposalID", proposalID).
		Str("voter", voter.String()).
		Str("option", string(option)).
		Msg("execute gov vote")
	return g.BuildAndSend([]original.Msg{msg}, baseTx)
}

// QueryProposal returns the proposal of the specified proposalID
func (g govClient) QueryProposal(proposalID uint64) (rpc.Proposal, original.Error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	res, err := g.Query("custom/gov/proposal", param)
	if err != nil {
		return nil, original.Wrap(err)
	}

	var proposal BasicProposal
	//var proposal rpc.Proposal
	if err = json.Unmarshal(res, &proposal); err != nil {
		return nil, original.Wrap(err)
	}

	return proposal.Convert().(rpc.Proposal), nil
}

// QueryProposals returns all proposals of the specified params
func (g govClient) QueryProposals(request rpc.ProposalRequest) ([]rpc.Proposal, original.Error) {
	var voter, depositor original.AccAddress
	var err error
	if len(request.Voter) != 0 {
		voter, err = original.AccAddressFromBech32(request.Voter)
		if err != nil {
			return nil, original.Wrap(err)
		}
	}

	if len(request.Depositor) != 0 {
		depositor, err = original.AccAddressFromBech32(request.Depositor)
		if err != nil {
			return nil, original.Wrap(err)
		}
	}

	param := struct {
		Voter          original.AccAddress
		Depositor      original.AccAddress
		ProposalStatus string
		Limit          uint64
	}{
		Voter:          voter,
		Depositor:      depositor,
		ProposalStatus: request.ProposalStatus,
		Limit:          request.Limit,
	}

	res, err := g.Query("custom/gov/proposals", param)
	if err != nil {
		return nil, original.Wrap(err)
	}

	var proposals proposals
	if err := cdc.UnmarshalJSON(res, &proposals); err != nil {
		return nil, original.Wrap(err)
	}

	var ps []rpc.Proposal
	for _, p := range proposals {
		ps = append(ps, p.Convert().(rpc.Proposal))
	}
	return ps, nil
}

// QueryVote returns the vote of the specified proposalID and voter
func (g govClient) QueryVote(proposalID uint64, voter string) (rpc.Vote, original.Error) {
	v, err := original.AccAddressFromBech32(voter)
	if err != nil {
		return rpc.Vote{}, original.Wrap(err)
	}

	param := struct {
		ProposalID uint64
		Voter      original.AccAddress
	}{
		ProposalID: proposalID,
		Voter:      v,
	}

	var vote vote
	if err := g.QueryWithResponse("custom/gov/vote", param, &vote); err != nil {
		return rpc.Vote{}, original.Wrap(err)
	}
	return vote.Convert().(rpc.Vote), nil
}

// QueryVotes returns all votes of the specified proposalID
func (g govClient) QueryVotes(proposalID uint64) ([]rpc.Vote, original.Error) {
	param := struct {
		ProposalID uint64
		Page       int
	}{
		ProposalID: proposalID,
		Page:       1, // A page number must be passed in (pass default page:1)
	}

	var vs votes
	if err := g.QueryWithResponse("custom/gov/votes", param, &vs); err != nil {
		return nil, original.Wrap(err)
	}
	return vs.Convert().([]rpc.Vote), nil
}

// QueryDeposit returns the deposit of the specified proposalID and depositor
func (g govClient) QueryDeposit(proposalID uint64, depositor string) (rpc.Deposit, original.Error) {
	d, err := original.AccAddressFromBech32(depositor)
	if err != nil {
		return rpc.Deposit{}, original.Wrap(err)
	}

	param := struct {
		ProposalID uint64
		Depositor  original.AccAddress
	}{
		ProposalID: proposalID,
		Depositor:  d,
	}

	var deposit deposit
	if err := g.QueryWithResponse("custom/gov/deposit", param, &deposit); err != nil {
		return rpc.Deposit{}, original.Wrap(err)
	}
	return deposit.Convert().(rpc.Deposit), nil
}

// QueryDeposits returns all deposits of the specified proposalID
func (g govClient) QueryDeposits(proposalID uint64) ([]rpc.Deposit, original.Error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	var deposits deposits
	err := g.QueryWithResponse("custom/gov/deposits", param, &deposits)
	if err != nil {
		return nil, original.Wrap(err)
	}
	return deposits.Convert().([]rpc.Deposit), nil
}

// QueryTally returns the result of proposal by the specified proposalID
func (g govClient) QueryTally(proposalID uint64) (rpc.TallyResult, original.Error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	var tally tallyResult
	err := g.QueryWithResponse("custom/gov/tally", param, &tally)
	if err != nil {
		return rpc.TallyResult{}, original.Wrap(err)
	}
	return tally.Convert().(rpc.TallyResult), nil
}

func (g govClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (g govClient) Name() string {
	return ModuleName
}