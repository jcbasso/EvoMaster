package dto

// WrappedResponseDto In REST, when we have an error, at most we would see an HTTP status code.
// But it can be very useful to get an actual description of the error.
// So, it is a common practice to have "Wrapped Responses", which can contain the error message (if any)
type WrappedResponseDto[T any] struct {
	// Data The actual payload we are sending and are "wrapping" here
	Data T `json:"data,omitempty"`
	// Error If this is not null, then "data" must be null.
	Error string `json:"error,omitempty"`
}

func NewWrappedResponseDtoWithData[T any](data T) WrappedResponseDto[T] {
	return WrappedResponseDto[T]{
		Data: data,
	}
}

func NewWrappedResponseDtoWithNoData() WrappedResponseDto[interface{}] {
	return WrappedResponseDto[interface{}]{}
}

func NewWrappedResponseDtoWithError(msg string) WrappedResponseDto[interface{}] {
	return WrappedResponseDto[interface{}]{
		Error: msg,
	}
}
