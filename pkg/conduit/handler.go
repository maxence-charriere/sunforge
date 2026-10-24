package conduit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/maxence-charriere/go-app/v11/pkg/logs"
)

// Mux is an HTTP request multiplexer. It matches the URL of each incoming
// request against request methods and a list of registered patterns and calls
// the handler for the pattern that most closely matches the URL.
type Mux struct {
	muxes map[string]*http.ServeMux
}

// NewMux creates a Mux.
func NewMux() Mux {
	return Mux{
		muxes: make(map[string]*http.ServeMux),
	}
}

// HandleFunc registers the handler function for the given method and pattern.
// The documentation for ServeMux explains how patterns are matched.
func (m Mux) HandleFunc(method string, pattern string, h http.HandlerFunc) {
	m.Handle(method, pattern, h)
}

// Handle registers the handler for the given method and pattern. The
// documentation for http.ServeMux explains how patterns are matched.
func (m Mux) Handle(method string, pattern string, h http.Handler) {
	mux, ok := m.muxes[method]
	if !ok {
		mux = http.NewServeMux()
		m.muxes[method] = mux
	}
	mux.Handle(pattern, h)
}

// ServeHTTP dispatches the request to the handler whose method and pattern most
// closely match the request URL.
func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux, ok := m.muxes[r.Method]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	mux.ServeHTTP(w, r)
}

// BaseHandler is a struct to embed that provides helper methods to handle a
// RESTful API.
type BaseHandler struct {
	// The function used to encode JSON. Default is json.Marshal.
	Encoder func(any) ([]byte, error)

	// The function to decode JSON. Default is json.Unmarshal.
	Decoder func([]byte, any) error
}

// Error writes an HTTP error.
func (h BaseHandler) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if code >= 500 {
		app.Log(errors.New("http request failed").
			WithTag("method", r.Method).
			WithTag("path", r.URL.Path).
			WithTag("code", code).
			Wrap(err))
	} else {
		app.Log(logs.New("http request failed").
			WithTag("method", r.Method).
			WithTag("path", r.URL.Path).
			WithTag("code", code).
			WithTag("error", err))
	}
	http.Error(w, "", code)
}

// Error writes a bad request HTTP error.
func (h BaseHandler) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, http.StatusBadRequest, err)
}

// Error writes an unauthorized HTTP error.
func (h BaseHandler) Unauthorized(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, http.StatusUnauthorized, err)
}

// Error writes a too many requests HTTP error.
func (h BaseHandler) TooManyRequests(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, http.StatusTooManyRequests, err)
}

// Error writes an internal server error HTTP error.
func (h BaseHandler) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, http.StatusInternalServerError, err)
}

// Error writes an not found error HTTP error.
func (h BaseHandler) NotFound(w http.ResponseWriter, r *http.Request, err error) {
	h.Error(w, r, http.StatusNotFound, err)
}

// Ok writes an ok response and encode the given values. Values are encoded as
// an array when there is more than 1.
func (h BaseHandler) Ok(w http.ResponseWriter, r *http.Request, v ...any) {
	var body []byte
	var err error

	if len(v) == 1 {
		body, err = h.Encode(v[0])
	} else if len(v) != 0 {
		body, err = h.Encode(v)
	}
	if err != nil {
		h.InternalServerError(w, r, errors.New("encoding response failed").Wrap(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	if len(body) == 0 {
		return
	}
	w.Write(body)
}

// PreventCaching sets response headers that disable caching in browsers and
// shared intermediaries such as CDNs.
func (h BaseHandler) PreventCaching(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, max-age=0, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Surrogate-Control", "no-store")
}

// NotModified writes a not modified response.
func (h BaseHandler) NotModified(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotModified)
}

// Encodes the given value to JSON.
func (h BaseHandler) Encode(v any) ([]byte, error) {
	if h.Encoder == nil {
		return json.Marshal(v)
	}
	return h.Encoder(v)
}

// Decodes the byte array and stores the result into the given value.
func (h BaseHandler) Decode(b []byte, v any) error {
	if h.Decoder == nil {
		return json.Unmarshal(b, v)
	}
	return h.Decoder(b, v)
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *responseWriter) Write(v []byte) (int, error) {
	w.body.Write(v)
	return w.ResponseWriter.Write(v)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
