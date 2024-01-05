package dto

// HeaderDto A HTTP header, which is a key=value pair
type HeaderDto struct {
	// Name The header name
	Name string `json:"name"`

	// Value The value of the header
	Value string `json:"value"`
}
