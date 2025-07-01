package odoo_service

import (
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"fmt"
	"github.com/kolo/xmlrpc"
)

type InMemoryClient struct {
	rpcClient *xmlrpc.Client
	config    config.OdooServiceConfig
}

func NewInMemoryClient(config config.OdooServiceConfig) (*InMemoryClient, error) {
	client, err := xmlrpc.NewClient(fmt.Sprintf("%s/xmlrpc/2/object", config.Url), nil)
	if err != nil {
		return &InMemoryClient{}, fmt.Errorf("OdooServiceRepository NewClient %w", err)
	}
	return &InMemoryClient{
		rpcClient: client,
		config:    config,
	}, nil
}

func (c *InMemoryClient) SendLoanApprovalRequest(loanApprovalRequest entity.LoanApprovalRequest) error {
	fmt.Printf("SendLoanApprovalRequest: %v\n", loanApprovalRequest)
	return nil
}
