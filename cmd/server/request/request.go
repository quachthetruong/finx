package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIdKey = "requestId"

func AddRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId, _ := uuid.NewRandom()
		c.Set(RequestIdKey, requestId.String())
	}
}
