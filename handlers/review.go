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
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create a new review
// POST /api/review/{productId}
func CreateReview(w http.ResponseWriter, r *http.Request) {
	var review *models.Review = &models.Review{}
	var user models.User

	userId := r.Header.Get("decoded")

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
	err = json.Unmarshal(bodyBytes, review)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.Models["user"].FindOne(r.Context(), bson.M{"_id": userObjectId}).Decode(&user)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	review, err = models.NewReview(review.Text, review.Rating, user)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	reviewResult, err := db.Models["review"].InsertOne(r.Context(), review)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(struct {
		Id primitive.ObjectID
	}{reviewResult.InsertedID.(primitive.ObjectID)})
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	review.Id = reviewResult.InsertedID.(primitive.ObjectID)
	productId := chi.URLParam(r, "productId")
	productObjectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
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
		utils.Error(w, r, productResult.Err().Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}
