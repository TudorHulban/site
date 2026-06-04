package domain

type AnalyticsRecord struct {
	// Timing (Stored as Unix timestamps)
	TimestampArrived  int64 `json:"timestamp_arrived"`  // Unix time (e.g., Milliseconds or Seconds)
	TimestampDeparted int64 `json:"timestamp_departed"` // Unix time (e.g., Milliseconds or Seconds)
	LatencyMicro      int64 `json:"latency_micro"`      // Latency in microseconds to avoid "0ms" results

	// Request Identity & Context
	ClientIP  string `json:"client_ip"` // Real IP (or masked IP for privacy)
	UserAgent string `json:"user_agent"`
	UserID    string `json:"user_id"` // Empty string if unauthenticated

	// Target & Routing
	Method       string `json:"method"`
	Path         string `json:"path"`
	MatchedRoute string `json:"matched_route"`
	RawQuery     string `json:"raw_query"`
	Referer      string `json:"referer"`

	// Outcome
	HTTPStatusCode int   `json:"http_status_code"` // Changed to match snake_case standard
	BytesWritten   int64 `json:"bytes_written"`
}
