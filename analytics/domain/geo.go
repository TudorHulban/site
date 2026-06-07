package domain

type Country struct {
	ID   string `db:"id"`
	Name string `db:"name"`
	IsEU bool   `db:"is_eu"`
}

type City struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	District  string `db:"district"`
	CountryID string `db:"country_id"`
}

// AsnEntity represents the parent organization/ISP.
type AsnEntity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// AsnNumber links a specific ASN string to a parent entity and a country.
type AsnNumber struct {
	ID        int    `db:"id"`
	ASN       string `db:"asn"`
	EntityID  int    `db:"entity_id"`
	CountryID string `db:"country_id"`
}

// The core log record that ties it all together via IDs.
type IPLog struct {
	ID        int64  `db:"id"`
	IPAddress string `db:"ip_address"`
	Postcode  string `db:"postcode"`
	CityID    int    `db:"city_id"`
	ASNID     int    `db:"asn_id"`
}
