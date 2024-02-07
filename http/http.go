// package http implements various http related helper types and functions.
package http

import "net/http"

// RequestUnmarshaler is any type capable to unmarshal data from HTTP request to itself.
type RequestUnmarshaler interface {
	UnmarshalHTTPRequest(r *http.Request) error
}

// UnmarshalRequest unmarshals HTTP request into m.
func UnmarshalRequest(r *http.Request, m RequestUnmarshaler) error {
	return m.UnmarshalHTTPRequest(r)
}

// ResponseUnmarshaler is any type capable to unmarshal data from HTTP request response to itself.
type ResponseUnmarshaler interface {
	UnmarshalHTTPResponse(r *http.Response) error
}

// UnmarshalResponse unmarshals HTTP response into m.
func UnmarshalResponse(r *http.Response, m ResponseUnmarshaler) error {
	return m.UnmarshalHTTPResponse(r)
}

// Marshaler is any type capable to marshal data to suitable HTTP request payload.
type Marshaler interface {
	MarshalHTTP(w http.ResponseWriter) error
}

// Marshal marshals v into w.
func Marshal(w http.ResponseWriter, m Marshaler) error {
	return m.MarshalHTTP(w)
}
