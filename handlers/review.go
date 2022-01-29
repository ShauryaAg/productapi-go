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
	var review *models.Review = &models.Review{}
	var user models.User

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
	err = json.Unmarshal(bodyBytes, review)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

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

	review, err = models.NewReview(review.Text, review.Rating, user)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	reviewResult, err := db.Models["review"].InsertOne(r.Context(), review)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jsonBytes, err := json.Marshal(struct {
		Id primitive.ObjectID
	}{reviewResult.InsertedID.(primitive.ObjectID)})
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	review.Id = reviewResult.InsertedID.(primitive.ObjectID)
	productId := chi.URLParam(r, "productId")
	productObjectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		fmt.Println("err", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	productResult := db.Models["product"].FindOneAndUpdate(r.Context(),
		bson.M{"_id": productObjectId},
		bson.M{
			"$push": bson.M{"reviews": review},
			"$inc":  bson.M{"ratingCount": 1, "ratingSum": review.Rating},
		},
	)
	if productResult.Err() != nil {
		fmt.Println("err", productResult.Err().Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(productResult.Err().Error()))
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}
