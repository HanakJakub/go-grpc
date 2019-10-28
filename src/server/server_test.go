package main

import (
	"testing"
)

func TestGivenPathShouldCreateUserStruct(t *testing.T) {

	f := loadUserFile("../../task/data/1.json")
	user := parseUser("1.json", f)

	if user.ID != 1 || len(user.Purchases) != 5 {
		t.Error("User is not parsed correctly")
	}
}
