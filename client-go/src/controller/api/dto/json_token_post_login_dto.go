package dto

type JsonTokenPostLoginDto struct {
	// UserId The id representing this user that is going to logi
	UserId string `json:"userId"`

	// Endpoint The endpoint where to execute the login
	Endpoint string `json:"endpoint"`

	// JsonPayload The payload to send, as stringified JSON object
	JsonPayload string `json:"jsonPayload"`

	// ExtractTokenField How to extract the token from a JSON response, as such JSON could have few fields, possibly nested.
	// It is expressed as a JSON Pointe
	ExtractTokenField string `json:"extractTokenField"`

	// HeaderPrefix When sending out the obtained token in the Authorization header,
	// specify if there should be any prefix (e.g., "Bearer " or "JWT "
	HeaderPrefix string `json:"headerPrefix"`
}
