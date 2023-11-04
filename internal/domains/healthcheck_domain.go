package domains

import (
	"context"

	"github.com/vareversat/chabo-api/internal/errors"
)

type SystemHealthOK struct {
	Message string `json:"message" example:"system is running properly"`
}

type SystemHealthNOK struct {
	Error string `json:"error" example:"system is not running properly"`
}

type HealthCheckRepository interface {
	GetDBHealth(ctx context.Context) error
}

type HealthCheckUsecase interface {
	GetHealth(ctx context.Context) errors.CustomError
}
