package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ShauryaAg/ProductAPI/models"
	"github.com/ShauryaAg/ProductAPI/models/db"
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateReview(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("decoded")

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

	var review models.Review
	err = json.Unmarshal(bodyBytes, &review)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var user models.User
	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = db.Models["user"].FindOne(r.Context(), bson.M{"_id": userObjectId}).Decode(&user)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	review = *models.NewReview(review.Text, review.Rating, user)
	result, err := db.Models["review"].InsertOne(r.Context(), review)
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
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	review.Id = result.InsertedID.(primitive.ObjectID)
	productId := chi.URLParam(r, "productId")
	productObjectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res := db.Models["product"].FindOneAndUpdate(r.Context(),
		bson.M{"_id": productObjectId},
		bson.M{"$push": bson.M{"reviews": review}},
	)
	if res.Err() != nil {
		fmt.Println("err", res.Err().Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(res.Err().Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}
