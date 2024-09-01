package metrics

const (
	MetricsSubsystem = "consulrpz"

	QueryStatusError       = "ERROR"
	QueryStatusDeny        = "DENY"
	QueryStatusFallthrough = "FALLTHROUGH"
	QueryStatusSuccess     = "SUCCESS"
	QueryStatusNoMatch     = "NOMATCH"
)
