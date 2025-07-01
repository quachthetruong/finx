package odoo_service

import (
	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_GetMarginPoolsByIds(t *testing.T) {
	defer gock.Off()
	odooServiceConfig := config.OdooServiceConfig{
		Url: "http://dnse-odoo-service",
	}
	client, _ := NewClient(odooServiceConfig, gock.NewTransport())
	t.Run(
		"send approval request2", func(t *testing.T) {

			gock.New(odooServiceConfig.Url).
				Post("/xmlrpc/2/object").
				Reply(200).JSON(nil)
			err := client.SendLoanApprovalRequest(entity.LoanApprovalRequest{})
			assert.Nil(t, err)
		},
	)
}
