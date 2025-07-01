package jwttoken

import (
	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	jwt.RegisteredClaims

	FullName       string   `json:"fullName,omitempty"`
	UserId         string   `json:"userId,omitempty"`
	CustomerEmail  string   `json:"customerEmail,omitempty"`
	CustomerMobile string   `json:"customerMobile,omitempty"`
	Username       string   `json:"username,omitempty"`
	Roles          []string `json:"roles,omitempty"`
	Status         string   `json:"status"`
	InvestorId     string   `json:"investorId"`
	Sub            string   `json:"sub"`
	CustodyCode    string   `json:"custodyCode"`
}
