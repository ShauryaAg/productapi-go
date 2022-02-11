package test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	testutils "github.com/ShauryaAg/ProductAPI/test/utils"
)

func TestReview(t *testing.T) {
	// Call the test method with the test cases
	// and the a custom request function
	test(t, "review_cases.json", reviewRequest)
}

// reviewRequest is a custom request function
// It sets the productId in the request context
func reviewRequest(method string, endpoint string, body *bytes.Buffer, headers interface{}) (*http.Request, error) {
	// get the default request
	req, err := defaultRequest(method, endpoint, body, headers)
	if err != nil {
		return nil, err
	}

	// add the productId to the request context
	if strings.Contains(endpoint, "{productId}") {
		req = testutils.SetUrlParamInContext(req, "productId", Id)
	}

	return req, nil
}
