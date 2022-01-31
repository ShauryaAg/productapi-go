package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShauryaAg/ProductAPI/handlers"
	"github.com/ShauryaAg/ProductAPI/middlewares"
	"github.com/ShauryaAg/ProductAPI/models"
	"github.com/ShauryaAg/ProductAPI/models/db"
	testutils "github.com/ShauryaAg/ProductAPI/test/utils"
	utils "github.com/ShauryaAg/ProductAPI/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestCase struct {
	Name         string      `json:"name"`
	Endpoint     string      `json:"endpoint"`
	Method       string      `json:"method"`
	Handler      string      `json:"handler"`
	InputBody    interface{} `json:"inputBody"`
	InputHeaders interface{} `json:"inputHeaders"`
	Expected     interface{} `json:"expected"`
}

const (
	baseUrl = "http://localhost:8080"
)

var handlersMap map[string]http.HandlerFunc = map[string]http.HandlerFunc{
	"LOGIN":          handlers.Login,
	"REGISTER":       handlers.Register,
	"GET_USER":       (middlewares.AuthMiddleware(http.HandlerFunc(handlers.GetUser))).ServeHTTP,
	"CREATE_PRODUCT": (middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreateProduct))).ServeHTTP,
	"SEARCH_PRODUCT": handlers.SearchProducts,
	"CREATE_REVIEW":  (middlewares.AuthMiddleware(http.HandlerFunc(handlers.CreateReview))).ServeHTTP,
}

var (
	JwtToken string
	Id       string
)

func setUp(ctx context.Context, db *mongo.Database) (string, string) {
	db.Collection("user").DeleteMany(ctx, bson.M{})
	db.Collection("product").DeleteMany(ctx, bson.M{})
	db.Collection("review").DeleteMany(ctx, bson.M{})

	id := primitive.NewObjectID()

	// Create User
	user, _ := models.NewUser("Shaurya", "abc@xyz.com", "password")
	user.Id = id
	_, _ = db.Collection("user").InsertOne(ctx, user)
	jwtToken, _ := utils.CreateToken(*user)

	// Create Product
	product, _ := models.NewProduct("Product1", "Description1", "http://image1.com")
	product.Id = id
	_, _ = db.Collection("product").InsertOne(ctx, product)

	// Create Review
	review, _ := models.NewReview("Review1", 3, *user)
	review.Id = id
	_, _ = db.Collection("review").InsertOne(ctx, review)

	return jwtToken, id.Hex()
}

func test(
	t *testing.T,
	casesFile string,
	newRequest func(method string, endpoint string, body *bytes.Buffer, headers interface{}) (*http.Request, error),
) {
	dbName := "test"
	ctx := context.TODO()
	client, err := db.InitDatabase(dbName, ctx)
	if err != nil {
		t.Error(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	JwtToken, Id = setUp(ctx, db)

	jsonFile, err := ioutil.ReadFile(casesFile)
	if err != nil {
		t.Error(err)
	}

	var testCases []TestCase
	err = json.Unmarshal(jsonFile, &testCases)
	if err != nil {
		t.Error(err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log(testCase.Name)
			jsonBody, err := json.Marshal(testCase.InputBody)
			if err != nil {
				t.Error(err)
			}

			req, err := newRequest(testCase.Method, baseUrl+testCase.Endpoint, bytes.NewBuffer(jsonBody), testCase.InputHeaders)
			if err != nil {
				t.Error(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlersMap[testCase.Handler])
			handler.ServeHTTP(rr, req)

			expected := testCase.Expected.(map[string]interface{})
			if status := rr.Code; status != int(expected["status"].(float64)) {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, expected["status"])
			}

			body := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &body)
			if err != nil {
				t.Error(err)
			}

			expectedResponse := expected["response"].(map[string]interface{})
			if !testutils.MatchMaps(expectedResponse, body) {
				t.Errorf("handler returned wrong body: got %v want %v",
					body, expected["response"])
			}
		})
	}

	t.Cleanup(func() {
		db.Collection("users").DeleteMany(ctx, bson.M{})
		db.Collection("product").DeleteMany(ctx, bson.M{})
		db.Collection("review").DeleteMany(ctx, bson.M{})
		db.Drop(ctx)
		client.Disconnect(ctx)
	})
}

func defaultRequest(method string, endpoint string, body *bytes.Buffer, headers interface{}) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if headers != nil {
		for k, v := range headers.(map[string]interface{}) {
			if k == "Authorization" {
				req.Header.Set("Authorization", "Bearer "+JwtToken)
			} else {
				req.Header.Set(k, v.(string))
			}
		}
	}

	return req, nil
}
