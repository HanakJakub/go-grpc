package lib

import "testing"

func TestMain(t *testing.T) {
	var res []Purchase

	u := User{
		ID: 5,
		Purchases: []Purchase{
			Purchase{Type: "test", Amount: 149},
			Purchase{Type: "test", Amount: 150},
			Purchase{Type: "test", Amount: 100},
			Purchase{Type: "test", Amount: 10000},
			Purchase{Type: "test", Amount: 10001},
		},
	}

	for _, p := range u.Purchases {
		if p.ShouldProcess() {
			res = append(res, p)
		}
	}

	if len(res) != 2 {
		t.Error("ShouldProccess should allow only 2 purchases to pass")
	}
}
