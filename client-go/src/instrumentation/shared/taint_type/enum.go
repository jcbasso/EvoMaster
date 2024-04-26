package taint_type

type TaintType int

const (
	NONE TaintType = iota
	FULL_MATCH
	PARTIAL_MATCH
	End
)

func (s TaintType) String() string {
	switch s {
	case NONE:
		return "NONE"
	case FULL_MATCH:
		return "FULL_MATCH"
	case PARTIAL_MATCH:
		return "PARTIAL_MATCH"
	default:
		return ""
	}
}
