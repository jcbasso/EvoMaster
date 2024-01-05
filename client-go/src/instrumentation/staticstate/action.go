package staticstate

// TODO: Review this struct

type Action struct {
	Index          int64
	InputVariables []string
}

func NewAction(index int64, inputVariables []string) *Action {
	return &Action{
		Index:          index,
		InputVariables: inputVariables,
	}
}
