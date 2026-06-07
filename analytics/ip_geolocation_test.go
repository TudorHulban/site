package analytics

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestGetLocationByIP(t *testing.T) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	apiKey := ""
	targetIP := "82.76.117.202"

	location, errGeo := GetLocationByIP(
		&httpClient,
		targetIP,
		apiKey,
	)
	if errGeo != nil {
		fmt.Printf("Error: %v\n", errGeo)

		return
	}

	fmt.Print(*location)
}
