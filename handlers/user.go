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
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user *models.User = &models.User{}

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

	err = json.Unmarshal(bodyBytes, user)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = models.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Models["user"].InsertOne(r.Context(), user)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := utils.CreateToken(*user)
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonBytes, err := json.Marshal(struct {
		Id    primitive.ObjectID
		Email string
		Token string
	}{result.InsertedID.(primitive.ObjectID), user.Email, token})
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	var user models.User

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
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.Models["user"].FindOne(r.Context(), bson.M{"email": data["email"]}).Decode(&user)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusNotFound)
		return
	}

	valid := user.VerifyPassword(data["password"])
	var token string
	if valid {
		token, err = utils.CreateToken(user)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		utils.Error(w, r, "Email/Password is incorrect", http.StatusBadRequest)
		return
	}

	jsonBytes, err := json.Marshal(struct {
		Id    primitive.ObjectID
		Name  string
		Email string
		Token string
	}{user.Id, user.Name, user.Email, token})
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Get user details using JWT
func GetUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	id := r.Header.Get("decoded")
	fmt.Print("id:", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.Models["user"].FindOne(r.Context(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonBytes, err := json.Marshal(struct {
		Id    primitive.ObjectID
		Name  string
		Email string
	}{user.Id, user.Name, user.Email})
	if err != nil {
		utils.Error(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
