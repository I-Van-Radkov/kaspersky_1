package handlers

import (
	"net/http"

	"github.com/I-Van-Radkov/kaspersky_1/internal/dto"
	"github.com/I-Van-Radkov/kaspersky_1/internal/utils"
	"github.com/gin-gonic/gin"
)

type QueueServiceProvider interface {
	AddToQueue(id, payload string, maxRetries int)
}

type EnqueueHandlers struct {
	queueService QueueServiceProvider
}

func NewEnqueueHandlers(queueService QueueServiceProvider) *EnqueueHandlers {
	return &EnqueueHandlers{
		queueService: queueService,
	}
}

func (h *EnqueueHandlers) EnqueueHandler(c *gin.Context) {
	req, err := dto.ToEnqueueRequest(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}
	defer c.Request.Body.Close()

	if !utils.ValidateParams(req.Id, req.Payload) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id, payload and max_retries are required",
		})
		return
	}

	if !utils.ValidateMaxRetries(req.MaxRetries) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "max_retries must be greater than 0",
		})
		return
	}

	h.queueService.AddToQueue(req.Id, req.Payload, req.MaxRetries)
}
