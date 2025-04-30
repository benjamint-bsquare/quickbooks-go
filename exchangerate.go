// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"encoding/json"
	"strconv"
)

type ExchangeRate struct {
	SyncToken          string               `json:"SyncToken,omitempty"`
	Domain             string               `json:"domain,omitempty"`
	AsOfDate           string               `json:"AsOfDate,omitempty"`
	SourceCurrencyCode string               `json:"SourceCurrencyCode,omitempty"`
	Rate               json.Number          `json:",omitempty"`
	Sparse             bool                 `json:"sparse,omitempty"`
	TargetCurrencyCode string               `json:"TargetCurrencyCode,omitempty"`
	MetaData           ModificationMetaData `json:",omitempty"`
}

// FindExchangeRates gets the full list of ExchangeRates in the QuickBooks account.
func (c *Client) FindExchangeRates(params RequestParameters) ([]ExchangeRate, error) {
	var resp struct {
		QueryResponse struct {
			ExchangeRates []ExchangeRate
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query(params, "SELECT COUNT(*) FROM ExchangeRate", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, nil
	}

	exchangeRates := make([]ExchangeRate, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += QueryPageSize {
		query := "SELECT * FROM ExchangeRate ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(QueryPageSize)

		if err := c.query(params, query, &resp); err != nil {
			return nil, err
		}

		exchangeRates = append(exchangeRates, resp.QueryResponse.ExchangeRates...)
	}

	return exchangeRates, nil
}

func (c *Client) FindExchangeRatesByPage(params RequestParameters, startPosition, pageSize int) ([]ExchangeRate, error) {
	var resp struct {
		QueryResponse struct {
			ExchangeRates []ExchangeRate
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	query := "SELECT * FROM ExchangeRate ORDERBY Id STARTPOSITION " + strconv.Itoa(startPosition) + " MAXRESULTS " + strconv.Itoa(pageSize)

	if err := c.query(params, query, &resp); err != nil {
		return nil, err
	}

	return resp.QueryResponse.ExchangeRates, nil
}

// FindExchangeRateById returns an exchangerate with a given Id.
func (c *Client) FindExchangeRateByCurrency(params RequestParameters, currencyCode string, asOf *Date) (*ExchangeRate, error) {
	var resp struct {
		ExchangeRate ExchangeRate
		Time         Date
	}

	qParams := map[string]string{
		"sourcecurrencycode": currencyCode,
	}
	if asOf != nil {
		qParams["asofdate"] = asOf.Format(dateFormat)
	}

	if err := c.get(params, "exchangerate", &resp, qParams); err != nil {
		return nil, err
	}

	return &resp.ExchangeRate, nil
}

// QueryExchangeRates accepts an SQL query and returns all exchangerate found using it
func (c *Client) QueryExchangeRates(params RequestParameters, query string) ([]ExchangeRate, error) {
	var resp struct {
		QueryResponse struct {
			ExchangeRates []ExchangeRate
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(params, query, &resp); err != nil {
		return nil, err
	}

	return resp.QueryResponse.ExchangeRates, nil
}
