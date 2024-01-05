package dto

// AllTargetsDto returns a list of all the targets
type AllTargetsDto struct {
	Files      []string `json:"files"`
	Lines      []string `json:"lines"`
	Branches   []string `json:"branches"`
	Statements []string `json:"statements"`
}

type CoverageDto struct {
	DescriptiveID string  `json:"descriptive_id"`
	MappedID      int64   `json:"mapped_id"`
	Value         float64 `json:"value"`
	ActionIndex   int64   `json:"action_index"`
}

type FullObjectiveCoverageDto struct {
	Files      []CoverageDto `json:"files"`
	Lines      []CoverageDto `json:"lines"`
	Branches   []CoverageDto `json:"branches"`
	Statements []CoverageDto `json:"statements"`
}
