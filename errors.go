// Copyright (c) 2020, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

// Failure is the outermost struct that holds an error response.
type Failure struct {
	Fault struct {
		Error []struct {
			Message string
			Detail  string
			Code    string `json:"code"`
			Element string `json:"element"`
		}
		Type string `json:"type"`
	}
	Time Date `json:"time"`
}

// Error implements the error interface.
func (f Failure) Error() string {
	text, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("unexpected error while marshalling error: %v", err)
	}

	return string(text)
}

// parseFailure takes a response reader and tries to parse a Failure.
func parseFailure(resp *http.Response) error {
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	msg, err := io.ReadAll(reader)
	if err != nil {
		return errors.New("When reading response body:" + err.Error())
	}

	var errStruct Failure

	if err = json.Unmarshal(msg, &errStruct); err != nil {
		return errors.New(strconv.Itoa(resp.StatusCode) + " " + string(msg))
	}

	return errStruct
}
