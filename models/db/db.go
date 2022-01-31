package db

import (
	"context"
	"os"

	"github.com/ShauryaAg/ProductAPI/models"
	"github.com/ShauryaAg/ProductAPI/utils"
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

// InitDatabase initializes the database connection
// and returns the intialized client object
func InitDatabase(database string, ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	Models, err = utils.Migrate(ctx, client.Database(database), models.User{}, models.Product{}, models.Review{})
	if err != nil {
		return nil, err
	}

	return client, nil
}
