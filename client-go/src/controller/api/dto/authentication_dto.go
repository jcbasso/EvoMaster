package dto

// AuthenticationDto To authenticate a user, would need specific settings,
// like specific values in the HTTP headers (eg, cookies)
type AuthenticationDto struct {
	// Name given name to this authentication info.
	// Just needed for display/debugging reasons
	Name string `json:"name"`

	// Headers The headers needed for authentication
	Headers []HeaderDto `json:"headers"`

	// CookieLogin If the login is based on cookies, need to provide info on how to get such a cookie
	CookieLogin CookieLoginDto `json:"cookieLogin"`

	JsonTokenPostLogin JsonTokenPostLoginDto `json:"jsonTokenPostLogin"`
}
