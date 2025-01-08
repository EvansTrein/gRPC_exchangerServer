package utils

import "fmt"

// data check
func ValidateCurrencyRequest(fromCurrency, toCurrency string) bool {
	if fromCurrency == "" || toCurrency == "" {
		return false
	}

	if fromCurrency == toCurrency {
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

// in the output we will get a map: the base currency and a slice of strings
// of all other currencies for which we need to make a rate
// example: key = USD value = [EUR, CNY, RUB] 
func GenerateCurrencyPairs(currencies []string) (map[string][]string, error) {
	pairs := make(map[string][]string)

	for _, baseCurrencieKey := range currencies {
		toCurrenciesValue := make([]string, 0, len(currencies)-1)

		for _, toCurrencie := range currencies {
			if baseCurrencieKey != toCurrencie {
				toCurrenciesValue = append(toCurrenciesValue, toCurrencie)
			}
		}
		// is a check that all currencies have been added as a pair to the base currency
		if len(toCurrenciesValue) != len(currencies)-1 {
			return nil, fmt.Errorf("slice of currency pairs of incorrect length")
		}

		pairs[baseCurrencieKey] = toCurrenciesValue
	}

	return pairs, nil
}
