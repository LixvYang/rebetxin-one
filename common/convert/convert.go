package convert

import "github.com/shopspring/decimal"

func String2Decimal(s string) decimal.Decimal {
	d, _ := decimal.NewFromString(s)
	return d
}
