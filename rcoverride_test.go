package traefik_plugin_response_code_override

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverrideWithBody(t *testing.T) {
	expectedBody := "This is a test body"
	expectedResponseCode := http.StatusTooManyRequests

	r, _ := http.NewRequest("GET", "https://google.com/whoami", nil)
	w := httptest.NewRecorder()

	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	})

	plugin := &responseCodeOverride{
		overrrides: map[int]int{200: expectedResponseCode},
		next:       nextHandler,
		removeBody: false,
	}

	// create the handler to test, using our custom "next" handler
	plugin.ServeHTTP(w, r)

	assert.Equal(t, expectedResponseCode, w.Code, "handler returned wrong status code")
	assert.Equal(t, expectedBody, w.Body.String(), "Body should be returned")
}

func TestOverrideWithoutBody(t *testing.T) {
	expectedBody := ""
	expectedResponseCode := http.StatusTooManyRequests

	r, _ := http.NewRequest("GET", "https://google.com/whoami", nil)
	w := httptest.NewRecorder()

	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Test", "test")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	})

	plugin := &responseCodeOverride{
		overrrides: map[int]int{200: expectedResponseCode},
		next:       nextHandler,
		removeBody: true,
	}

	// create the handler to test, using our custom "next" handler
	plugin.ServeHTTP(w, r)

	assert.Equal(t, expectedResponseCode, w.Code, "handler returned wrong status code")
	assert.Equal(t, expectedBody, w.Body.String(), "Body should be returned")

}

func TestOverrideWithoutHeader(t *testing.T) {
	expectedBody := ""
	expectedResponseCode := http.StatusTooManyRequests

	clientResponseHeader := "X-Test"
	clientResponseHeaderVal := "test"

	r, _ := http.NewRequest("GET", "https://google.com/whoami", nil)
	w := httptest.NewRecorder()

	// create a handler to use as "next" which will verify the request
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(clientResponseHeader, clientResponseHeaderVal)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expectedBody))
	})

	plugin := &responseCodeOverride{
		overrrides:      map[int]int{200: expectedResponseCode},
		next:            nextHandler,
		removeBody:      false,
		headersToRemove: []string{clientResponseHeader},
	}

	// create the handler to test, using our custom "next" handler
	plugin.ServeHTTP(w, r)

	resp := w.Result()

	assert.Equal(t, expectedResponseCode, w.Code, "handler returned wrong status code")
	assert.Equal(t, expectedBody, w.Body.String(), "Body should be returned")
	assert.Nil(t, resp.Header[clientResponseHeader], "Header should be removed from response")

}
