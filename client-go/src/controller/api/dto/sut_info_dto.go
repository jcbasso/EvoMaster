package dto

import (
	"github.com/jcbasso/EvoMaster/client-go/src/controller/api/dto/problem"
)

type SutInfoDto struct {
	// RestProblem If the SUT is a RESTful API, here there will be the info
	// on how to interact with it
	RestProblem problem.RestProblemDto `json:"restProblem"`

	// GraphQLProblem If the SUT is a GraphQL API, here there will be the info
	// on how to interact with it
	//GraphQLProblem problem.GraphQLProblemDto `json:"graphQLProblem,omitempty"`

	// IsSutRunning Whether the SUT is running or not
	IsSutRunning bool `json:"isSutRunning"`

	// DefaultOutputFormat When generating test cases for this SUT, specify the default
	// preferred output format (eg JUnit 4 in Java)
	DefaultOutputFormat string `json:"defaultOutputFormat"`

	// BaseUrlOfSUT The base URL of the running SUT (if any).
	// E.g., "http://localhost:8080"
	// It should only contain the protocol and the hostname/port
	BaseUrlOfSUT string `json:"baseUrlOfSUT,omitempty"`

	// InfoForAuthentication There is no way a testing system can guess passwords, even
	// if given full access to the database storing them (ie, reversing
	// hash values). As such, the SUT might need to provide a set of valid credentials
	InfoForAuthentication []AuthenticationDto `json:"infoForAuthentication"`

	// UnitsInfoDto Information about the "units" in the SUT.
	UnitsInfoDto UnitsInfoDto `json:"unitsInfoDto"`
}

// OutputFormat Note: this enum must be kept in sync with what declared in
// org.evomaster.core.output.OutputFormat
type OutputFormat int

const (
	JavaJunit5 OutputFormat = iota
	JavaJunit4
	KotlinJunit4
	KotlinJunit5
	JsJest
	CsharpXunit
	Go
)

func (e OutputFormat) String() string {
	switch e {
	case JavaJunit5:
		return "JAVA_JUNIT_5"
	case JavaJunit4:
		return "JAVA_JUNIT_4"
	case KotlinJunit4:
		return "KOTLIN_JUNIT_4"
	case KotlinJunit5:
		return "KOTLIN_JUNIT_5"
	case JsJest:
		return "JS_JEST"
	case CsharpXunit:
		return "CSHARP_XUNIT"
	case Go:
		return "GO"
	default:
		return ""
	}
}
