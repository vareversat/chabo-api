package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/errors"
)

type boatUseCase struct {
	boatRepository domains.BoatRepository
	contextTimeout time.Duration
}

func NewBoatUseCase(
	boatRepository domains.BoatRepository,
	timeout time.Duration,
) domains.BoatUseCase {
	return &boatUseCase{
		boatRepository: boatRepository,
		contextTimeout: timeout,
	}
}

func (bU *boatUseCase) InsertOne(ctx context.Context, boat domains.Boat) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, bU.contextTimeout)
	defer cancel()
	err, _ := bU.boatRepository.InsertOne(ctx, boat)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
