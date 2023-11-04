package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
	"github.com/vareversat/chabo-api/internal/errors"
)

type syncUsecase struct {
	syncRepository domains.SyncRepository
	contextTimeout time.Duration
}

func NewSyncUsecase(
	syncRepository domains.SyncRepository,
	timeout time.Duration,
) domains.SyncUsecase {
	return &syncUsecase{
		syncRepository: syncRepository,
		contextTimeout: timeout,
	}
}

func (rU *syncUsecase) InsertOne(ctx context.Context, sync domains.Sync) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	err := rU.syncRepository.InsertOne(ctx, sync)

	if err != nil {
		return errors.NewNotFoundError("No sync status exists in database")
	}
	return nil
}

func (rU *syncUsecase) GetLast(ctx context.Context, sync *domains.Sync) errors.CustomError {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	err := rU.syncRepository.GetLast(ctx, sync)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
