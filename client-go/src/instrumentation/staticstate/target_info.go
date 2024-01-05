package staticstate

// TODO: Review this struct

type TargetInfo struct {
	MappedID      int64
	DescriptiveID string
	Value         float64
	ActionIndex   int64
}

func NewNotReachedTargetInfo(mappedID int64) TargetInfo {
	return TargetInfo{
		MappedID:    mappedID,
		Value:       0,  // Only will be 0 when it is not reached. The other ones will be at least 0.01
		ActionIndex: -1, // This means that it was not reached
	}
}
