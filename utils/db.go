package utils

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
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

		exists, err := collectionExists(ctx, db, modelName)
		if err != nil {
			return nil, err
		}

		if !exists {
			err := db.CreateCollection(ctx, modelName)
			if err != nil {
				return nil, err
			}
		}

		collection := db.Collection(modelName)
		collections[modelName] = collection

	loop:
		for i := 0; i < modelType.NumField(); i++ {
			field := modelType.Field(i)
			fieldName := strings.ToLower(field.Name)
			fieldTag := field.Tag.Get("mongo")
			fieldOptions := strings.Split(fieldTag, ",")

			if len(fieldOptions) == 0 {
				continue loop
			}

			// check if index exists already
			indexes, err := collection.Indexes().ListSpecifications(ctx)
			if err != nil {
				return nil, err
			}

			re := regexp.MustCompile(fmt.Sprintf(`^%s_\d+`, fieldName))
			for _, index := range indexes {
				if re.Match([]byte(index.Name)) {
					continue loop
				}
				var unique = true
				if index.Unique == nil {
					unique = false
				}
				if (unique) && (strings.Contains(fieldTag, "unique")) {
					continue loop
				}
			}

			indexOptions := options.Index()
			for _, option := range fieldOptions {
				switch option {
				case "unique":
					indexOptions = indexOptions.SetUnique(true)
				case "sparse":
					indexOptions = indexOptions.SetSparse(true)
					// TODO: add more options for max,min,expire
				}
			}
			_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.M{fieldName: 1},
				Options: indexOptions,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	return collections, nil
}

func collectionExists(ctx context.Context, db *mongo.Database, collection string) (bool, error) {
	names, err := db.ListCollectionNames(ctx, bson.M{})
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
