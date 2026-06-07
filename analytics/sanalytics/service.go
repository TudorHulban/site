package sanalytics

import (
	"encoding/gob"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/tudorhulban/site/analytics/domain"
	"github.com/tudorhulban/site/analytics/shared"
)

type AnalyticsService struct {
	DC domain.DataCenter
}

func (s *AnalyticsService) HandleConnection(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)

	for {
		var req shared.FlatRequest
		// Decodes the incoming binary stream straight back into RAM
		if err := decoder.Decode(&req); err != nil {
			break // Connection closed or error
		}

		// Feed the flat request details directly into your matrix logic
		s.RecordEvent(&req)
	}
}

func (s *AnalyticsService) RecordEvent(ev *shared.FlatRequest) error {
	// 1. Instant IP Extraction from the string header
	host, _, errHost := net.SplitHostPort(ev.RemoteAddr)
	if errHost != nil {
		host = ev.RemoteAddr
	}

	ip, errParseIP := netip.ParseAddr(host)
	if errParseIP != nil {
		return errParseIP
	}

	// 2. High-speed User-Agent parsing for your top 7 browsers
	userAgent := ev.Header["User-Agent"]
	var uaString string
	if len(userAgent) > 0 {
		uaString = userAgent[0]
	}

	var browser domain.Browser

	switch {
	case strings.Contains(uaString, "Chrome"):
		browser = domain.Chrome

	case strings.Contains(uaString, "Safari") && !strings.Contains(uaString, "Chrome"):
		browser = domain.Safari

	case strings.Contains(uaString, "Edg"):
		browser = domain.Edge

	case strings.Contains(uaString, "Firefox"):
		browser = domain.Firefox

	case strings.Contains(uaString, "Brave"):
		browser = domain.Brave

	default:
		browser = 0 // Unknown
	}

	// 3. Extract Geo details (Passed downstream from your edge proxy headers)
	var country, city string

	if reqCountry := ev.Header["Cf-Ipcountry"]; len(reqCountry) > 0 {
		country = reqCountry[0]
	}
	if reqCity := ev.Header["X-Client-City"]; len(reqCity) > 0 {
		city = reqCity[0]
	}

	// 4. Resolve current time slots
	now := time.Now()
	dayIdx := now.Day()
	hourIdx := now.Hour()

	// 5. Route directly down to the raw bits in RAM
	// Using ev.Host as the DataCenter identifier
	return s.DC.AddEvent(
		&domain.ParamsAddEvent{
			DayIdx:  dayIdx,
			HourIdx: hourIdx,
			IP:      ip,
			Browser: browser,
			Country: country,
			City:    city,
		},
	)
}
