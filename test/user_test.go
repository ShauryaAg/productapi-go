package test

import "testing"

func TestUser(t *testing.T) {
	// Call the test method with the test cases
	// and the default request method
	test(t, "user_cases.json", defaultRequest)
}
