package flex_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"financing-offer/internal/config"
)

type Client struct {
	httpClient *http.Client
	config     config.FlexOpenApiConfig
}

func NewClient(config config.FlexOpenApiConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) IsHOActive(ctx context.Context) (bool, error) {
	errorFormat := "IsHOActive %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/system/checkHOStatus", c.config.Url), nil)
	if err != nil {
		return false, fmt.Errorf(errorFormat, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return false, fmt.Errorf("IsHOActive %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return false, fmt.Errorf("IsHOActive %s got error Status %d, message: %s", req.URL.String(), res.StatusCode, string(rawMessage))
	}
	dest := IsHOActiveResponse{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return false, fmt.Errorf(errorFormat, err)
	}
	if dest.ErrorCode != "ok" {
		return false, fmt.Errorf("got error code %s", dest.ErrorCode)
	}
	return dest.HOStatus == "1", nil
}

func (c *Client) NewRequest(ctx context.Context, method, url string, body any) (*http.Request, error) {
	marshaledBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("NewRequest %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(marshaledBody))
	if err != nil {
		return nil, fmt.Errorf("NewRequest %w", err)
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
