package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"financing-offer/internal/appcontext"
	"financing-offer/internal/apperrors"
	"financing-offer/internal/jwttoken"
)

func (middleware *Middleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Authorization")
		authorizationHeader := c.GetHeader("Authorization")
		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token, err := middleware.parseClaims(headerParts[1])
				if err != nil {
					return
				}
				if claims, ok := token.Claims.(*jwttoken.AdminClaims); ok && token.Valid {
					c.Set(appcontext.UserInformation, claims)
				} else {
					return
				}
			}
		}
		c.Next()
	}
}

func (middleware *Middleware) parseClaims(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenString, &jwttoken.AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			key, err := jwttoken.GetPublicKey(middleware.Config.Jwt.PublicKey)
			if err != nil {
				return nil, err
			}
			return key, nil
		},
		jwt.WithValidMethods([]string{"RS256"}),
	)
}

func (middleware *Middleware) RequireAuthenticatedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticatedUser := appcontext.ContextGetCustomerInfo(c)
		if authenticatedUser == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

func (middleware *Middleware) RequireOneOfRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authenticatedUser := appcontext.ContextGetCustomerInfo(c)
		if authenticatedUser == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
			return
		}
		if !matchRole(authenticatedUser.Roles, allowedRoles) {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized, gin.H{"Error": "You are not allowed to perform this action"},
			)
			return
		}
		c.Next()
	}
}

func matchRole(userRoles []string, allowedRoles []string) bool {
	for _, userRole := range userRoles {
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(userRole, allowedRole) {
				return true
			}
		}
	}
	return false
}

func (middleware *Middleware) RequireFeatureEnable(featureName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		customerInfo := appcontext.ContextGetCustomerInfo(c)
		if customerInfo == nil || customerInfo.InvestorId == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "You're not allowed to perform this action"})
			return
		}
		if isEnable, err := middleware.FeatureFlagUseCase.IsFeatureEnable(
			featureName, customerInfo.InvestorId,
		); err != nil || !isEnable {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Error": "You're not allowed to perform this action"})
			return
		}
	}
}

func (middleware *Middleware) RequireHOActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		isHOActive, err := middleware.FlexRepo.IsHOActive(c)
		if err != nil {
			middleware.Logger.Error("Error checking HO status ", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error": "an error happened, please try again later"})
			return
		}
		if !isHOActive {
			c.AbortWithStatusJSON(http.StatusBadRequest, apperrors.ErrorHONotActive)
			return
		}
		c.Next()
	}
}
