package quickbooks

import (
	"errors"
	"strconv"
)

type PaymentMethod struct {
	MetaData  ModificationMetaData `json:",omitempty"`
	Id        string               `json:",omitempty"`
	Name      string               `json:",omitempty"`
	SyncToken string               `json:",omitempty"`
	Type      string               `json:",omitempty"`
	Active    bool                 `json:",omitempty"`
	Domain    string               `json:"domain,omitempty"`
	Status    string               `json:"status,omitempty"`
}

// CreatePaymentMethod creates the given PaymentMethod on the QuickBooks server, returning
// the resulting PaymentMethod object.
func (c *Client) CreatePaymentMethod(params RequestParameters, paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	var resp struct {
		PaymentMethod PaymentMethod
		Time          Date
	}

	if err := c.post(params, "paymentmethod", paymentMethod, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.PaymentMethod, nil
}

// FindPaymentMethods gets the full list of PaymentMethods in the QuickBooks account.
func (c *Client) FindPaymentMethods(params RequestParameters) ([]PaymentMethod, error) {
	var resp struct {
		QueryResponse struct {
			PaymentMethods []PaymentMethod `json:"PaymentMethod"`
			MaxResults     int
			StartPosition  int
			TotalCount     int
		}
	}

	if err := c.query(params, "SELECT COUNT(*) FROM PaymentMethod", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, nil
	}

	paymentMethods := make([]PaymentMethod, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += QueryPageSize {
		query := "SELECT * FROM PaymentMethod ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(QueryPageSize)

		if err := c.query(params, query, &resp); err != nil {
			return nil, err
		}

		paymentMethods = append(paymentMethods, resp.QueryResponse.PaymentMethods...)
	}

	return paymentMethods, nil
}

func (c *Client) FindPaymentMethodsByPage(params RequestParameters, startPosition, pageSize int) ([]PaymentMethod, error) {
	var resp struct {
		QueryResponse struct {
			PaymentMethods []PaymentMethod `json:"PaymentMethod"`
			MaxResults     int
			StartPosition  int
			TotalCount     int
		}
	}

	query := "SELECT * FROM PaymentMethod ORDERBY Id STARTPOSITION " + strconv.Itoa(startPosition) + " MAXRESULTS " + strconv.Itoa(pageSize)

	if err := c.query(params, query, &resp); err != nil {
		return nil, err
	}

	return resp.QueryResponse.PaymentMethods, nil
}

// FindPaymentMethodById finds the estimate by the given id
func (c *Client) FindPaymentMethodById(params RequestParameters, id string) (*PaymentMethod, error) {
	var resp struct {
		PaymentMethod PaymentMethod
		Time          Date
	}

	if err := c.get(params, "paymentmethod/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.PaymentMethod, nil
}

// QueryPaymentMethods accepts an SQL query and returns all estimates found using it
func (c *Client) QueryPaymentMethods(params RequestParameters, query string) ([]PaymentMethod, error) {
	var resp struct {
		QueryResponse struct {
			PaymentMethods []PaymentMethod `json:"PaymentMethod"`
			StartPosition  int
			MaxResults     int
		}
	}

	if err := c.query(params, query, &resp); err != nil {
		return nil, err
	}

	return resp.QueryResponse.PaymentMethods, nil
}

// UpdatePaymentMethod full updates the payment method, meaning that missing writable fields will be set to nil/null
func (c *Client) UpdatePaymentMethod(params RequestParameters, paymentMethod *PaymentMethod) (*PaymentMethod, error) {
	if paymentMethod.Id == "" {
		return nil, errors.New("missing estimate id")
	}

	existingPaymentMethod, err := c.FindPaymentMethodById(params, paymentMethod.Id)
	if err != nil {
		return nil, err
	}

	paymentMethod.SyncToken = existingPaymentMethod.SyncToken

	payload := struct {
		*PaymentMethod
	}{
		PaymentMethod: paymentMethod,
	}

	var paymentMethodData struct {
		PaymentMethod PaymentMethod
		Time          Date
	}

	if err = c.post(params, "estimate", payload, &paymentMethodData, nil); err != nil {
		return nil, err
	}

	return &paymentMethodData.PaymentMethod, err
}
