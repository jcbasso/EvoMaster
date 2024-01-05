package dto

type ControllerInfoDto struct {
	// FullName The full qualifying name of the controller.
	// This will be needed when tests are generated, as those
	// will instantiate and start the controller directly
	FullName string `json:"fullName"`

	// IsInstrumentationOn Whether the system under test is running with instrumentation
	// to collect data about its execution
	IsInstrumentationOn bool `json:"isInstrumentationOn"`
}
