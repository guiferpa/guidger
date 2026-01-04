package mathutil

import (
	"math/big"
	"testing"
)

func TestNewBig(t *testing.T) {
	tests := []struct {
		amount string
		want   *Big
	}{
		{amount: "1.0", want: &Big{amount: big.NewRat(1, 1), precision: 1}},
		{amount: "1.00", want: &Big{amount: big.NewRat(1, 1), precision: 2}},
		{amount: "1.000", want: &Big{amount: big.NewRat(1, 1), precision: 3}},
		{amount: "1.0000", want: &Big{amount: big.NewRat(1, 1), precision: 4}},
		{amount: "1.00000", want: &Big{amount: big.NewRat(1, 1), precision: 5}},
		{amount: "1.000000", want: &Big{amount: big.NewRat(1, 1), precision: 6}},
		{amount: "1.0000000", want: &Big{amount: big.NewRat(1, 1), precision: 7}},
		{amount: "1.00000000", want: &Big{amount: big.NewRat(1, 1), precision: 8}},
		{amount: "1.000000000", want: &Big{amount: big.NewRat(1, 1), precision: 9}},
		{amount: "1.0000000000", want: &Big{amount: big.NewRat(1, 1), precision: 10}},
		{amount: "1.00000000000", want: &Big{amount: big.NewRat(1, 1), precision: 11}},
		{amount: "1.000000000000", want: &Big{amount: big.NewRat(1, 1), precision: 12}},
		{amount: "1.0000000000000", want: &Big{amount: big.NewRat(1, 1), precision: 13}},
		{amount: "1.00000000000000", want: &Big{amount: big.NewRat(1, 1), precision: 14}},
		{amount: "1.000000000000000", want: &Big{amount: big.NewRat(1, 1), precision: 15}},
	}
	for _, test := range tests {
		got, err := NewBig(test.amount)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.String() != test.want.String() {
			t.Fatalf("expected %s, got %s", test.want.String(), got.String())
		}
	}
}

func TestBigAdd(t *testing.T) {
	tests := []struct {
		a    string
		b    string
		want string
	}{
		{a: "1.0", b: "1.0", want: "2.0"},
		{a: "1.00", b: "1.00", want: "2.00"},
		{a: "1.000", b: "1.000", want: "2.000"},
		{a: "1.0000", b: "1.0000", want: "2.0000"},
		{a: "1.00000", b: "1.00000", want: "2.00000"},
		{a: "1.000000", b: "1.000000", want: "2.000000"},
		{a: "1.0000000", b: "1.0000000", want: "2.0000000"},
		{a: "1.00000000", b: "1.00000000", want: "2.00000000"},
		{a: "1.000000000", b: "1.000000000", want: "2.000000000"},

		// Tests for different precisions
		{a: "1.0", b: "1.00", want: "2.00"},
		{a: "1.00", b: "1.0", want: "2.00"},
		{a: "1.00", b: "1.000", want: "2.000"},
		{a: "1.000", b: "1.00", want: "2.000"},
		{a: "1.000", b: "1.0000", want: "2.0000"},
		{a: "1.0000", b: "1.000", want: "2.0000"},
		{a: "1.0000", b: "1.00000", want: "2.00000"},
		{a: "1.00000", b: "1.0000", want: "2.00000"},

		// Tests with subtraction and different precisions
		{a: "1.0", b: "-0.5", want: "0.5"},
		{a: "1.00", b: "-0.50", want: "0.50"},
		{a: "1.000", b: "-0.500", want: "0.500"},
		{a: "1.0000", b: "-0.5000", want: "0.5000"},
	}
	for _, test := range tests {
		a, err := NewBig(test.a)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		b, err := NewBig(test.b)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		got := a.Add(b)
		if got.String() != test.want {
			t.Fatalf("expected %s, got %s", test.want, got.String())
		}
	}
}
