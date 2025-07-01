package order_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
)

type Client struct {
	httpClient *http.Client
	config     config.OrderServiceConfig
}

func NewClient(config config.OrderServiceConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.Token))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) GetAllAccountLoanPackages(ctx context.Context, accountNo string) ([]entity.AccountLoanPackage, error) {
	errorFormat := "GetAllAccountLoanPackages %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/v2/accounts/%s/loan-packages", c.config.Url, accountNo), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := resp.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		_ = json.NewDecoder(resp.Body).Decode(&errorResponse)
		return nil, fmt.Errorf(
			"GetAllAccountLoanPackages got error Status %d, Message: %s", resp.StatusCode, errorResponse.Message,
		)
	}
	dest := LoanPackagesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	return dest.LoanPackages, nil
}

func (c *Client) GetAccountByAccountNoAndCustodyCode(ctx context.Context, custodyCode string, accountNo string) (entity.OrderServiceAccount, error) {
	errorFormat := "GetAccountByAccountNoAndCustodyCode %w"
	request, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/internal/investors/%s/accounts/%s", c.config.Url, custodyCode, accountNo), nil)
	if err != nil {
		return entity.OrderServiceAccount{}, fmt.Errorf(errorFormat, err)
	}
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return entity.OrderServiceAccount{}, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := resp.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	if resp.StatusCode == http.StatusNotFound {
		return entity.OrderServiceAccount{}, apperrors.ErrAccountNoInvalid
	}
	if resp.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return entity.OrderServiceAccount{}, fmt.Errorf(errorFormat, err)
		}
		return entity.OrderServiceAccount{}, fmt.Errorf(
			"GetAccountByAccountNoAndCustodyCode got error Status %d, Message: %s", resp.StatusCode, errorResponse.Message,
		)
	}
	dest := entity.OrderServiceAccount{}
	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return entity.OrderServiceAccount{}, fmt.Errorf(errorFormat, err)
	}
	return dest, nil
}
