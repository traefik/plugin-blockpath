package plugin_blockpath

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc             string
		regexps          []string
		regexpsWhitelist []string
		expErr           bool
	}{
		{
			desc:             "should return no error",
			regexps:          []string{`^/foo/(.*)`},
			regexpsWhitelist: []string{`^/foo/(.*)`},
			expErr:           false,
		},
		{
			desc:             "should return an error",
			regexps:          []string{"*"},
			regexpsWhitelist: []string{"*"},
			expErr:           true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Regex:          test.regexps,
				RegexWhitelist: test.regexpsWhitelist,
			}

			if _, err := New(context.Background(), nil, cfg, "name"); test.expErr && err == nil {
				t.Errorf("expected error on bad regexp format")
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc             string
		regexps          []string
		regexpsWhitelist []string
		reqPath          string
		expNextCall      bool
		expStatusCode    int
	}{
		{
			desc:          "Should return forbidden status",
			regexps:       []string{"/test"},
			reqPath:       "/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{"/test", "/toto"},
			reqPath:       "/toto",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:             "should return forbidden status",
			regexps:          []string{"/test", "/toto"},
			regexpsWhitelist: []string{"/tests", "/totos"},
			reqPath:          "/toto",
			expNextCall:      false,
			expStatusCode:    http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexps:       []string{"/test", "/toto"},
			reqPath:       "/plop",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:             "should return ok status",
			regexps:          []string{"^/wp-admin(.*)"},
			regexpsWhitelist: []string{"^/wp-admin/admin-ajax\\.php(.*)"},
			reqPath:          "/wp-admin/admin-ajax.php",
			expNextCall:      true,
			expStatusCode:    http.StatusOK,
		},
		{
			desc:          "should return ok status",
			reqPath:       "/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{`^/bar(.*)`},
			reqPath:       "/bar/foo",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexps:       []string{`^/bar(.*)`},
			reqPath:       "/foo/bar",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Regex:          test.regexps,
				RegexWhitelist: test.regexpsWhitelist,
			}

			nextCall := false
			next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				nextCall = true
			})

			handler, err := New(context.Background(), next, cfg, "blockpath")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("http://localhost%s", test.reqPath)
			req := httptest.NewRequest(http.MethodGet, url, nil)

			handler.ServeHTTP(recorder, req)

			if nextCall != test.expNextCall {
				t.Errorf("next handler should not be called")
			}

			if recorder.Result().StatusCode != test.expStatusCode {
				t.Errorf("%s: got status code %d, want %d", test.desc, recorder.Code, test.expStatusCode)
			}
		})
	}
}
