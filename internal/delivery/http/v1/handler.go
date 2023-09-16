package v1

import (
	"fio_service/internal/service"
	"fio_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Services
	logger  logger.Logger
}

func NewHandler(service *service.Services, logger *logger.Logger) *Handler {
	return &Handler{
		service: *service,
		logger:  *logger,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		go h.consumeMessages()
		h.initPersonRoutes(v1)

	}
}
