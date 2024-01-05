package dto

type SutRunDto struct {
	// Run Whether the SUT should be running
	Run bool `json:"run"`

	// ResetState Whether the internal state of the SUT should be reset
	ResetState bool `json:"resetState"`

	// CalculateSqlHeuristics Whether SQL heuristics should be computed.
	// Note: those can be very expensive
	CalculateSqlHeuristics bool `json:"calculateSqlHeuristics"`

	// ExtractSqlExecutionInfo Whether SQL execution info should be saved.
	ExtractSqlExecutionInfo bool `json:"extractSqlExecutionInfo"`
}
