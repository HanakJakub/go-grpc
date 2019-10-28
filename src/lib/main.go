package lib

import (
	"encoding/json"
)

// User struct to define user
type User struct {
	ID        uint64
	Purchases []Purchase
}

// Purchase struct to define purchase
type Purchase struct {
	Type   string
	Amount int64
}

// ShouldProcess will check if purchase should be processed or skipped
// and return bool
func (p Purchase) ShouldProcess() bool {
	if p.Amount < 150 || p.Amount > 10000 {
		return false
	}

	return true
}

// ToString will convert purchase to string
func (p Purchase) ToString() (string, error) {
	b, err := json.Marshal(p)

	return string(b), err
}
