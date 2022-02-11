package test

import (
	"testing"
)

func TestProduct(t *testing.T) {
	// Call the test method with the test cases
	// and the defaultRequest
	test(t, "product_cases.json", defaultRequest)
}
