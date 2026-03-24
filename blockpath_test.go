package plugin_blockpath

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc    string
		regexps []string
		expErr  bool
	}{
		{
			desc:    "should return no error",
			regexps: []string{`^/foo/(.*)`},
			expErr:  false,
		},
		{
			desc:    "should return an error",
			regexps: []string{"*"},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Regex:      test.regexps,
				StatusCode: http.StatusForbidden,
			}

			_, err := New(context.Background(), nil, cfg, "name")
			if test.expErr && err == nil {
				t.Errorf("expected error on bad regexp format")
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc          string
		regexps       []string
		reqPath       string
		statusCode    int
		expNextCall   bool
		expStatusCode int
	}{
		{
			desc:          "should return forbidden status",
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
			desc:          "should return ok status",
			regexps:       []string{"/test", "/toto"},
			reqPath:       "/plop",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
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
		{
			desc:          "should return defined statuscode not found",
			regexps:       []string{`^/bar(.*)`},
			reqPath:       "/bar/foo",
			statusCode:    http.StatusNotFound,
			expNextCall:   false,
			expStatusCode: http.StatusNotFound,
		},
		{
			desc:          "should return ok status with not matching and statuscode set",
			regexps:       []string{"/test", "/toto"},
			reqPath:       "/plop",
			statusCode:    http.StatusNotFound,
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status with no regex and statuscode set",
			reqPath:       "/test",
			statusCode:    http.StatusForbidden,
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := CreateConfig()

			cfg.Regex = test.regexps
			if test.statusCode != 0 {
				cfg.StatusCode = test.statusCode
			}

			nextCall := false
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				nextCall = true
			})

			handler, err := New(context.Background(), next, cfg, "blockpath")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			url := "http://localhost" + test.reqPath
			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)

			handler.ServeHTTP(recorder, req)

			if nextCall != test.expNextCall {
				t.Errorf("next handler should not be called")
			}

			if recorder.Result().StatusCode != test.expStatusCode {
				t.Errorf("got status code %d, want %d", recorder.Code, test.expStatusCode)
			}
		})
	}
}
