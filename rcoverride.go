package traefik_plugin_response_code_override

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
)

// Config the plugin configuration.
type Config struct {
	Overrides       map[int]int `json:"overrides,omitempty"`
	RemoveBody      bool        `json:"remove_body,omitempty"`
	HeadersToRemove []string    `json:"headers_to_remove,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Overrides:       map[int]int{},
		RemoveBody:      false,
		HeadersToRemove: []string{},
	}
}

// ReturnCodeOverride plugin
type responseCodeOverride struct {
	next            http.Handler
	name            string
	overrrides      map[int]int
	headersToRemove []string
	removeBody      bool
	logger          log.Logger
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &responseCodeOverride{
		overrrides:      config.Overrides,
		next:            next,
		removeBody:      config.RemoveBody,
		headersToRemove: config.HeadersToRemove,
		logger:          *log.New(os.Stdout, "plugin:responseCodeOverride ", log.Ldate|log.Ltime),
	}, nil
}

/*

ResponseWriter has some constraints which are must to know:

* Header() http.Header:

	Changing the header map after a call to WriteHeader (or Write) has no effect unless the modified headers are trailers.

* WriteHeader(statusCode int):

	WriteHeader sends an HTTP response header with the provided status code.

* Write([]byte) (int, error):

	If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK) before writing the data.

*/

type responseCodeOverrideWriter struct {
	http.ResponseWriter

	buffer     bytes.Buffer
	overridden bool
	plugin     *responseCodeOverride //back reference to access the configuration
}

func (w *responseCodeOverrideWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *responseCodeOverrideWriter) WriteHeader(responseCode int) {

	// We need to remove specific headers accoording to
	// RFC specification
	if w.plugin.removeBody {
		w.Header().Del("Content-Length")
		w.Header().Del("Content-Type")

	}

	for _, headerToRemove := range w.plugin.headersToRemove {
		w.Header().Del(headerToRemove)
	}

	if val, ok := w.plugin.overrrides[responseCode]; ok {
		w.overridden = true
		w.ResponseWriter.WriteHeader(val)
		return
	}

	w.ResponseWriter.WriteHeader(responseCode)
	return
}

func (w *responseCodeOverrideWriter) Write(b []byte) (int, error) {
	if !w.overridden {
		// If not overriden yet and we do have a config for such
		// then we need to override that
		if val, ok := w.plugin.overrrides[200]; ok {
			w.overridden = true
			w.ResponseWriter.WriteHeader(val)
		}
	}

	return w.buffer.Write(b)
}

func (e *responseCodeOverride) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	respCodeOverrideWriter := &responseCodeOverrideWriter{ResponseWriter: rw, plugin: e}
	e.next.ServeHTTP(respCodeOverrideWriter, req)

	bodyBytes := respCodeOverrideWriter.buffer.Bytes()

	if e.removeBody {
		rw.Write(make([]byte, 0))
	} else {
		rw.Write(bodyBytes)
	}

}
