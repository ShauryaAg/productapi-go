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

// Migrate function will create the collections if they do not exist.
// Along with creating the collections, it will also create indexes on the collections.
func Migrate(ctx context.Context, db *mongo.Database, models ...interface{}) (map[string]*mongo.Collection, error) {
	collections := make(map[string]*mongo.Collection)
	for _, model := range models {
		modelType := reflect.TypeOf(model)
		modelName := strings.ToLower(modelType.Name())

		exists, err := collectionExists(ctx, db, modelName)
		if err != nil {
			return nil, err
		}

		// If the collection does not exist, create it.
		if !exists {
			err := db.CreateCollection(ctx, modelName)
			if err != nil {
				return nil, err
			}
		}

		// Get the collection.
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

			// Check if the index already exists
			indexes, err := collection.Indexes().ListSpecifications(ctx)
			if err != nil {
				return nil, err
			}

			// regex to check if index such as "fieldName_1" exists
			re := regexp.MustCompile(fmt.Sprintf(`^%s_\d+`, fieldName))
			for _, index := range indexes {
				if re.Match([]byte(index.Name)) {
					continue loop
				}
				// Check if existing index is unique
				if (index.Unique != nil) && (strings.Contains(fieldTag, "unique")) {
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

			// Create the index with the options
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

// collectionExists checks if a collection exists in the database.
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
