package test

import "testing"

func TestUser(t *testing.T) {
	test(t, "user_cases.json", defaultRequest)
}
