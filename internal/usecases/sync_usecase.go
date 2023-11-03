package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
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

func (rU *syncUsecase) InsertOne(ctx context.Context, sync domains.Sync) error {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	return rU.syncRepository.InsertOne(ctx, sync)
}

func (rU *syncUsecase) GetLast(ctx context.Context, sync *domains.Sync) error {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	return rU.syncRepository.GetLast(ctx, sync)
}
