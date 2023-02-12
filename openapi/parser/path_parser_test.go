package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/morozovcookie/ogen/jsonschema"
	"github.com/morozovcookie/ogen/openapi"
)

func TestPathParser(t *testing.T) {
	var (
		bar = &openapi.Parameter{
			Name:   "bar",
			Schema: &jsonschema.Schema{Type: jsonschema.Integer},
			In:     openapi.LocationPath,
		}
		baz = &openapi.Parameter{
			Name:   "baz",
			Schema: &jsonschema.Schema{Type: jsonschema.String},
			In:     openapi.LocationPath,
		}
	)

	tests := []struct {
		Path      string
		Params    []*openapi.Parameter
		Expect    openapi.Path
		ExpectErr string
	}{
		{
			Path:   "/foo/{bar}",
			Params: []*openapi.Parameter{bar},
			Expect: openapi.Path{
				{Raw: "/foo/"},
				{Param: bar},
			},
		},
		{
			Path:   "/foo.{bar}",
			Params: []*openapi.Parameter{bar},
			Expect: openapi.Path{
				{Raw: "/foo."},
				{Param: bar},
			},
		},
		{
			Path:   "/foo.{bar}.{baz}abc/def",
			Params: []*openapi.Parameter{bar, baz},
			Expect: openapi.Path{
				{Raw: "/foo."},
				{Param: bar},
				{Raw: "."},
				{Param: baz},
				{Raw: "abc/def"},
			},
		},
		{
			Path: "/foo%20bar",
			Expect: openapi.Path{
				{Raw: "/foo%20bar"},
			},
		},
		{
			Path: "/foo%3fbar",
			Expect: openapi.Path{
				{Raw: "/foo%3Fbar"},
			},
		},
		{
			Path: "/foo%3Fbar",
			Expect: openapi.Path{
				{Raw: "/foo%3Fbar"},
			},
		},
		{
			Path: "/user/ern%41do",
			Expect: openapi.Path{
				{Raw: "/user/ernAdo"},
			},
		},
		{
			Path: "/user/ern%61do",
			Expect: openapi.Path{
				{Raw: "/user/ernado"},
			},
		},
		{
			Path:      "\xF1",
			ExpectErr: errInvalidPathUTF8.Error(),
		},
		{
			Path:      "\x00",
			ExpectErr: `parse "\x00": net/url: invalid control character in URL`,
		},
		{
			Path:      "/foo/{bar}/{baz}",
			Params:    []*openapi.Parameter{bar},
			ExpectErr: `path parameter not specified: "baz"`,
		},
		{
			Path:      "/foo/{",
			ExpectErr: `invalid path: /foo/{`,
		},
		{
			Path:      "/foo/{{",
			ExpectErr: `invalid path: /foo/{{`,
		},
		{
			Path:      "/foo/{{}",
			ExpectErr: `invalid path: /foo/{{}`,
		},
		{
			Path:      "/foo/{}}",
			ExpectErr: `invalid path: /foo/{}}`,
		},
		{
			Path:      "/foo/{{}}",
			ExpectErr: `invalid path: /foo/{{}}`,
		},
		{
			Path:      "/foo/{/",
			ExpectErr: `invalid path: /foo/{/`,
		},
		{
			Path:      "foo/",
			ExpectErr: `path MUST begin with a forward slash`,
		},
		{
			Path:      "/foo?k=v",
			ExpectErr: `path MUST NOT contain a query string`,
		},
		{
			Path:      "https://i-want-to-die.org/foo",
			ExpectErr: `path MUST be relative`,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Test%d", i+1), func(t *testing.T) {
			result, err := parsePath(test.Path, test.Params)
			if test.ExpectErr != "" {
				require.EqualError(t, err, test.ExpectErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.Expect, result)
		})
	}
}
