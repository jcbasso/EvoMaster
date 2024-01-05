package controller

type SutHandler interface {
	// StartSut Start a new instance of the SUT.
	// This method must be blocking until the SUT is initialized.
	// returns the base URL of the running SUT. (for example: "http://localhost:8080")
	StartSut() string

	// StopSut Stop the SUT
	StopSut()

	// ResetStateOfSUT Make sure the SUT is in a clean state (eg, reset data in database).
	// A possible (likely very inefficient) way to implement this would be to
	// call #StopSUT followed by #StartSUT
	ResetStateOfSUT()

	// SetHost Sets the Host to be used by the SUT. (for example: "localhost")
	SetHost(string)

	// SetPort Sets the Port to be used by the SUT. (for example: 8080)
	SetPort(int)
}
