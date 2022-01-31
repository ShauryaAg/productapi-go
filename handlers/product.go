package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ShauryaAg/ProductAPI/models"
	"github.com/ShauryaAg/ProductAPI/models/db"
	"github.com/ShauryaAg/ProductAPI/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product *models.Product = &models.Product{}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	ct := r.Header.Get("content-type")
	if !strings.Contains(ct, "application/json") {
		utils.Error(
			w, r,
			fmt.Sprintf("Need content-type: 'application/json', but got %s", ct),
			http.StatusUnsupportedMediaType,
		)
		return
	}

	err = json.Unmarshal(bodyBytes, product)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	product, err = models.NewProduct(product.Name, product.Description, product.ThumbnailImageUrl)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Models["product"].InsertOne(r.Context(), product)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(struct {
		Id primitive.ObjectID
	}{result.InsertedID.(primitive.ObjectID)})
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	filter := bson.M{}
	if q := r.URL.Query().Get("q"); q != "" {
		filter = bson.M{"name": bson.M{"$regex": q}}
	}

	paginationOptions, err := utils.Pagination(r, options.Find())
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	cursor, err := db.Models["product"].Find(
		r.Context(),
		filter,
		paginationOptions,
	)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(r.Context())

	err = cursor.All(r.Context(), &products)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	type productResult struct {
		Id                primitive.ObjectID
		Name              string
		Description       string
		ThumbnailImageUrl string
		Rating            float64
	}

	var productResults []productResult
	for _, product := range products {
		productResults = append(productResults, productResult{
			Id:                product.Id,
			Name:              product.Name,
			Description:       product.Description,
			ThumbnailImageUrl: product.ThumbnailImageUrl,
			Rating:            product.Rating(),
		})
	}

	jsonBytes, err := json.Marshal(productResults)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
