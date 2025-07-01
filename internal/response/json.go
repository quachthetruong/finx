package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONWithHeaders(c *gin.Context, status int, data any, headers http.Header) {
	for key, value := range headers {
		c.Writer.Header()[key] = value
	}
	c.JSON(status, data)
}
