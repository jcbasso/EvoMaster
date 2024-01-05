package problem

type GraphQLProblemDto struct {
	// Endpoint The endpoint in the SUT that expect incoming GraphQL queries
	Endpoint string `json:"endpoint"`
}
