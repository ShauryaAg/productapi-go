package db

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/ShauryaAg/ProductAPI/models"
	"go.mongodb.org/mongo-driver/bson"
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

	err = migrate(ctx, client.Database(database), models.User{}, models.Product{}, models.Review{})
	if err != nil {
		return nil, err
	}

	Models = getCollections(*client.Database(database), "user", "product", "review")
	return client, nil
}

func getCollections(db mongo.Database, collections ...string) map[string]*mongo.Collection {
	models := make(map[string]*mongo.Collection)
	for _, collection := range collections {
		models[collection] = db.Collection(collection)
	}

	return models
}

func migrate(ctx context.Context, db *mongo.Database, models ...interface{}) error {
	for _, model := range models {
		modelType := reflect.TypeOf(model)
		modelName := strings.ToLower(modelType.Name())

		exists, err := collectionExists(db, modelName)
		if err != nil {
			return err
		}

		if !exists {
			err := db.CreateCollection(ctx, modelName)
			if err != nil {
				return err
			}
		}

		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			fieldOptions := strings.Split(field.Tag.Get("mongo"), ",")

		loop:
			for _, option := range fieldOptions {
				switch option {
				case "unique":
					_, err := db.Collection(modelName).Indexes().CreateOne(ctx, mongo.IndexModel{
						Keys:    bson.M{fmt.Sprintf("%s_unique_index", field.Name): 1},
						Options: options.Index().SetUnique(true),
					})
					if err != nil {
						return err
					}
					break loop
				case "index":
					_, err := db.Collection(modelName).Indexes().CreateOne(ctx, mongo.IndexModel{
						Keys:    bson.M{fmt.Sprintf("%s_index", field.Name): 1},
						Options: options.Index(),
					})
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func collectionExists(db *mongo.Database, collection string) (bool, error) {
	names, err := db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		return false, err
	}

	for _, name := range names {
		if name == collection {
			return true, nil
		}
	}

	return false, nil
}
