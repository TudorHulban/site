package analytics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ResponseGeoIP struct {
	Location struct {
		City        string `json:"city"`
		District    string `json:"district"`
		CountryCode string `json:"country_code3"`
		Postcode    string `json:"zipcode"`
		IsEU        bool   `json:"is_eu"`
	} `json:"location"`

	ASN struct {
		AsNumber     string `json:"as_number"`
		Organization string `json:"organization"`
		Country      string `json:"country"`
	} `json:"asn"`
}

func GetLocationByIP(client *http.Client, ipAddress, apiKey string) (*ResponseGeoIP, error) {
	// 1. Construct the URL safely
	baseURL, errParse := url.Parse("https://api.ipgeolocation.io/v3/ipgeo")
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
