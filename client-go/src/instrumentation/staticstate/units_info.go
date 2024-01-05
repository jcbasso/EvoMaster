package staticstate

// TODO: Review this struct

type UnitsInfo struct {
	UnitNamesSet   map[string]bool
	LinesCount     int64
	BranchCount    int64
	StatementCount int64
}

func NewUnitsInfo() *UnitsInfo {
	return &UnitsInfo{
		UnitNamesSet: map[string]bool{},
	}
}
