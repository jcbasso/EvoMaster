package dto

// UnitsInfoDto Information about the "units" in the SUT.
// In case of OO languages like Java and Kotlin, those will be "classes"
type UnitsInfoDto struct {

	// UnitNames Then name of all the units (eg classes) in the SUT
	UnitNames []string `json:"unitNames"`

	// NumberOfLines The total number of lines/statements/instructions in all units of the whole SUT
	NumberOfLines int `json:"numberOfLines"`

	// NumberOfBranches The total number of branches in all units of the whole SUT
	NumberOfBranches int `json:"numberOfBranches"`

	// Number of replaced method testability transformations. But only for SUT units.
	NumberOfReplacedMethodsInSut int `json:"numberOfReplacedMethodsInSut,omitempty"`

	// Number of replaced method testability transformations.
	// But only for third-party library units (ie all units not in the SUT).
	NumberOfReplacedMethodsInThirdParty int `json:"numberOfReplacedMethodsInThirdParty,omitempty"`

	// Number of tracked methods. Those are special methods for which
	// we explicitly keep track of how they are called (eg their inputs).
	NumberOfTrackedMethods int `json:"numberOfTrackedMethods,omitempty"`
}
