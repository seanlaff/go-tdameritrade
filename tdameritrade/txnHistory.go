package tdameritrade

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

// Transactions is a slice of transactions
type Transactions []*Transaction

// Transaction represents a single transaction
type Transaction struct {
	Type                          string          `json:"type"`
	ClearingReferenceNumber       string          `json:"clearingReferenceNumber"`
	SubAccount                    string          `json:"subAccount"`
	SettlementDate                string          `json:"settlementDate"`
	OrderID                       string          `json:"orderId"`
	SMA                           float64         `json:"sma"`
	RequirementReallocationAmount float64         `json:"requirementReallocationAmount"`
	DayTradeBuyingPowerEffect     float64         `json:"dayTradeBuyingPowerEffect"`
	NetAmount                     float64         `json:"netAmount"`
	TransactionDate               string          `json:"transactionDate"`
	OrderDate                     string          `json:"orderDate"`
	TransactionSubType            string          `json:"transactionSubType"`
	TransactionID                 int64           `json:"transactionId"`
	CashBalanceEffectFlag         bool            `json:"cashBalanceEffectFlag"`
	ACHStatus                     string          `json:"achStatus"`
	AccruedInterest               float64         `json:"accruedInterest"`
	Fees                          TransactionFees `json:"fees"`
	TransactionItem               TransactionItem `json:"transactionItem"`
}

// TransactionFees contains fees related to the transaction
type TransactionFees struct {
	AdditionalFee float64 `json:"additionalFee"`
	CdscFee       float64 `json:"cdscFee"`
	Commission    float64 `json:"commission"`
	OptRegFee     float64 `json:"optRegFee"`
	OtherCharges  float64 `json:"otherCharges"`
	RFee          float64 `json:"rFee"`
	RegFee        float64 `json:"regFee"`
	SecFee        float64 `json:"secFee"`
}

// TransactionItem is an item within a transaction response
type TransactionItem struct {
	AccountID            int32                 `json:"accountId"`
	Amount               float64               `json:"amount"`
	Price                float64               `json:"price"`
	Cost                 float64               `json:"cost"`
	ParentOrderKey       int32                 `json:"parentOrderKey"`
	ParentChildIndicator string                `json:"parentChildIndicator"`
	Instruction          string                `json:"instruction"`
	PositionEffect       string                `json:"positionEffect"`
	Instrument           TransactionInstrument `json:"instrument"`
}

// TransactionInstrument is the instrumnet traded within a transaction
type TransactionInstrument struct {
	Symbol               string  `json:"symbol"`
	UnderlyingSymbol     string  `json:"underlyingSymbol"`
	OptionExpirationDate string  `json:"optionExpirationDate"`
	OptionStrikePrice    float64 `json:"optionStrikePrice"`
	PutCall              string  `json:"putCall"`
	CUSIP                string  `json:"cusip"`
	Description          string  `json:"description"`
	AssetType            string  `json:"assetType"`
	BondMaturityDate     string  `json:"bondMaturityDate"`
	BondInterestRate     float64 `json:"bondInterestRate"`
}

// TransactionHistoryOptions is parsed and translated to query options in the https request
type TransactionHistoryOptions struct {
	Type      string    `url:"type,omitempty"`
	Symbol    string    `url:"symbol,omitempty"`
	EndDate   time.Time `url:"endDate,omitempty"`
	StartDate time.Time `url:"startDate,omitempty"`
}

// EncodeValues handles converting TransactionHistoryOptions into the right shape (date formats)
func (t *TransactionHistoryOptions) EncodeValues(key string, v *url.Values) error {
	switch key {
	case "Type":
		v.Add(key, t.Type)
	case "Symbol":
		v.Add(key, t.Symbol)
	case "EndDate":
		v.Add(key, t.EndDate.Format("2006-01-02"))
	case "StartDate":
		v.Add(key, t.StartDate.Format("2006-01-02"))
	}
	return nil
}

// TransactionHistoryService handles communication with the transaction history related methods of
// the TDAmeritrade API.
//
// TDAmeritrade API docs: https://developer.tdameritrade.com/transaction-history/apis
type TransactionHistoryService struct {
	client *Client
}

// GetTransaction gets a specific transaction by account
// TDAmeritrade API Docs: https://developer.tdameritrade.com/transaction-history/apis/get/accounts/%7BaccountId%7D/transactions/%7BtransactionId%7D-0
func (s *TransactionHistoryService) GetTransaction(ctx context.Context, accountID string, transactionID string) (*Transaction, *Response, error) {
	u := fmt.Sprintf("accounts/%s/transactions/%s", accountID, transactionID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	txn := new(Transaction)
	resp, err := s.client.Do(ctx, req, txn)
	if err != nil {
		return nil, resp, err
	}
	return txn, resp, nil
}

// GetTransactions gets all transaction by account
// TDAmeritrade API Docs: https://developer.tdameritrade.com/transaction-history/apis/get/accounts/%7BaccountId%7D/transactions-0
func (s *TransactionHistoryService) GetTransactions(ctx context.Context, accountID string, opts *TransactionHistoryOptions) (*Transactions, *Response, error) {
	u := fmt.Sprintf("accounts/%s/transactions", accountID)
	if opts != nil {
		q, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		u = fmt.Sprintf("%s?%s", u, q.Encode())
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	txns := new(Transactions)
	resp, err := s.client.Do(ctx, req, txns)
	if err != nil {
		return nil, resp, err
	}
	return txns, resp, nil
}
