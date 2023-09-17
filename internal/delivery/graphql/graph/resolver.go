package graph

import (
	"fio_service/internal/service"
	"fio_service/pkg/logger"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Services service.Services
	Logger   logger.Logger
}
