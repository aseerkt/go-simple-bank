package utils

import "github.com/aseerkt/go-simple-bank/pkg/constants"

func IsSupportedCurrency(currency string) bool {
	_, ok := constants.Currency[currency]
	return ok
}
