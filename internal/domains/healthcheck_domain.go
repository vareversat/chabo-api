package domains

import "context"

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
	GetHealth(ctx context.Context) error
}
