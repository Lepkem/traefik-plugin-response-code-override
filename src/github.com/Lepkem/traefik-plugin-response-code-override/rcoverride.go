package traefik-plugin-response-code-override

// Config the plugin configuration.
type Config struct {
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// ReturnCodeOverride plugin
type RcOverride struct {
	next     http.Handler
	name     string
    // ...
}

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// ...
	return &RcOverride{
		// ...
	}, nil
}

func (e *RcOverride) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// ...
	e.next.ServeHTTP(rw, req)
}