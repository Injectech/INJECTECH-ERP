package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"backend/internal/transport/http/handler"
	"backend/internal/transport/http/middleware"
	usecaseauth "backend/internal/usecase/auth"
)

// Build constructs the gin engine with middleware and routes.
func Build(
	log *zap.Logger,
	authUC *usecaseauth.Service,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	roleHandler *handler.RoleHandler,
	permissionHandler *handler.PermissionHandler,
	inventoryHandler *handler.InventoryHandler,
	auditHandler *handler.AuditHandler,
	locationHandler *handler.LocationHandler,
	corsOrigins []string,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger(log))
	r.Use(middleware.CORS(corsOrigins))

	// Swagger UI (served from docs/swagger.json)
	r.StaticFile("/swagger.json", "docs/swagger.json")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger.json")))

	// Preflight handler for CORS
	r.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok"})
	})

	api := r.Group("/api/v1")
	authHandler.RegisterRoutes(api.Group("/auth"))

	secured := api.Group("")
	secured.Use(middleware.Auth(authUC))
	userHandler.RegisterRoutes(secured)
	productHandler.RegisterRoutes(secured)
	roleHandler.RegisterRoutes(secured)
	permissionHandler.RegisterRoutes(secured)
	inventoryHandler.RegisterRoutes(secured)
	auditHandler.RegisterRoutes(secured)
	locationHandler.RegisterRoutes(secured)

	return r
}

func requestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		log.Info("http request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
		)
	}
}
