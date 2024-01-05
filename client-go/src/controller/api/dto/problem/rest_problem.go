package problem

type RestProblemDto struct {
	// OpenApiUrl The full URL of where the OpenAPI schema can be located
	OpenApiUrl string `json:"openApiUrl"`

	// EndpointsToSkip When testing a REST API, there might be some endpoints that are not so important to test.
	// For example, in Spring, health-check endpoints like "/heapdump" are not so interesting to test,
	// and they can be very expensive to run.
	EndpointsToSkip []string `json:"endpointsToSkip,omitempty"`
}
