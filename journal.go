package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

type GlobalTaxCalculationEnum string

const (
	JournalEntryLineDetailTypeTaxExcluded  GlobalTaxCalculationEnum = "TaxExcluded"
	JournalEntryLineDetailTypeTaxInclusive GlobalTaxCalculationEnum = "TaxInclusive"
)

type JournalEntry struct {
	SyncToken            *string                   `json:"SyncToken,omitempty"`
	Domain               string                    `json:"domain"`
	TxnDate              Date                      `json:",omitempty"`
	Sparse               bool                      `json:"sparse"`
	Line                 []JournalEntryLine        `json:"Line"`
	CurrencyRef          *ReferenceType            `json:",omitempty"`
	Adjustment           bool                      `json:"Adjustment"`
	Id                   *string                   `json:"Id,omitempty"`
	TxnTaxDetail         *TxnTaxDetail             `json:",omitempty"`
	MetaData             *ModificationMetaData     `json:",omitempty"`
	GlobalTaxCalculation *GlobalTaxCalculationEnum `json:",omitempty"`
	DocNumber            *string                   `json:"DocNumber,omitempty"`
	PrivateNote          *string                   `json:"PrivateNote,omitempty"`
	ExchangeRate         json.Number               `json:",omitempty"`
}

type JournalEntryLineDetailTypeEnum string

const (
	JournalEntryLineDetailType JournalEntryLineDetailTypeEnum = "JournalEntryLineDetail"
)

type JournalEntryLine struct {
	Description            *string                        `json:"Description,omitempty"`
	LineNum                *json.Number                   `json:"LineNum,omitempty"`
	JournalEntryLineDetail JournalEntryLineDetail         `json:"JournalEntryLineDetail"`
	DetailType             JournalEntryLineDetailTypeEnum `json:"DetailType"`
	ProjectRef             *ReferenceType                 `json:"ProjectRef,omitempty"`
	Amount                 float64                        `json:"Amount"`
	Id                     *string                        `json:"Id,omitempty"`
}

type JournalEntryLineDetailPostingTypeEnum string

const (
	JournalEntryLineDetailPostingTypeCredit JournalEntryLineDetailPostingTypeEnum = "Credit"
	JournalEntryLineDetailPostingTypeDebit  JournalEntryLineDetailPostingTypeEnum = "Debit"
)

type JournalEntryLineDetailTaxApplicableOnEnum string

const (
	JournalEntryLineDetailTaxApplicableOnSales    JournalEntryLineDetailTaxApplicableOnEnum = "Sales"
	JournalEntryLineDetailTaxApplicableOnPurchase JournalEntryLineDetailTaxApplicableOnEnum = "Purchase"
)

type JournalEntryLineDetail struct {
	PostingType     JournalEntryLineDetailPostingTypeEnum      `json:"PostingType"`
	AccountRef      ReferenceType                              `json:"AccountRef"`
	Entity          *JournalEntryLineDetailEntity              `json:"Entity,omitempty"`
	JournalCodeRef  *ReferenceType                             `json:"JournalCodeRef,omitempty"`
	TaxApplicableOn *JournalEntryLineDetailTaxApplicableOnEnum `json:"TaxApplicableOn,omitempty"`
}

type JournalEntryLineDetailEntityTypeEnum string

const (
	JournalEntryLineDetailEntityTypeEnumVendor   JournalEntryLineDetailEntityTypeEnum = "Vendor"
	JournalEntryLineDetailEntityTypeEnumEmployee JournalEntryLineDetailEntityTypeEnum = "Employee"
	JournalEntryLineDetailEntityTypeEnumCustomer JournalEntryLineDetailEntityTypeEnum = "Customer"
)

type JournalEntryLineDetailEntity struct {
	Type      JournalEntryLineDetailEntityTypeEnum `json:"Type"`
	EntityRef ReferenceType                        `json:"EntityRef"`
}

// CreateAccount creates the given account within QuickBooks
func (c *Client) CreateJournalEntry(params RequestParameters, entry *JournalEntry) (*JournalEntry, error) {
	var resp struct {
		JournalEntry JournalEntry
		Time         Date
	}

	if err := c.post(params, "journalentry", entry, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.JournalEntry, nil
}

// FindJournalEntries gets the full list of JournalEntries in the QuickBooks account.
func (c *Client) FindJournalEntries(params RequestParameters) ([]JournalEntry, error) {
	var resp struct {
		QueryResponse struct {
			JournalEntry  []JournalEntry
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query(params, "SELECT COUNT(*) FROM JournalEntry", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, nil
	}

	accounts := make([]JournalEntry, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += QueryPageSize {
		query := "SELECT * FROM JournalEntry ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(QueryPageSize)

		if err := c.query(params, query, &resp); err != nil {
			return nil, err
		}

		accounts = append(accounts, resp.QueryResponse.JournalEntry...)
	}

	return accounts, nil
}

// FindJournalEntryById returns an account with a given Id.
func (c *Client) FindJournalEntryById(params RequestParameters, id string) (*JournalEntry, error) {
	var resp struct {
		JournalEntry JournalEntry
		Time         Date
	}

	if err := c.get(params, "journalentry/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.JournalEntry, nil
}

// UpdateJournalEntry full updates the JournalEntry, meaning that missing writable fields will be set to nil/null
func (c *Client) UpdateJournalEntry(params RequestParameters, entry *JournalEntry) (*JournalEntry, error) {
	if entry.Id == nil {
		return nil, errors.New("missing entry id")
	}

	existingEntry, err := c.FindJournalEntryById(params, *entry.Id)
	if err != nil {
		return nil, err
	}

	entry.SyncToken = existingEntry.SyncToken

	payload := struct {
		*JournalEntry
	}{
		JournalEntry: entry,
	}

	var entryData struct {
		JournalEntry JournalEntry
		Time         Date
	}

	if err = c.post(params, "journalentry", payload, &entryData, nil); err != nil {
		return nil, err
	}

	return &entryData.JournalEntry, err
}
