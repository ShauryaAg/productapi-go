package test

import (
	"bytes"
	"net/http"
	"strings"
	"testing"

	testutils "github.com/ShauryaAg/ProductAPI/test/utils"
)

func TestReview(t *testing.T) {
	test(t, "review_cases.json", reviewRequest)
}

func reviewRequest(method string, endpoint string, body *bytes.Buffer, headers interface{}) (*http.Request, error) {
	req, err := defaultRequest(method, endpoint, body, headers)
	if err != nil {
		return nil, err
	}

	if strings.Contains(endpoint, "{productId}") {
		req = testutils.SetUrlParamInContext(req, "productId", Id)
	}

	return req, nil
}
