package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/errors"
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

func (rU *healthcheckUsecase) GetHealth(ctx context.Context) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()

	err := rU.healthcheckRepository.GetDBHealth(ctx)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
