package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/errors"
)

type healthcheckUseCase struct {
	healthcheckRepository domains.HealthCheckRepository
	contextTimeout        time.Duration
}

func NewHealthCheckUseCase(
	healthcheckRepository domains.HealthCheckRepository,
	timeout time.Duration,
) domains.HealthCheckUsecase {
	return &healthcheckUseCase{
		healthcheckRepository: healthcheckRepository,
		contextTimeout:        timeout,
	}
}

func (rU *healthcheckUseCase) GetHealth(ctx context.Context) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()

	err := rU.healthcheckRepository.GetDBHealth(ctx)

	if err != nil {
		return errors.NewInternalServerError()
	}
	return nil
}
