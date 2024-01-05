package dto

// ExtraHeuristicsDto Represents possible extra heuristics related to the code
// execution and that do apply to all the reached testing targets.
// Example: rewarding SQL "select" operations that return non-empty sets
type ExtraHeuristicsDto struct {

	// List of extra heuristic values we want to optimize
	Heuristics []HeuristicEntryDto `json:"heuristics"`

	// TODO
	// databaseExecutionDto: ExecutionDto;
}
