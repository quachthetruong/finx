package financialproduct

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/financialproduct/repository"
	"financing-offer/internal/funcs"
)

var _ repository.FinancialProductRepository = (*Client)(nil)

const maxFetchSize = 500

type Client struct {
	httpClient *http.Client
	config     config.FinancialProductConfig
}

func (c *Client) GetMarginBasketsByIds(ctx context.Context, ids []int64) ([]entity.MarginBasket, error) {
	var (
		errorGroup errgroup.Group
		mu         sync.Mutex
		baskets    = make([]entity.MarginBasket, 0, len(ids))
	)
	for _, id := range ids {
		errorGroup.Go(
			func() error {
				basket, err := c.GetMarginBasketDetail(ctx, id)
				if err != nil {
					return err
				}
				mu.Lock()
				baskets = append(baskets, basket)
				mu.Unlock()
				return nil
			},
		)
	}
	if err := errorGroup.Wait(); err != nil {
		return nil, fmt.Errorf("GetMarginBasketsByIds %w", err)
	}
	return baskets, nil
}

func (c *Client) GetMarginBasketDetail(ctx context.Context, id int64) (entity.MarginBasket, error) {
	errorFormat := "GetMarginBasketDetail %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/v2/margin-baskets/%d/detail", c.config.Url, id), nil)
	if err != nil {
		return entity.MarginBasket{}, fmt.Errorf(errorFormat, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.MarginBasket{}, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	dest := entity.MarginBasket{}
	if err := c.handleResponse(res, errorFormat, &dest); err != nil {
		return entity.MarginBasket{}, err
	}
	return dest, nil
}

func (c *Client) GetLoanRatesByIds(ctx context.Context, loanRateIds []int64) ([]entity.LoanRate, error) {
	errorFormat := "GetLoanRatesByIds %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/loan-rates", c.config.Url), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	q := req.URL.Query()
	q.Add("ids", strings.Join(funcs.Map(loanRateIds, func(id int64) string { return strconv.FormatInt(id, 10) }), ","))
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
	dest := ListResponse[entity.LoanRate]{}
	if err := c.handleResponse(res, errorFormat, &dest); err != nil {
		return nil, err
	}
	return dest.Data, nil
}

func (c *Client) AssignLoanPackageOrGetLoanPackageAccountId(ctx context.Context, accountNo string, loanId int64, assetType entity.AssetType) (loanPackageAccountId int64, err error) {
	loanPackageAccountId, err = c.AssignLoanPackage(ctx, accountNo, loanId, assetType)
	if err != nil {
		if errors.Is(err, apperrors.ErrLoanPackageAccountAlreadyExisted(accountNo, loanId)) {
			existedLoanPackageAccountId, scopedErr := c.GetLoanPackageAccountIdByAccountNoAndLoanPackageId(
				ctx, accountNo, loanId, assetType,
			)
			if scopedErr != nil {
				return 0, fmt.Errorf("AssignLoanPackageOrGetLoanPackageAccountId %w", scopedErr)
			}
			return existedLoanPackageAccountId, nil
		} else if err != nil {
			return 0, fmt.Errorf("AssignLoanPackageOrGetLoanPackageAccountId %w", err)
		}
	}
	return loanPackageAccountId, nil
}

func NewClient(config config.FinancialProductConfig) *Client {
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

func (c *Client) GetAllAccountDetail(ctx context.Context, investorId string) (accounts []entity.FinancialAccountDetail, err error) {
	errorFormat := "GetAllAccountDetail %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/accounts", c.config.Url), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	q := req.URL.Query()
	q.Add("investorId", investorId)
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
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("GetAllAccountDetail %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return nil, fmt.Errorf("GetAllAccountDetail %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := FinancialAccountsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	return dest.Accounts, nil
}

func (c *Client) GetAllAccountDetailByCustodyCode(ctx context.Context, custodyCode string) (accounts []entity.FinancialAccountDetail, err error) {
	errorFormat := "GetAllAccountDetail %w"
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("%s/accounts", c.config.Url), nil)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	q := req.URL.Query()
	q.Add("custodyCode", custodyCode)
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
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("GetAllAccountDetailByCustodyCode %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return nil, fmt.Errorf("GetAllAccountDetailByCustodyCode %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := FinancialAccountsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}
	return dest.Accounts, nil
}

func (c *Client) AssignLoanPackage(ctx context.Context, accountNo string, loanId int64, assetType entity.AssetType) (int64, error) {
	errorFormat := "AssignLoanPackage %w"
	url := fmt.Sprintf("%s/v2/loan-package-accounts", c.config.Url)
	if assetType == entity.AssetTypeDerivative {
		url = fmt.Sprintf("%s/derivatives/package-accounts", c.config.Url)
	}
	req, err := c.NewRequest(
		ctx,
		http.MethodPost,
		url,
		AssignLoanIdAccountNoRequest{
			LoanPackageId: strconv.FormatInt(loanId, 10),
			AccountNo:     accountNo,
		},
	)
	if err != nil {
		return 0, fmt.Errorf(errorFormat, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		errorResponse := ErrorResponse{}
		_ = json.NewDecoder(res.Body).Decode(&errorResponse)
		if strings.Contains(errorResponse.Message, "already existed") || strings.Contains(
			errorResponse.Message, "already exists",
		) {
			return 0, apperrors.ErrLoanPackageAccountAlreadyExisted(accountNo, loanId)
		}
		if strings.Contains(errorResponse.Message, "not existed") {
			return 0, apperrors.ErrLoanPackageAccountNotExisted(loanId)
		}
		return 0, fmt.Errorf("AssignLoanPackage %s got error Status %d, Message: %s", req.URL.String(), res.StatusCode, errorResponse.Message)
	}
	dest := AssignLoanPackageResponse{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return 0, fmt.Errorf(errorFormat, err)
	}
	return dest.Id, nil
}

func (c *Client) GetLoanPackageAccountIdByAccountNoAndLoanPackageId(ctx context.Context, accountNo string, loanPackageId int64, assetType entity.AssetType) (int64, error) {
	errorFormat := "GetLoanPackageAccountIdByAccountNoAndLoanPackageId %w"
	url := fmt.Sprintf(
		"%s/v2/loan-package-accounts?accountNo=%s&loanPackageId=%d", c.config.Url, accountNo, loanPackageId,
	)
	if assetType == entity.AssetTypeDerivative {
		url = fmt.Sprintf(
			"%s/derivatives/package-accounts?accountNo=%s&loanPackageId=%d", c.config.Url, accountNo, loanPackageId,
		)
	}
	req, err := c.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf(errorFormat, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return 0, fmt.Errorf("GetLoanPackageAccountIdByAccountNoAndLoanPackageId %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return 0, fmt.Errorf("GetLoanPackageAccountIdByAccountNoAndLoanPackageId %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := struct {
		Data []struct {
			Id int64 `json:"id"`
		} `json:"data"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return 0, fmt.Errorf("GetLoanPackageAccountIdByAccountNoAndLoanPackageId %w", err)
	}
	if len(dest.Data) == 0 {
		return 0, fmt.Errorf("GetLoanPackageAccountIdByAccountNoAndLoanPackageId got empty response")
	}
	return dest.Data[0].Id, nil
}

func (c *Client) GetLoanPackageDetail(ctx context.Context, loanPackageId int64) (entity.FinancialProductLoanPackage, error) {
	errorTemplate := "GetLoanPackageDetail %w"
	url := fmt.Sprintf("%s/v2/loan-packages/%d", c.config.Url, loanPackageId)
	req, err := c.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entity.FinancialProductLoanPackage{}, fmt.Errorf(errorTemplate, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.FinancialProductLoanPackage{}, fmt.Errorf(errorTemplate, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorTemplate, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return entity.FinancialProductLoanPackage{}, fmt.Errorf("GetLoanPackageDetail %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return entity.FinancialProductLoanPackage{}, fmt.Errorf("GetLoanPackageDetail %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := entity.FinancialProductLoanPackage{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return entity.FinancialProductLoanPackage{}, fmt.Errorf(errorTemplate, err)
	}
	return dest, nil
}

func (c *Client) GetLoanPackageDerivative(ctx context.Context, loanPackageId int64) (entity.FinancialProductLoanPackageDerivative, error) {
	errorTemplate := "GetLoanPackageDerivative %w"
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/derivatives/packages/%d", c.config.Url, loanPackageId), nil,
	)
	if err != nil {
		return entity.FinancialProductLoanPackageDerivative{}, fmt.Errorf(errorTemplate, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.FinancialProductLoanPackageDerivative{}, fmt.Errorf(errorTemplate, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorTemplate, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return entity.FinancialProductLoanPackageDerivative{}, fmt.Errorf("GetLoanPackageDerivative %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return entity.FinancialProductLoanPackageDerivative{}, fmt.Errorf("GetLoanPackageDerivative %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := entity.FinancialProductLoanPackageDerivative{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return entity.FinancialProductLoanPackageDerivative{}, fmt.Errorf(errorTemplate, err)
	}
	return dest, nil
}

func (c *Client) GetLoanPackageDetails(ctx context.Context, loanPackageIds []int64) ([]entity.FinancialProductLoanPackage, error) {
	errorTemplate := "GetLoanPackageDetails %w"
	loanPackageIdsString := funcs.Map(
		loanPackageIds, func(id int64) string {
			return strconv.FormatInt(id, 10)
		},
	)
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v2/loan-packages?ids=%s", c.config.Url, strings.Join(loanPackageIdsString, ",")), nil,
	)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorTemplate, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("GetLoanPackageDetails %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return nil, fmt.Errorf("GetLoanPackageDetails %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := LoanPackageDetailsResponse{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	return dest.Data, nil
}

func (c *Client) GetLoanRateDetail(ctx context.Context, loanRateId int64) (entity.LoanRate, error) {
	errorTemplate := "GetLoanPackageDetail %w"
	url := fmt.Sprintf("%s/loan-rates/%d", c.config.Url, loanRateId)
	req, err := c.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return entity.LoanRate{}, fmt.Errorf(errorTemplate, err)
	}
	res, err := c.httpClient.Do(req)
	if err != nil {
		return entity.LoanRate{}, fmt.Errorf(errorTemplate, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorTemplate, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return entity.LoanRate{}, fmt.Errorf("GetLoanRateDetail %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return entity.LoanRate{}, fmt.Errorf("GetLoanRateDetail %s got error Status %d, Message %s", req.URL.String(), res.StatusCode, rawMessage)
	}
	dest := entity.LoanRate{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return entity.LoanRate{}, fmt.Errorf(errorTemplate, err)
	}
	return dest, nil
}

func (c *Client) GetLoanProducts(ctx context.Context, filter entity.MarginProductFilter) ([]entity.MarginProduct, error) {
	errorTemplate := "GetLoanProducts %w"
	req, err := c.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/loan-products", c.config.Url),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	q := req.URL.Query()
	q.Add("pageSize", strconv.Itoa(maxFetchSize))
	if filter.Symbol != "" {
		q.Add("symbol", filter.Symbol)
	}
	req.URL.RawQuery = q.Encode()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(errorTemplate, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorTemplate, scopedErr)
		}
	}()
	dest := ListResponse[entity.MarginProduct]{}
	if err := c.handleResponse(res, errorTemplate, &dest); err != nil {
		return nil, err
	}
	return dest.Data, nil
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
