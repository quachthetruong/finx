package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	internalroutes "financing-offer/cmd/server/api/internal-routes"
	v1 "financing-offer/cmd/server/api/v1"
	"financing-offer/cmd/server/request"
	"financing-offer/internal/config"
	"financing-offer/pkg/docs"
)

// NewRoutes godoc
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func (app *Application) routes() http.Handler {
	if app.Config.Env == config.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	middleware := app.Middleware
	docs.SwaggerInfo.Title = "Financing offer"
	docs.SwaggerInfo.Description = "Financing offer API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	r := gin.New()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	r.Use(
		cors.New(corsConfig),
		gin.Recovery(),
		request.AddRequestId(),
		middleware.JsonLoggerMiddleware(),
	)
	r.NoRoute(app.notFound)
	r.NoMethod(app.methodNotAllowed)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET(
		"/status", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)
	apiGroup := r.Group("/api", middleware.Authenticate())
	v1.NewRoutes(apiGroup, middleware, app.Injector)
	publicGroup := r.Group("/public")
	v1.NewPublicRoutes(publicGroup, middleware, app.Injector)
	internalroutes.NewRoutes(apiGroup, middleware, app.Injector)
	return r
}
