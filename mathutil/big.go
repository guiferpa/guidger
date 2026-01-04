package mathutil

import (
	"math/big"
	"strings"
)

// InvalidAmountError represents an error when parsing an invalid amount
type InvalidAmountError struct {
	Amount string
}

func (e *InvalidAmountError) Error() string {
	return "invalid amount format: " + e.Amount
}

type Big struct {
	amount    *big.Rat
	precision int
}

// String returns the decimal representation with high precision,
// removing trailing zeros to reflect actual precision
func (b *Big) String() string {
	// Use high precision (100 decimals) to preserve all possible decimal places
	return b.amount.FloatString(b.precision)
}

// Add adds another Big value to this one and returns a new Big
// The precision is set to the maximum of the two precisions
func (b *Big) Add(other *Big) *Big {
	result := new(big.Rat).Add(b.amount, other.amount)
	maxPrecision := b.precision
	if other.precision > maxPrecision {
		maxPrecision = other.precision
	}
	return &Big{
		amount:    result,
		precision: maxPrecision,
	}
}

// GetPrecision returns the number of decimal places in the amount string
func GetPrecision(amount string) int {
	parts := strings.Split(amount, ".")
	if len(parts) == 1 {
		return 0
	}
	return len(parts[1])
}

// NewBig creates a new Big from a decimal string
// It preserves the precision (number of decimal places) from the input
func NewBig(amount string) (*Big, error) {
	amount = strings.TrimSpace(amount)
	if amount == "" {
		return nil, &InvalidAmountError{Amount: amount}
	}

	precision := GetPrecision(amount)
	rat := new(big.Rat)
	if _, ok := rat.SetString(amount); !ok {
		return nil, &InvalidAmountError{Amount: amount}
	}

	return &Big{
		amount:    rat,
		precision: precision,
	}, nil
}
