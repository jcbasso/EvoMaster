package dto

type TestResultsDto struct {
	Targets            []TargetInfoDto      `json:"targets"`
	AdditionalInfoList []AdditionalInfoDto  `json:"additionalInfoList"`
	ExtraHeuristics    []ExtraHeuristicsDto `json:"extraHeuristics"`
}
