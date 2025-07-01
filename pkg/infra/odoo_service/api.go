package odoo_service

import (
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/odoo_service/repository"
	"fmt"
	"github.com/kolo/xmlrpc"
	"net/http"
)

var _ repository.OdooServiceRepository = (*Client)(nil)

type Client struct {
	rpcClient *xmlrpc.Client
	config    config.OdooServiceConfig
}

func NewClient(config config.OdooServiceConfig, transport http.RoundTripper) (*Client, error) {
	client, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", config.Url), transport)
	if err != nil {
		return &Client{}, fmt.Errorf("OdooServiceRepository NewClient %w", err)
	}
	return &Client{
		rpcClient: client,
		config:    config,
	}, nil
}

func (c *Client) SendLoanApprovalRequest(loanApprovalRequest entity.LoanApprovalRequest) error {
	errorTemplate := "OdooServiceRepository SendApprovalLoanRequest %w"
	args := []any{
		c.config.Db, c.config.Uid, c.config.Password,
		"approval.finx", "create",
		loanApprovalRequest.ToOdooFormat(),
	}
	err := c.rpcClient.Call("execute_kw", args, nil)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}
