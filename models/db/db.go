package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DBCon  *mongo.Database
	Models map[string]*mongo.Collection
)

var (
	connectionUri = os.Getenv("MONGODB_URI")
)

func InitDatabase(database string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	Models = getCollections(*client.Database(database), "user")
	return client, nil
}

func getCollections(db mongo.Database, collections ...string) map[string]*mongo.Collection {
	models := make(map[string]*mongo.Collection)
	for _, collection := range collections {
		models[collection] = db.Collection(collection)
	}

	return models
}
