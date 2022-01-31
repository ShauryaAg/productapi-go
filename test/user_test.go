package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ShauryaAg/ProductAPI/handlers"
	"github.com/ShauryaAg/ProductAPI/middlewares"
	"github.com/ShauryaAg/ProductAPI/models/db"
	"github.com/ShauryaAg/ProductAPI/test/utils"
	"go.mongodb.org/mongo-driver/bson"
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

func TestMain(t *testing.T) {
	client, err := db.InitDatabase("test", context.TODO())
	if err != nil {
		t.Error(err)
	}

	jsonFile, err := ioutil.ReadFile("cases.json")
	if err != nil {
		t.Error(err)
	}

	var testCases []TestCase
	err = json.Unmarshal(jsonFile, &testCases)
	if err != nil {
		t.Error(err)
	}

	var jwtToken string
	var productId string
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Log(testCase.Name)
			jsonBody, err := json.Marshal(testCase.InputBody)
			if err != nil {
				t.Error(err)
			}

			req, err := http.NewRequest(testCase.Method, baseUrl+testCase.Endpoint, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Error(err)
			}

			if strings.Contains(testCase.Endpoint, "{productId}") {
				req = utils.SetUrlParamInContext(req, "productId", productId)
			}

			req.Header.Set("Content-Type", "application/json")
			if testCase.InputHeaders != nil {
				for k, v := range testCase.InputHeaders.(map[string]interface{}) {
					if k == "Authorization" {
						req.Header.Set("Authorization", "Bearer "+jwtToken)
					} else {
						req.Header.Set(k, v.(string))
					}
				}
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
			if !utils.MatchMaps(expectedResponse, body) {
				t.Errorf("handler returned wrong body: got %v want %v",
					body, expected["response"])
			} else {
				if testCase.Name == "CreateUserSuccess" {
					jwtToken = body["Token"].(string)
				}
				if testCase.Name == "CreateProductSuccess" {
					productId = body["Id"].(string)
				}
			}
		})
	}

	t.Cleanup(func() {
		client.Database("test").Collection("users").DeleteMany(context.TODO(), bson.M{})
		client.Database("test").Drop(context.TODO())
		client.Disconnect(context.TODO())
	})
}
