package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
)

type healthcheckUsecase struct {
	healthcheckRepository domains.HealthCheckRepository
	contextTimeout        time.Duration
}

func NewHealthCheckUsecase(
	healthcheckRepository domains.HealthCheckRepository,
	timeout time.Duration,
) domains.HealthCheckUsecase {
	return &healthcheckUsecase{
		healthcheckRepository: healthcheckRepository,
		contextTimeout:        timeout,
	}
}

func (rU *healthcheckUsecase) GetHealth(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	return rU.healthcheckRepository.GetDBHealth(ctx)
}
