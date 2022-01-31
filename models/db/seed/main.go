package seed

// func to seed data to mongo database

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"
	"strings"

	"github.com/ShauryaAg/ProductAPI/models"
	"github.com/ShauryaAg/ProductAPI/models/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// Seed is a wrapper function to seed data to mongo database
func Seed() {
	ctx := context.Background()
	client, err := db.InitDatabase("mongo", ctx)
	if err != nil {
		log.Fatal(err)
	}

	dataFiles := map[string]interface{}{
		"models/db/seed/user_data.json":    models.User{},
		"models/db/seed/product_data.json": models.Product{},
		"models/db/seed/review_data.json":  models.Review{},
	}

	err = seedDB(ctx, client, dataFiles)
	if err != nil {
		log.Fatal(err)
	}
}

// seedDB seeds data to mongo database
func seedDB(ctx context.Context, client *mongo.Client, dataFiles map[string]interface{}) error {
	db := client.Database("mongo")
	for file, model := range dataFiles {
		err := seedModel(ctx, db, model, file)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedModel seeds data to each mongo collection
func seedModel(ctx context.Context, db *mongo.Database, model interface{}, file string) error {
	modelName := strings.ToLower(reflect.TypeOf(model).Name())
	collection := db.Collection(modelName)

	jsonFile, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	var obj []interface{}
	err = json.Unmarshal(jsonFile, &obj)
	if err != nil {
		return err
	}

	collection.InsertMany(ctx, obj)
	return nil
}
