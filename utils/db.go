package utils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Migrate(ctx context.Context, db *mongo.Database, models ...interface{}) (map[string]*mongo.Collection, error) {
	collections := make(map[string]*mongo.Collection)
	for _, model := range models {
		modelType := reflect.TypeOf(model)
		modelName := strings.ToLower(modelType.Name())

		exists, err := collectionExists(db, modelName)
		if err != nil {
			return nil, err
		}

		if !exists {
			err := db.CreateCollection(ctx, modelName)
			if err != nil {
				return nil, err
			}
		}

		collections[modelName] = db.Collection(modelName)

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
						return nil, err
					}
					break loop
				case "index":
					_, err := db.Collection(modelName).Indexes().CreateOne(ctx, mongo.IndexModel{
						Keys:    bson.M{fmt.Sprintf("%s_index", field.Name): 1},
						Options: options.Index(),
					})
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return collections, nil
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
