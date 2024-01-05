package dto

// Type The type of extra heuristic.
// Note: for the moment, we only have heuristics on SQL commands
type Type int

const (
	Sql Type = iota
)

func (t Type) String() string {
	switch t {
	case Sql:
		return "SQL"
	default:
		return ""
	}
}

// Objective Should we try to minimize or maximize the heuristic?
type Objective int

const (
	// MINIMIZE_TO_ZERO The lower, the better.
	// Minimum is 0. It can be considered as a "distance" to minimize.
	MINIMIZE_TO_ZERO Objective = iota

	// MAXIMIZE The higher, the better.
	// Note: given x, we could rather consider the value 1/x to minimize. But that wouldn't work for negative x,
	// and also would make debugging more difficult (ie better to look at the raw, non-transformed values).
	MAXIMIZE
)

func (o Objective) String() string {
	switch o {
	case MINIMIZE_TO_ZERO:
		return "MINIMIZE_TO_ZERO"
	case MAXIMIZE:
		return "MAXIMIZE"
	default:
		return ""
	}
}

type HeuristicEntryDto struct {
	Type Type `json:"type"`

	Objective Objective `json:"objective"`

	// Id An id representing these heuristics.
	// For example, for SQL, it could be a SQL command
	Id string `json:"id"`

	// Value The actual value of the heuristic
	Value int `json:"value"`
}
