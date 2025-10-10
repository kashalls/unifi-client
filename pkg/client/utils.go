package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UnifiErrorResponse struct {
	Code      string                 `json:"code"`
	Details   map[string]interface{} `json:"details"`
	ErrorCode int                    `json:"errorCode"`
	Message   string                 `json:"message"`
}

func (c *Unifi) Do(method, path string, body io.Reader) (*http.Response, error) {
	endpoint, err := c.Endpoint(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Config.ApiKey == "" {
		if csrf := resp.Header.Get("X-CSRF-Token"); csrf != "" {
			c.csrf = csrf
		}

		if resp.StatusCode == http.StatusUnauthorized {
			if err := c.Login(); err != nil {
				return nil, err
			}
			c.setHeaders(req)

			resp, err = c.Client.Do(req)
			if err != nil {
				return nil, err
			}
		}
	}

	// It is unknown at this time if the UniFi API returns anything other than 200 for these types of requests.
	if resp.StatusCode != http.StatusOK {
		body, bodyErr := io.ReadAll(io.LimitReader(resp.Body, 512))
		if bodyErr != nil {
			return nil, bodyErr
		}

		var apiError UnifiErrorResponse
		if err := json.Unmarshal(body, &apiError); err != nil {
			return nil, fmt.Errorf("failed to decode json: %w", err)
		}

		return nil, fmt.Errorf("%s request to %s returned %d: %s", method, path, resp.StatusCode, apiError.Message)
	}

	return resp, nil
}

func (c *Unifi) setHeaders(req *http.Request) {
	if c.Config.ApiKey != "" {
		req.Header.Set("X-API-KEY", c.Config.ApiKey)
	} else {
		req.Header.Set("X-CSRF-Token", c.csrf)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
}
