package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ShauryaAg/ProductAPI/models"
	"github.com/ShauryaAg/ProductAPI/models/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product *models.Product = &models.Product{}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if !strings.Contains(ct, "application/json") {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type: 'application/json', but got %s", ct)))
		return
	}

	err = json.Unmarshal(bodyBytes, product)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	product, err = models.NewProduct(product.Name, product.Description, product.ThumbnailImageUrl)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := db.Models["product"].InsertOne(r.Context(), product)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jsonBytes, err := json.Marshal(struct {
		Id primitive.ObjectID
	}{result.InsertedID.(primitive.ObjectID)})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	var products []models.Product

	cursor, err := db.Models["product"].Find(
		r.Context(),
		bson.M{"name": bson.M{"$regex": r.URL.Query().Get("q")}},
	)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer cursor.Close(r.Context())

	err = cursor.All(r.Context(), &products)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
