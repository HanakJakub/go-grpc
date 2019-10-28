package main

import (
	"go-grpc/src/lib"
	"testing"
)

func TestMedian(t *testing.T) {
	var res []lib.Purchase

	u := lib.User{
		ID: 5,
		Purchases: []lib.Purchase{
			lib.Purchase{Type: "test", Amount: 149},
			lib.Purchase{Type: "test", Amount: 150},
			lib.Purchase{Type: "test", Amount: 100},
			lib.Purchase{Type: "test", Amount: 10000},
			lib.Purchase{Type: "test", Amount: 10000},
			lib.Purchase{Type: "test", Amount: 10000},
			lib.Purchase{Type: "test", Amount: 10000},
		},
	}

	for _, p := range u.Purchases {
		if p.ShouldProcess() {
			res = append(res, p)
		}
	}

	if median(res) != float64(10000) {
		t.Error("Median should be 10000")
	}
}

func TestMean(t *testing.T) {
	var res []lib.Purchase

	u := lib.User{
		ID: 5,
		Purchases: []lib.Purchase{
			lib.Purchase{Type: "test", Amount: 149},
			lib.Purchase{Type: "test", Amount: 150},
			lib.Purchase{Type: "test", Amount: 100},
			lib.Purchase{Type: "test", Amount: 10000},
			lib.Purchase{Type: "test", Amount: 10000},
			lib.Purchase{Type: "test", Amount: 10000},
			lib.Purchase{Type: "test", Amount: 10000},
		},
	}

	for _, p := range u.Purchases {
		if p.ShouldProcess() {
			res = append(res, p)
		}
	}

	if mean(res) != float64(8030) {
		t.Error("Mean should be 8030")
	}
}
