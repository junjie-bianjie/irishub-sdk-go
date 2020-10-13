package service

import (
	"encoding/hex"
	json2 "encoding/json"
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
	ModuleName = "service"

	eventTypeNewBatchRequest         = "new_batch_request"
	eventTypeNewBatchRequestProvider = "new_batch_request_provider"
	attributeKeyRequests             = "requests"
	attributeKeyRequestID            = "request_id"
	attributeKeyRequestContextID     = "request_context_id"
	attributeKeyServiceName          = "service_name"
	attributeKeyProvider             = "provider"

	requestIDLen = 58
	contextIDLen = 40
)

var (
	_ original.Msg = MsgDefineService{}
	_ original.Msg = MsgBindService{}
	_ original.Msg = MsgUpdateServiceBinding{}
	_ original.Msg = MsgSetWithdrawAddress{}
	_ original.Msg = MsgDisableServiceBinding{}
	_ original.Msg = MsgEnableServiceBinding{}
	_ original.Msg = MsgRefundServiceDeposit{}
	_ original.Msg = MsgCallService{}
	_ original.Msg = MsgRespondService{}
	_ original.Msg = MsgPauseRequestContext{}
	_ original.Msg = MsgStartRequestContext{}
	_ original.Msg = MsgKillRequestContext{}
	_ original.Msg = MsgUpdateRequestContext{}
	_ original.Msg = MsgWithdrawEarnedFees{}

	cdc = original.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

// MsgDefineService defines a message to define a service
type MsgDefineService struct {
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Tags              []string            `json:"tags"`
	Author            original.AccAddress `json:"author"`
	AuthorDescription string              `json:"author_description"`
	Schemas           string              `json:"schemas"`
}

func (msg MsgDefineService) Route() string { return ModuleName }

func (msg MsgDefineService) Type() string {
	return "define_service"
}

func (msg MsgDefineService) ValidateBasic() error {
	if len(msg.Author) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Name) == 0 {
		return errors.New("author missing")
	}

	if len(msg.Schemas) == 0 {
		return errors.New("schemas missing")
	}

	return nil
}

func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgDefineService) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Author}
}

// MsgBindService defines a message to bind a service
type MsgBindService struct {
	ServiceName string              `json:"service_name"`
	Provider    original.AccAddress `json:"provider"`
	Deposit     original.Coins      `json:"deposit"`
	Pricing     string              `json:"pricing"`
	MinRespTime uint64              `json:"min_resp_time"`
}

func (msg MsgBindService) Type() string {
	return "bind_service"
}

func (msg MsgBindService) Route() string { return ModuleName }

func (msg MsgBindService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Pricing) == 0 {
		return errors.New("pricing missing")
	}
	return nil
}

func (msg MsgBindService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgBindService) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

// MsgCallService defines a message to request a service
type MsgCallService struct {
	ServiceName       string                `json:"service_name"`
	Providers         []original.AccAddress `json:"providers"`
	Consumer          original.AccAddress   `json:"consumer"`
	Input             string                `json:"input"`
	ServiceFeeCap     original.Coins        `json:"service_fee_cap"`
	Timeout           int64                 `json:"timeout"`
	SuperMode         bool                  `json:"super_mode"`
	Repeated          bool                  `json:"repeated"`
	RepeatedFrequency uint64                `json:"repeated_frequency"`
	RepeatedTotal     int64                 `json:"repeated_total"`
}

func (msg MsgCallService) Route() string { return ModuleName }

func (msg MsgCallService) Type() string {
	return "request_service"
}

func (msg MsgCallService) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	if len(msg.Providers) == 0 {
		return errors.New("providers missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("serviceName missing")
	}

	if len(msg.Input) == 0 {
		return errors.New("input missing")
	}
	return nil
}

func (msg MsgCallService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgCallService) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Consumer}
}

// MsgRespondService defines a message to respond a service request
type MsgRespondService struct {
	RequestID bytes.HexBytes      `json:"request_id"`
	Provider  original.AccAddress `json:"provider"`
	Result    string              `json:"result"`
	Output    string              `json:"output"`
}

func (msg MsgRespondService) Route() string { return ModuleName }

func (msg MsgRespondService) Type() string {
	return "respond_service"
}

func (msg MsgRespondService) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.Result) == 0 {
		return errors.New("result missing")
	}

	if len(msg.Output) > 0 {
		if !json2.Valid([]byte(msg.Output)) {
			return errors.New("output is not valid JSON")
		}
	}

	return nil
}

func (msg MsgRespondService) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

func (msg MsgRespondService) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgUpdateServiceBinding defines a message to update a service binding
type MsgUpdateServiceBinding struct {
	ServiceName string              `json:"service_name"`
	Provider    original.AccAddress `json:"provider"`
	Deposit     original.Coins      `json:"deposit"`
	Pricing     string              `json:"pricing"`
}

func (msg MsgUpdateServiceBinding) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return "update_service_binding" }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return errors.New(fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSetWithdrawAddress defines a message to set the withdrawal address for a service binding
type MsgSetWithdrawAddress struct {
	Provider        original.AccAddress `json:"provider"`
	WithdrawAddress original.AccAddress `json:"withdraw_address"`
}

func (msg MsgSetWithdrawAddress) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.WithdrawAddress) == 0 {
		return errors.New("withdrawal address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgDisableServiceBinding defines a message to disable a service binding
type MsgDisableServiceBinding struct {
	ServiceName string              `json:"service_name"`
	Provider    original.AccAddress `json:"provider"`
}

func (msg MsgDisableServiceBinding) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgDisableServiceBinding) Type() string { return "disable_service" }

// GetSignBytes implements Msg.
func (msg MsgDisableServiceBinding) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgDisableServiceBinding) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgEnableServiceBinding defines a message to enable a service binding
type MsgEnableServiceBinding struct {
	ServiceName string              `json:"service_name"`
	Provider    original.AccAddress `json:"provider"`
	Deposit     original.Coins      `json:"deposit"`
}

func (msg MsgEnableServiceBinding) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgEnableServiceBinding) Type() string { return "enable_service" }

// GetSignBytes implements Msg.
func (msg MsgEnableServiceBinding) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableServiceBinding) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	if !msg.Deposit.Empty() {
		return errors.New(fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableServiceBinding) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRefundServiceDeposit defines a message to refund deposit from a service binding
type MsgRefundServiceDeposit struct {
	ServiceName string              `json:"service_name"`
	Provider    original.AccAddress `json:"provider"`
}

func (msg MsgRefundServiceDeposit) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return "refund_service_deposit" }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	if len(msg.ServiceName) == 0 {
		return errors.New("service name missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgPauseRequestContext defines a message to suspend a request context
type MsgPauseRequestContext struct {
	RequestContextID bytes.HexBytes      `json:"request_context_id"`
	Consumer         original.AccAddress `json:"consumer"`
}

func (msg MsgPauseRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgPauseRequestContext) Type() string { return "pause_request_context" }

// GetSignBytes implements Msg.
func (msg MsgPauseRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgPauseRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgPauseRequestContext) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgStartRequestContext defines a message to resume a request context
type MsgStartRequestContext struct {
	RequestContextID bytes.HexBytes      `json:"request_context_id"`
	Consumer         original.AccAddress `json:"consumer"`
}

func (msg MsgStartRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgStartRequestContext) Type() string { return "start_request_context" }

// GetSignBytes implements Msg.
func (msg MsgStartRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgStartRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgStartRequestContext) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgKillRequestContext defines a message to terminate a request context
type MsgKillRequestContext struct {
	RequestContextID bytes.HexBytes      `json:"request_context_id"`
	Consumer         original.AccAddress `json:"consumer"`
}

func (msg MsgKillRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgKillRequestContext) Type() string { return "kill_request_context" }

// GetSignBytes implements Msg.
func (msg MsgKillRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgKillRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgKillRequestContext) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgUpdateRequestContext defines a message to update a request context
type MsgUpdateRequestContext struct {
	RequestContextID  bytes.HexBytes        `json:"request_context_id"`
	Providers         []original.AccAddress `json:"providers"`
	ServiceFeeCap     original.Coins        `json:"service_fee_cap"`
	Timeout           int64                 `json:"timeout"`
	RepeatedFrequency uint64                `json:"repeated_frequency"`
	RepeatedTotal     int64                 `json:"repeated_total"`
	Consumer          original.AccAddress   `json:"consumer"`
}

func (msg MsgUpdateRequestContext) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgUpdateRequestContext) Type() string { return "update_request_context" }

// GetSignBytes implements Msg.
func (msg MsgUpdateRequestContext) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateRequestContext) ValidateBasic() error {
	if len(msg.Consumer) == 0 {
		return errors.New("consumer missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateRequestContext) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgWithdrawEarnedFees defines a message to withdraw the fees earned by the provider
type MsgWithdrawEarnedFees struct {
	Provider original.AccAddress `json:"provider"`
}

func (msg MsgWithdrawEarnedFees) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgWithdrawEarnedFees) Type() string { return "withdraw_earned_fees" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawEarnedFees) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawEarnedFees) ValidateBasic() error {
	if len(msg.Provider) == 0 {
		return errors.New("provider missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawEarnedFees) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgWithdrawTax defines a message to withdraw the service tax
type MsgWithdrawTax struct {
	Trustee     original.AccAddress `json:"trustee"`
	DestAddress original.AccAddress `json:"dest_address"`
	Amount      original.Coins      `json:"amount"`
}

func (msg MsgWithdrawTax) Route() string { return ModuleName }

// Type implements Msg.
func (msg MsgWithdrawTax) Type() string { return "withdraw_tax" }

// GetSignBytes implements Msg.
func (msg MsgWithdrawTax) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawTax) ValidateBasic() error {
	if len(msg.Trustee) == 0 {
		return errors.New("trustee missing")
	}

	if len(msg.DestAddress) == 0 {
		return errors.New("destination address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawTax) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Trustee}
}

//==========================================for QueryWithResponse==========================================

// serviceDefinition represents a service definition
type serviceDefinition struct {
	Name              string              `json:"name"`
	Description       string              `json:"description"`
	Tags              []string            `json:"tags"`
	Author            original.AccAddress `json:"author"`
	AuthorDescription string              `json:"author_description"`
	Schemas           string              `json:"schemas"`
}

func (r serviceDefinition) Convert() interface{} {
	return rpc.ServiceDefinition{
		Name:              r.Name,
		Description:       r.Description,
		Tags:              r.Tags,
		Author:            r.Author,
		AuthorDescription: r.AuthorDescription,
		Schemas:           r.Schemas,
	}
}

// serviceBinding defines a struct for service binding
type serviceBinding struct {
	ServiceName  string              `json:"service_name"`
	Provider     original.AccAddress `json:"provider"`
	Deposit      original.Coins      `json:"deposit"`
	Pricing      string              `json:"pricing"`
	Qos          uint64              `json:"qos"`
	Available    bool                `json:"available"`
	DisabledTime time.Time           `json:"disabled_time"`
	Owner        original.AccAddress `json:"owner"`
}

func (b serviceBinding) Convert() interface{} {
	return rpc.ServiceBinding{
		ServiceName:  b.ServiceName,
		Provider:     b.Provider.String(),
		Deposit:      b.Deposit,
		Pricing:      b.Pricing,
		Qos:          b.Qos,
		Available:    b.Available,
		DisabledTime: b.DisabledTime,
		Owner:        b.Owner.String(),
	}

}

type serviceBindings []serviceBinding

func (bs serviceBindings) Convert() interface{} {
	bindings := make([]rpc.ServiceBinding, len(bs))
	for i, v := range bs {
		bindings[i] = v.Convert().(rpc.ServiceBinding)
	}
	return bindings
}

// request defines a request which contains the detailed request data
type request struct {
	ID                         string              `json:"id"`
	ServiceName                string              `json:"service_name"`
	Provider                   original.AccAddress `json:"provider"`
	Consumer                   original.AccAddress `json:"consumer"`
	Input                      string              `json:"input"`
	ServiceFee                 original.Coins      `json:"service_fee"`
	SuperMode                  bool                `json:"super_mode"`
	RequestHeight              int64               `json:"request_height"`
	ExpirationHeight           int64               `json:"expiration_height"`
	RequestContextID           bytes.HexBytes      `json:"request_context_id"`
	RequestContextBatchCounter uint64              `json:"request_context_batch_counter"`
}

func (r request) Empty() bool {
	return len(r.ServiceName) == 0
}

func (r request) Convert() interface{} {
	return rpc.ServiceRequest{
		ID:                         r.ID,
		ServiceName:                r.ServiceName,
		Provider:                   r.Provider,
		Consumer:                   r.Consumer,
		Input:                      r.Input,
		ServiceFee:                 r.ServiceFee,
		SuperMode:                  r.SuperMode,
		RequestHeight:              r.RequestHeight,
		ExpirationHeight:           r.ExpirationHeight,
		RequestContextID:           r.RequestContextID.String(),
		RequestContextBatchCounter: r.RequestContextBatchCounter,
	}
}

type requests []request

func (rs requests) Convert() interface{} {
	requests := make([]rpc.ServiceRequest, len(rs))
	for i, v := range rs {
		requests[i] = v.Convert().(rpc.ServiceRequest)
	}
	return requests
}

// ServiceResponse defines a response
type response struct {
	Provider                   original.AccAddress `json:"provider"`
	Consumer                   original.AccAddress `json:"consumer"`
	Output                     string              `json:"output"`
	Result                     string              `json:"result"`
	RequestContextID           bytes.HexBytes      `json:"request_context_id"`
	RequestContextBatchCounter uint64              `json:"request_context_batch_counter"`
}

func (r response) Empty() bool {
	return len(r.Provider) == 0
}

func (r response) Convert() interface{} {
	return rpc.ServiceResponse{
		Provider:                   r.Provider,
		Consumer:                   r.Consumer,
		Output:                     r.Output,
		Result:                     r.Result,
		RequestContextID:           r.RequestContextID.String(),
		RequestContextBatchCounter: r.RequestContextBatchCounter,
	}
}

type responses []response

func (rs responses) Convert() interface{} {
	responses := make([]rpc.ServiceResponse, len(rs))
	for i, v := range rs {
		responses[i] = v.Convert().(rpc.ServiceResponse)
	}
	return responses
}

// requestContext defines a context which holds request-related data
type requestContext struct {
	ServiceName        string                `json:"service_name"`
	Providers          []original.AccAddress `json:"providers"`
	Consumer           original.AccAddress   `json:"consumer"`
	Input              string                `json:"input"`
	ServiceFeeCap      original.Coins        `json:"service_fee_cap"`
	ModuleName         string                `json:"module_name"`
	Timeout            int64                 `json:"timeout"`
	SuperMode          bool                  `json:"super_mode"`
	Repeated           bool                  `json:"repeated"`
	RepeatedFrequency  uint64                `json:"repeated_frequency"`
	RepeatedTotal      int64                 `json:"repeated_total"`
	BatchCounter       uint64                `json:"batch_counter"`
	BatchRequestCount  uint32                `json:"batch_request_count"`
	BatchResponseCount uint32                `json:"batch_response_count"`
	ResponseThreshold  uint32                `json:"response_threshold"`
	BatchState         int32                 `json:"batch_state"`
	State              int32                 `json:"state"`
}

// Empty returns true if empty
func (r requestContext) Empty() bool {
	return len(r.ServiceName) == 0
}

func (r requestContext) Convert() interface{} {
	return rpc.RequestContext{
		ServiceName:        r.ServiceName,
		Providers:          r.Providers,
		Consumer:           r.Consumer,
		Input:              r.Input,
		ServiceFeeCap:      r.ServiceFeeCap,
		Timeout:            r.Timeout,
		SuperMode:          r.SuperMode,
		Repeated:           r.Repeated,
		RepeatedFrequency:  r.RepeatedFrequency,
		RepeatedTotal:      r.RepeatedTotal,
		BatchCounter:       r.BatchCounter,
		BatchRequestCount:  r.BatchRequestCount,
		BatchResponseCount: r.BatchResponseCount,
		BatchState:         r.BatchState,
		State:              r.State,
		ResponseThreshold:  r.ResponseThreshold,
		ModuleName:         r.ModuleName,
	}
}

// earnedFees defines a struct for the fees earned by the provider
type earnedFees struct {
	Address original.AccAddress `json:"address"`
	Coins   original.Coins      `json:"coins"`
}

func (e earnedFees) Convert() interface{} {
	return rpc.EarnedFees{
		Address: e.Address,
		Coins:   e.Coins,
	}
}

// CompactRequest defines a compact request with a request context ID
type compactRequest struct {
	RequestContextID           bytes.HexBytes
	RequestContextBatchCounter uint64
	Provider                   original.AccAddress
	ServiceFee                 original.Coins
	RequestHeight              int64
}

// service params
type Params struct {
	MaxRequestTimeout    int64          `json:"max_request_timeout"`
	MinDepositMultiple   int64          `json:"min_deposit_multiple"`
	MinDeposit           original.Coins `json:"min_deposit"`
	ServiceFeeTax        original.Dec   `json:"service_fee_tax"`
	SlashFraction        original.Dec   `json:"slash_fraction"`
	ComplaintRetrospect  time.Duration  `json:"complaint_retrospect"`
	ArbitrationTimeLimit time.Duration  `json:"arbitration_time_limit"`
	TxSizeLimit          uint64         `json:"tx_size_limit"`
	BaseDenom            string         `json:"base_denom"`
}

func (p Params) Convert() interface{} {
	return p
}

func registerCodec(cdc original.Codec) {
	cdc.RegisterConcrete(&MsgDefineService{}, "irismod/service/MsgDefineService")
	cdc.RegisterConcrete(&MsgBindService{}, "irismod/service/MsgBindService")
	cdc.RegisterConcrete(&MsgUpdateServiceBinding{}, "irismod/service/MsgUpdateServiceBinding")
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "irismod/service/MsgSetWithdrawAddress")
	cdc.RegisterConcrete(&MsgDisableServiceBinding{}, "irismod/service/MsgDisableServiceBinding")
	cdc.RegisterConcrete(&MsgEnableServiceBinding{}, "irismod/service/MsgEnableServiceBinding")
	cdc.RegisterConcrete(&MsgRefundServiceDeposit{}, "irismod/service/MsgRefundServiceDeposit")
	cdc.RegisterConcrete(&MsgCallService{}, "irismod/service/MsgCallService")
	cdc.RegisterConcrete(&MsgRespondService{}, "irismod/service/MsgRespondService")
	cdc.RegisterConcrete(&MsgPauseRequestContext{}, "irismod/service/MsgPauseRequestContext")
	cdc.RegisterConcrete(&MsgStartRequestContext{}, "irismod/service/MsgStartRequestContext")
	cdc.RegisterConcrete(&MsgKillRequestContext{}, "irismod/service/MsgKillRequestContext")
	cdc.RegisterConcrete(&MsgUpdateRequestContext{}, "irismod/service/MsgUpdateRequestContext")
	cdc.RegisterConcrete(&MsgWithdrawEarnedFees{}, "irismod/service/MsgWithdrawEarnedFees")
}

func actionTagKey(key ...string) original.EventKey {
	return original.EventKey(strings.Join(key, "."))
}

func hexBytesFrom(requestID string) bytes.HexBytes {
	v, _ := hex.DecodeString(requestID)
	return bytes.HexBytes(v)
}