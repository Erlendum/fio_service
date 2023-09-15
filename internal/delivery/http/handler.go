package http

import (
	v1 "fio_service/internal/delivery/http/v1"
	"fio_service/internal/service"
	"fio_service/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services service.Services
	logger   logger.Logger
}

func NewHandler(services *service.Services, logger *logger.Logger) *Handler {
	return &Handler{
		services: *services,
		logger:   *logger,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	h.initAPI(router)
	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(&h.services, &h.logger)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
