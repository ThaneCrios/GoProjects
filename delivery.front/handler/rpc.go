package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	httpClient = http.Client{
		Timeout: 30 * time.Second,
	}
)

func (h *Handler) RPC(method, url string, params map[string]string, payload, response interface{}) error {
	var req *http.Request
	var err error

	sParams := ""
	if params != nil {
		for k, v := range params {
			if sParams == "" {
				sParams += "?" + k + "=" + v
			} else {
				sParams += "&" + k + "=" + v
			}
		}
	}

	fmt.Print(url + sParams)

	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return errors.Wrap(err, "failed to marshal a payload")
		}
		req, err = http.NewRequest(method, url+sParams, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, url+sParams, nil)
	}
	if err != nil {
		return errors.Wrap(err, "failed to create an http request")
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to make a %s request", method))
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Errorf("%s request to %s failed with status: %d", method, url, resp.StatusCode)
		}
		return errors.Errorf("%s request to %s failed with status: %d and body: %s", method, url, resp.StatusCode, string(body))
	}

	if response != nil {
		return json.NewDecoder(resp.Body).Decode(response)
	}

	return nil
}
