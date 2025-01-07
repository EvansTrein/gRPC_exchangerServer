package utils

func ValidateCurrencyRequest(fromCurrency, toCurrency string) bool {
	if fromCurrency == "" || toCurrency == "" {
		return false
	}

	if len(fromCurrency) > 5 || len(fromCurrency) < 3 {
		return false
	}

	if len(toCurrency) > 5 || len(toCurrency) < 3 {
		return false
	}

	return true
}
