package dto

type AdditionalInfoDto struct {

	// QueryParameters In REST APIs, it can happen that some query parameters do not appear in the schema if they are
	// indirectly accessed via objects like WebRequest.
	// But we can track at runtime when such kind of objects are used to access the query parameters
	QueryParameters []string `json:"queryParameters"`

	// Headers In REST APIs, it can happen that some HTTP headers do not appear in the schema if they are indirectly
	// accessed via objects like WebRequest.
	// But we can track at runtime when such kind of objects are used to access the query parameters
	Headers []string `json:"headers"`

	// StringSpecializations Information for taint analysis.
	// When some string inputs are recognized of a specific type (eg, they are used as integers or dates),
	// we keep track of it.
	// The key in this map is the value of the tainted input.
	// The associated list is its possible specializations (which usually will be at most 1).
	// TODO: Implement string specialization
	//StringSpecializations map[string]StringSpecializationInfoDto

	// LastExecutedStatement Keep track of the last executed statement done in the SUT.
	// But not in the third-party libraries, just the business logic of the SUT.
	// The statement is represented with a descriptive unique id, like the class name and line number.
	LastExecutedStatement string `json:"lastExecutedStatement"`
}
