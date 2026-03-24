// Package plugin_blockpath a plugin to block a path.
package plugin_blockpath

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
)

// Config holds the plugin configuration.
type Config struct {
	Regex      []string `json:"regex,omitempty"`
	StatusCode int
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Regex:      make([]string, 0),
		StatusCode: http.StatusForbidden,
	}
}

type blockPath struct {
	name       string
	next       http.Handler
	regexps    []*regexp.Regexp
	statuscode int
}

// New creates and returns a plugin instance.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	regexps := make([]*regexp.Regexp, len(config.Regex))

	for regexIndex, regex := range config.Regex {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", regex, err)
		}

		regexps[regexIndex] = re
	}

	return &blockPath{
		name:       name,
		next:       next,
		regexps:    regexps,
		statuscode: config.StatusCode,
	}, nil
}

func (b *blockPath) ServeHTTP(responseWriter http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.EscapedPath()

	for _, re := range b.regexps {
		if re.MatchString(currentPath) {
			responseWriter.WriteHeader(b.statuscode)
			return
		}
	}

	b.next.ServeHTTP(responseWriter, req)
}
