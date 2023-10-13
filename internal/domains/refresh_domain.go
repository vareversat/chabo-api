package domains

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	RefreshCollection = os.Getenv("MONGO_REFRESHES_COLLECTION_NAME")
)

type Refresh struct {
	ID        primitive.ObjectID `json:"-"           bson:"_id,omitempty"`
	ItemCount int                `json:"item_count"  bson:"item_count"    example:"10"`
	Duration  time.Duration      `json:"duration_ns" bson:"duration_ns"   example:"348872934"                   swaggertype:"primitive,integer"`
	Timestamp time.Time          `json:"timestamp"   bson:"timestamp"     example:"2021-05-25T00:53:16.535668Z"                                 format:"date-time"`
}

type RefreshRepository interface {
	InsertOne(ctx context.Context, refresh Refresh) error
	GetLast(ctx context.Context, refresh *Refresh) error
}

type RefreshUsecase interface {
	InsertOne(ctx context.Context, refresh Refresh) error
	GetLast(ctx context.Context, refresh *Refresh) error
}
