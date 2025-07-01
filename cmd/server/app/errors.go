package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/response"
)

func (app *Application) errorMessage(c *gin.Context, status int, message string, headers http.Header) {
	response.JSONWithHeaders(c, status, gin.H{"ErrorTrace": message}, headers)
}

func (app *Application) notFound(c *gin.Context) {
	message := "The requested resource could not be found"
	app.errorMessage(c, http.StatusNotFound, message, nil)
}

func (app *Application) methodNotAllowed(c *gin.Context) {
	message := fmt.Sprintf("The %s method is not supported for this resource", c.Request.Method)
	app.errorMessage(c, http.StatusMethodNotAllowed, message, nil)
}
