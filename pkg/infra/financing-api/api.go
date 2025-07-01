package financing_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"financing-offer/internal/config"
)

type Client struct {
	httpClient *http.Client
	config     config.FinancingApiConfig
}

func NewClient(config config.FinancingApiConfig) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetDateAfter(date time.Time, workingDays int) (time.Time, error) {
	errorFormat := "GetDateAfter %w"
	req, err := http.NewRequest(
		http.MethodGet, fmt.Sprintf("%s/internal/api/business_dates/%s", c.config.Url, date.Format("2006-01-02")), nil,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf(errorFormat, err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.Token))
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("expandedDays", strconv.Itoa(workingDays))
	req.URL.RawQuery = q.Encode()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return time.Time{}, fmt.Errorf(errorFormat, err)
	}
	defer func() {
		if scopedErr := res.Body.Close(); scopedErr != nil {
			err = fmt.Errorf(errorFormat, scopedErr)
		}
	}()
	if res.StatusCode != http.StatusOK {
		rawMessage, err := io.ReadAll(res.Body)
		if err != nil {
			return time.Time{}, fmt.Errorf("GetDateAfter %s got error Status %d", req.URL.String(), res.StatusCode)
		}
		return time.Time{}, fmt.Errorf("GetDateAfter %s got error Status %d, message: %s", req.URL.String(), res.StatusCode, string(rawMessage))
	}
	dest := DateAfterResponse{}
	if err := json.NewDecoder(res.Body).Decode(&dest); err != nil {
		return time.Time{}, fmt.Errorf(errorFormat, err)
	}
	destTime, err := time.Parse("2006-01-02", dest.Data)
	if err != nil {
		return time.Time{}, fmt.Errorf(errorFormat, err)
	}
	return time.Date(
		destTime.Year(), destTime.Month(), destTime.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(),
		date.Location(),
	), nil
}
