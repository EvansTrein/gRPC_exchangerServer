package currencyrateapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/EvansTrein/gRPC_exchangerServer/internal/storages"
)

// example url API https://api.fxratesapi.com/latest?api_key=YOUR_ACCESS_TOKEN&currencies=EUR,RUB,CNY&base=USD

const urlAPI = "https://api.fxratesapi.com/latest"
const myApiKey = "fxr_live_80e48706cea1f3527bc8bbb9b1e43137023a"

type exchangeRatesResponse struct {
	Success     bool               `json:"success"`
	Error       string             `json:"error"`
	Description string             `json:"description"`
	Base        string             `json:"base"`
	Rates       map[string]float32 `json:"rates"`
}

func DownloadExchangeRateData(baseCurrencyCode string, toCurrencysCodes []string) ([]storages.Rate, error) {
	var builder strings.Builder
	var url string

	builder.WriteString(urlAPI)

	builder.WriteString("?api_key=")
	builder.WriteString(myApiKey)

	builder.WriteString("&currencies=")
	builder.WriteString(strings.Join(toCurrencysCodes, ","))

	builder.WriteString("&base=")
	builder.WriteString(baseCurrencyCode)

	url = builder.String()

	// you can send it that way, too, without the key
	// url = urlAPI + "?" + strings.Join(toCurrencysCodes, ",") + "&" + baseCurrencyCode

	log.Printf("sending a request for %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute http request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response body: %v", err)
	}

	var response exchangeRatesResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse the response into JSON: %v", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("failed to get data from API\nerror: %s\ndescription: %s", response.Error, response.Description)
	}

	var rates []storages.Rate
	for toCurrency, rate := range response.Rates {
		rates = append(rates, storages.Rate{
			BaseCurrency: response.Base,
			ToCurrency:   toCurrency,
			Rate:         rate,
		})
	}

	return rates, nil
}
