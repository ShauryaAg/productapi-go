package test

import (
	"testing"
)

func TestProduct(t *testing.T) {

	test(t, "product_cases.json", defaultRequest)
}
