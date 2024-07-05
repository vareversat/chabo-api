package domains

import (
	"context"
	"os"
	"time"

	"github.com/vareversat/chabo-api/internal/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	SyncCollection = os.Getenv("MONGO_SYNCS_COLLECTION_NAME")
)

type Sync struct {
	ID        primitive.ObjectID `json:"-"           bson:"_id,omitempty"`
	ItemCount int                `json:"item_count"  bson:"item_count"    example:"10"`
	Duration  time.Duration      `json:"duration_ms" bson:"duration_ms"   example:"130"                         swaggertype:"primitive,integer"`
	Timestamp time.Time          `json:"timestamp"   bson:"timestamp"     example:"2021-05-25T00:53:16.535668Z"                                 format:"date-time"`
}

type SyncRepository interface {
	InsertOne(ctx context.Context, sync Sync) error
	GetLast(ctx context.Context, sync *Sync) error
}

type SyncUseCase interface {
	InsertOne(ctx context.Context, sync Sync) errors.CustomError
	GetLast(ctx context.Context, sync *Sync) errors.CustomError
}
