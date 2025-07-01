package appcontext

import (
	"context"

	"financing-offer/internal/jwttoken"
)

const (
	UserInformation = "UserInformation"
)

func ContextGetCustomerInfo(ctx context.Context) *jwttoken.AdminClaims {
	if user, ok := ctx.Value(UserInformation).(*jwttoken.AdminClaims); ok {
		return user
	}
	return nil
}

func ContextGetUserName(ctx context.Context) string {
	if user := ContextGetCustomerInfo(ctx); user != nil {
		return user.Username
	}
	return ""
}
