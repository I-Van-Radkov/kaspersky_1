package http

import (
	"github.com/I-Van-Radkov/kaspersky_1/internal/http/handlers"
	middlwares "github.com/I-Van-Radkov/kaspersky_1/internal/http/middlewares"
	"github.com/gin-gonic/gin"
)

func NewRouterGin(
	enqueueHandlers *handlers.EnqueueHandlers,
) *gin.Engine {
	router := gin.Default()

	router.Use(middlwares.CorsMiddleware())

	router.POST("/enqueue", enqueueHandlers.EnqueueHandler)

	router.GET("/healthz", handlers.HealthzHandler)

	return router
}
