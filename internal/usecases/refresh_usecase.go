package usecases

import (
	"context"
	"time"

	"github.com/vareversat/chabo-api/internal/domains"
)

type refreshUsecase struct {
	refreshRepository domains.RefreshRepository
	contextTimeout    time.Duration
}

func NewRefreshUsecase(
	refreshRepository domains.RefreshRepository,
	timeout time.Duration,
) domains.RefreshUsecase {
	return &refreshUsecase{
		refreshRepository: refreshRepository,
		contextTimeout:    timeout,
	}
}

func (rU *refreshUsecase) InsertOne(ctx context.Context, refresh domains.Refresh) error {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	return rU.refreshRepository.InsertOne(ctx, refresh)
}

func (rU *refreshUsecase) GetLast(ctx context.Context, refresh *domains.Refresh) error {
	ctx, cancel := context.WithTimeout(ctx, rU.contextTimeout)
	defer cancel()
	return rU.refreshRepository.GetLast(ctx, refresh)
}
