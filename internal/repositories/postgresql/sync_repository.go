package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vareversat/chabo-api/internal/domains"
)

type syncRepository struct {
	connectionPool *pgxpool.Pool
}

func NewSyncRepository(connectionPool *pgxpool.Pool) domains.SyncRepository {
	return &syncRepository{
		connectionPool: connectionPool,
	}
}

func (rR *syncRepository) InsertOne(ctx context.Context, sync domains.Sync) error {
	query := `
        INSERT INTO syncs (item_count, duration, created_at)
        VALUES (@item_count, @duration, @created_at)
        RETURNING sync_id
    `
	args := pgx.NamedArgs{
		"item_count": sync.ItemCount,
		"duration":   sync.Duration,
		"created_at": sync.CreatedAt,
	}
	var id string
	return rR.connectionPool.QueryRow(ctx, query, args).Scan(&id)
}

func (rR *syncRepository) GetLast(ctx context.Context, sync *domains.Sync) error {
	query := `
		SELECT * FROM syncs
		ORDER BY created_at DESC
		LIMIT 1
	`
	return rR.connectionPool.QueryRow(ctx, query).Scan(&sync.ID, &sync.ItemCount, &sync.Duration, &sync.CreatedAt)
}
