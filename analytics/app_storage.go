package analytics

import (
	"database/sql"
	"fmt"
	"time"
)

// AggregateFunc handles the heavy lifting for a single metric domain
type AggregateFunc func(tx *sql.Tx, activeDB *sql.DB, start, end time.Time) error

// AppStorage manages database pool dependencies
type AppStorage struct {
	ActiveDB    *sql.DB
	AnalyticsDB *sql.DB
	Registry    map[string]AggregateFunc
}

// NewAppStorage initializes the storage layer and hooks up the registry
func NewAppStorage(activeDB, analyticsDB *sql.DB) *AppStorage {
	return &AppStorage{
		ActiveDB:    activeDB,
		AnalyticsDB: analyticsDB,

		Registry: map[string]AggregateFunc{
			"City Aggregates": AggregateCities,
			"ASN Aggregates":  AggregateASNs,
		},
	}
}

func (s *AppStorage) ProcessIntervalMetrics(start time.Time, end time.Time) error {
	tx, err := s.AnalyticsDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// The map key naturally replaces the .Name() method or struct property
	for name, aggregate := range s.Registry {
		if errAggregate := aggregate(tx, s.ActiveDB, start, end); errAggregate != nil {
			// Transparent, out-of-the-box contextual error logging
			// e.g., logger.Errorf("processor [%s] failed: %v", name, err)
			return fmt.Errorf(
				"aggregator %s: %w",
				name,
				errAggregate,
			)
		}
	}

	return tx.Commit()
}

func AggregateCities(tx *sql.Tx, activeDB *sql.DB, start, end time.Time) error {
	rows, err := activeDB.Query(`
		SELECT city_id, COUNT(id) FROM raw_ip_hits 
		WHERE created_at >= ? AND created_at < ? AND city_id IS NOT NULL
		GROUP BY city_id`, start, end)
	if err != nil {
		return err
	}
	defer rows.Close()

	// ... rest of SQL insertion logic ...
	return nil
}

func AggregateASNs(tx *sql.Tx, activeDB *sql.DB, start, end time.Time) error {
	rows, err := activeDB.Query(`
		SELECT asn_id, COUNT(id) FROM raw_ip_hits 
		WHERE created_at >= ? AND created_at < ? AND asn_id IS NOT NULL
		GROUP BY asn_id`, start, end)
	if err != nil {
		return err
	}
	defer rows.Close()

	// ... rest of SQL insertion logic ...
	return nil
}
