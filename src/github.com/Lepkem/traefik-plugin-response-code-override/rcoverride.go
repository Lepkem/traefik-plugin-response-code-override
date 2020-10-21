package traefik_plugin_response_code_override

import (
	"bytes"
	"context"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Overrides  map[int]int `json:"overrides,omitempty"`
	RemoveBody bool        `json:"remove_body,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Overrides:  map[int]int{},
		RemoveBody: false,
	}
}

// ReturnCodeOverride plugin
type responseCodeOverride struct {
	next       http.Handler
	name       string
	overrrides map[int]int
	removeBody bool
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &responseCodeOverride{
		overrrides: config.Overrides,
		next:       next,
		removeBody: config.RemoveBody,
	}, nil
}

// Type we use to signal that nothing else writes and we will be able to change the status
// code
type responseCodeOverrideWriter struct {
	http.ResponseWriter

	buffer     bytes.Buffer
	overridden bool
	plugin     *responseCodeOverride //back reference to access the configuration
}

func (w *responseCodeOverrideWriter) WriteHeader(responseCode int) {

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
