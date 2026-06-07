package domain

import (
	"errors"
	"net/netip"
)

type Metric struct {
	RecordsPerPeriod uint32
	TopIPs           Meta[netip.Addr]
	TopBrowsers      Meta[Browser]
	TopASN           Meta[AsnEntity]
	TopCountries     Meta[string]
	TopCities        Meta[string]
}

type Day [24]Metric

type Registry struct {
	CurrentMonth [31]Day
	History      [7][31]Day
}

type DataCenter map[string]*Registry

type ParamsAddEvent struct {
	SiteKey string
	DayIdx  int // 1-31
	HourIdx int // 0-23
	IP      netip.Addr
	Browser Browser
	ASN     AsnEntity
	Country string
	City    string
}

func (dc DataCenter) AddEvent(params *ParamsAddEvent) error {
	// 1. Defensive Boundary Checks (Crucial for fixed array indices)
	if params.DayIdx < 1 || params.DayIdx > 31 {
		return errors.New("day index out of bounds (1-31)")
	}
	if params.HourIdx < 0 || params.HourIdx > 23 {
		return errors.New("hour index out of bounds (0-23)")
	}

	// 2. O(1) Fetch or Allocate the Registry for this site/datacenter
	reg, exists := dc[params.SiteKey]
	if !exists {
		reg = &Registry{}
		dc[params.SiteKey] = reg
	}

	// 3. Drill down directly to the exact memory address in the matrix
	// We do dayIdx-1 to map human days (1-31) to zero-based array indices (0-30)
	metricSlot := &reg.CurrentMonth[params.DayIdx-1][params.HourIdx]

	// 4. Update the flat metrics using pure pointer execution
	metricSlot.RecordsPerPeriod++

	metricSlot.TopIPs.Increment(params.IP)
	metricSlot.TopBrowsers.Increment(params.Browser)
	metricSlot.TopASN.Increment(params.ASN)

	if params.Country != "" {
		metricSlot.TopCountries.Increment(params.Country)
	}
	if params.City != "" {
		metricSlot.TopCities.Increment(params.City)
	}

	return nil
}
