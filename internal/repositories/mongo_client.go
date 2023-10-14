package repositories

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDSN = os.Getenv("MONGO_DSN")

func NewMongoClient() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().
		ApplyURI(mongoDSN).
		SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	} else {
		// Send a ping command to test the connection
		if err := client.Ping(context.TODO(), nil); err != nil {
			panic(err)
		}
	}

	return client
}
