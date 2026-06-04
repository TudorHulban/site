package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ResponseGeoIP struct {
	City struct {
		Name string `json:"name"`
	} `json:"city"`

	Country struct {
		Name string `json:"name"`
		Code string `json:"iso_code"`
	} `json:"country"`

	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`

	Postcode string `json:"postcode"`
}

func GetLocationByIP(client *http.Client, ipAddress, apiKey string) (*ResponseGeoIP, error) {
	// 1. Construct the URL safely
	baseURL, errParse := url.Parse("https://api.geoapify.com/v1/ipinfo")
	if errParse != nil {
		return nil,
			fmt.Errorf("failed parsing base URL: %w", errParse)
	}

	params := url.Values{}
	params.Add("ip", ipAddress)
	params.Add("apiKey", apiKey)
	baseURL.RawQuery = params.Encode()

	// 2. Execute the HTTP GET request
	resp, errGet := client.Get(baseURL.String())
	if errGet != nil {
		return nil,
			fmt.Errorf("http request failed: %w", errGet)
	}
	defer resp.Body.Close()

	// 3. Handle non-200 responses safely
	if resp.StatusCode != http.StatusOK {
		return nil,
			fmt.Errorf("geoapify API returned status: %d", resp.StatusCode)
	}

	// 4. Decode the JSON stream directly into the struct (more efficient than io.ReadAll)
	var geoData ResponseGeoIP

	if errDecoder := json.NewDecoder(resp.Body).Decode(&geoData); errDecoder != nil {
		return nil,
			fmt.Errorf("failed decoding json response: %w", errDecoder)
	}

	return &geoData, nil
}
