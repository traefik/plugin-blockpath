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
	Regex          []string `json:"regex,omitempty"`
	RegexWhitelist []string `json:"regexwhitelist,omitempty"`
}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

type blockPath struct {
	name             string
	next             http.Handler
	regexps          []*regexp.Regexp
	regexpsWhitelist []*regexp.Regexp
}

// New creates and returns a plugin instance. Ensure that all provided regex strings compile correctly.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	regexps := make([]*regexp.Regexp, len(config.Regex))
	regexpsWhitelist := make([]*regexp.Regexp, len(config.RegexWhitelist))

	for i, regex := range config.Regex {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", regex, err)
		}

		regexps[i] = re
	}

	for i, regex := range config.RegexWhitelist {
		re, err := regexp.Compile(regex)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex %q: %w", regex, err)
		}

		regexpsWhitelist[i] = re
	}

	return &blockPath{
		name:             name,
		next:             next,
		regexps:          regexps,
		regexpsWhitelist: regexpsWhitelist,
	}, nil
}

func (b *blockPath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.EscapedPath()
	isBlocked := false

	// Check if the request should be blocked
	for _, re := range b.regexps {
		if re.MatchString(currentPath) {
			isBlocked = true
			break
		}
	}

	// Only check for whitelist if the request was blocked
	if isBlocked {
		for _, re := range b.regexpsWhitelist {
			if re.MatchString(currentPath) {
				isBlocked = false
				break
			}
		}
	}

	// If still blocked, send a forbidden status
	if isBlocked {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	b.next.ServeHTTP(rw, req)
}
