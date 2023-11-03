package domains

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	SyncCollection = os.Getenv("MONGO_SYNCS_COLLECTION_NAME")
)

type Sync struct {
	ID        primitive.ObjectID `json:"-"           bson:"_id,omitempty"`
	ItemCount int                `json:"item_count"  bson:"item_count"    example:"10"`
	Duration  time.Duration      `json:"duration_ns" bson:"duration_ns"   example:"348872934"                   swaggertype:"primitive,integer"`
	Timestamp time.Time          `json:"timestamp"   bson:"timestamp"     example:"2021-05-25T00:53:16.535668Z"                                 format:"date-time"`
}

type SyncRepository interface {
	InsertOne(ctx context.Context, sync Sync) error
	GetLast(ctx context.Context, sync *Sync) error
}

type SyncUsecase interface {
	InsertOne(ctx context.Context, sync Sync) error
	GetLast(ctx context.Context, sync *Sync) error
}
