package dto

// CookieLoginDto Information on how to do a login based on username/password, from which we then get a cookie back
type CookieLoginDto struct {
	// Username The id of the user
	Username string `json:"username"`

	// Password The password of the user. This must NOT be hashed.
	Password string `json:"password"`

	// UsernameField The name of the field in the body payload containing the username
	UsernameField string `json:"usernameField"`

	// PasswordField The name of the field in the body payload containing the password
	PasswordField string `json:"passwordField"`

	// LoginEndpointUrl The URL of the endpoint, e.g., "/login"
	LoginEndpointUrl string `json:"loginEndpointUrl"`

	// HttpVerb The HTTP verb used to send the data. Usually a "POST".
	HttpVerb string `json:"httpVerb"`

	// ContentType The encoding type used to specify how the data is sent
	ContentType string `json:"contentType"`
}

type ContentType int

const (
	ContentTypeJSON ContentType = iota
	ContentTypeXWWWFormUrlEncoded
)

func (e ContentType) String() string {
	switch e {
	case ContentTypeJSON:
		return "JSON"
	case ContentTypeXWWWFormUrlEncoded:
		return "X_WWW_FORM_URLENCODED"
	default:
		return ""
	}
}

type HttpVerb int

const (
	HttpVerbGet HttpVerb = iota
	HttpVerbPost
)

func (e HttpVerb) String() string {
	switch e {
	case HttpVerbGet:
		return "GET"
	case HttpVerbPost:
		return "POST"
	default:
		return ""
	}
}
