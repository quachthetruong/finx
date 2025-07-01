package mo_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/marginoperation/repository"
	"financing-offer/internal/funcs"
)

var _ repository.MarginOperationRepository = (*Client)(nil)

type Client struct {
	httpClient *http.Client
	config     config.MoServiceConfig
}

func NewClient(config config.MoServiceConfig) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		config: config,
	}
}

func (c *Client) GetMarginPoolById(ctx context.Context, marginPoolGroupId int64) (entity.MarginPool, error) {
	errorFormat := "GetMarginPoolGroupsByIds %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/pools/%d", c.config.Url, marginPoolGroupId), nil)
	if err != nil {
		return entity.MarginPool{}, fmt.Errorf(errorFormat, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.MarginPool{}, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	dest := entity.MarginPool{}
	if err := c.handleResponse(res, errorFormat, &dest); err != nil {
		return entity.MarginPool{}, fmt.Errorf(errorFormat, err)
	}
	return dest, nil
}

func (c *Client) GetMarginPoolGroupsByIds(ctx context.Context, marginPoolGroupIds []int64) ([]entity.MarginPoolGroup, error) {
	errorFormat := "GetMarginPoolGroupsByIds %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/pool-groups", c.config.Url), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	q := req.URL.Query()
	q.Add(
		"ids",
		strings.Join(funcs.Map(marginPoolGroupIds, func(id int64) string { return strconv.FormatInt(id, 10) }), ","),
	)
	req.URL.RawQuery = q.Encode()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	dest := ListResponse[entity.MarginPoolGroup]{}
	if err := c.handleResponse(res, errorFormat, &dest); err != nil {
		return nil, err
	}
	return dest.Data, nil
}

func (c *Client) GetMarginPoolsByIds(ctx context.Context, marginPoolIds []int64) ([]entity.MarginPool, error) {
	errorFormat := "GetMarginPoolsByIds %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/pools", c.config.Url), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	q := req.URL.Query()
	q.Add(
		"ids", strings.Join(funcs.Map(marginPoolIds, func(id int64) string { return strconv.FormatInt(id, 10) }), ","),
	)
	req.URL.RawQuery = q.Encode()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	dest := ListResponse[entity.MarginPool]{}
	if err := c.handleResponse(res, errorFormat, &dest); err != nil {
		return nil, err
	}
	return dest.Data, nil
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

func (c *Client) handleResponse(res *http.Response, errorFormat string, dest any) error {
	// close response body
	defer func() {
		res.Body.Close()
	}()
	if res.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		if err := json.NewDecoder(res.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf(errorFormat, err)
		}
		return apperrors.AppError{
			Err:     errorResponse,
			Code:    res.StatusCode,
			Message: errorFormat,
		}
	}
	if dest == nil {
		return nil
	}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return fmt.Errorf(errorFormat, err)
	}
	return nil
}
