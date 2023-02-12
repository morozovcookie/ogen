package ogen_test

import (
	"net/url"
	"path"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/morozovcookie/ogen"
	"github.com/morozovcookie/ogen/gen"
	"github.com/morozovcookie/ogen/gen/genfs"
	"github.com/morozovcookie/ogen/gen/ir"
	"github.com/morozovcookie/ogen/internal/integration/customformats/phonetype"
	"github.com/morozovcookie/ogen/internal/integration/customformats/rgbatype"
	"github.com/morozovcookie/ogen/internal/location"
	"github.com/morozovcookie/ogen/jsonschema"
	"github.com/morozovcookie/ogen/openapi/parser"
)

func testGenerate(t *testing.T, dir, filename string, data []byte, aliases ctAliases, ignore ...string) {
	t.Helper()
	t.Parallel()
	log := zaptest.NewLogger(t)

	spec, err := ogen.Parse(data)
	require.NoError(t, err)

	notImplemented := map[string]struct{}{}
	opt := gen.Options{
		InferSchemaType:      true,
		IgnoreNotImplemented: ignore,
		NotImplementedHook: func(name string, err error) {
			notImplemented[name] = struct{}{}
		},
		ContentTypeAliases: aliases,
		File:               location.NewFile(filename, filename, data),
		Logger:             log,
	}

	switch filename {
	case "file_reference.yml": // HACK
		opt.AllowRemote = true
		opt.RootURL = &url.URL{
			Scheme: "file",
			Path:   "/" + path.Join(dir, filename),
		}
		opt.Remote = gen.RemoteOptions{
			ReadFile: func(p string) ([]byte, error) {
				p = strings.TrimPrefix(p, "/")
				return testdata.ReadFile(p)
			},
			URLToFilePath: func(u *url.URL) (string, error) {
				// By default, urlpath.URLToFilePath output depends on the OS.
				//
				// But we use virtual filesystem, so we should use the fs.FS path.
				if u.Path == "" {
					return u.Opaque, nil
				}
				return u.Path, nil
			},
		}
	case "custom_formats.json":
		opt.CustomFormats = gen.CustomFormatsMap{
			jsonschema.String: {
				"phone": gen.CustomFormat[
					phonetype.Phone,
					phonetype.JSONPhoneEncoding,
					phonetype.TextPhoneEncoding,
				](),
				"rgba": gen.CustomFormat[
					rgbatype.RGBA,
					rgbatype.JSONRGBAEncoding,
					rgbatype.TextRGBAEncoding,
				](),
			},
		}
	}

	t.Run("Gen", func(t *testing.T) {
		defer func() {
			if rr := recover(); rr != nil {
				t.Fatalf("panic: %+v\n%s", rr, debug.Stack())
			}
		}()

		g, err := gen.NewGenerator(spec, opt)
		require.NoError(t, err)
		require.NoError(t, g.WriteSource(genfs.CheckFS{}, "api"))

		if len(opt.IgnoreNotImplemented) > 0 {
			// Check that all ignore rules are necessary.
			for _, feature := range ignore {
				if _, ok := notImplemented[feature]; !ok {
					t.Errorf("Ignore rule %q hasn't been used", feature)
				}
			}
		}
	})
	t.Run("Full", func(t *testing.T) {
		t.Skipf("Ignoring: [%s]", strings.Join(opt.IgnoreNotImplemented, ", "))
	})
}

type ctAliases = map[string]ir.Encoding

func runPositive(root string,
	aliases map[string]ctAliases,
	skipSets map[string][]string,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		// Ensure that all skipSets schemas are present.
		for file := range skipSets {
			_, err := testdata.ReadFile(path.Join(root, file))
			require.NoErrorf(t, err, "skip file %s", file)
		}

		walkTestdata(t, root, func(t *testing.T, file string, data []byte) {
			dir := path.Dir(file)
			if parent := path.Base(dir); parent == "file_reference_external" {
				t.Skip("Special directory for testing remote references.")
				return
			}

			file = strings.TrimPrefix(file, root+"/")
			skip := skipSets[file]
			testGenerate(t, dir, file, data, aliases[file], skip...)
		})
	}
}

func TestGenerate(t *testing.T) {
	t.Run("Positive", runPositive("_testdata/positive", nil,
		map[string][]string{
			"sample.json": {
				"enum format",
			},
			"content_header_response.json": {
				"parameter content encoding",
			},
		}))

	t.Run("Examples", runPositive("_testdata/examples",
		map[string]ctAliases{
			"autorest/ApiManagementClient-openapi.json": {
				"text/json":                        ir.EncodingJSON,
				"application/vnd.swagger.doc+json": ir.EncodingJSON,
			},
			"api.github.com.json": {
				"text/x-markdown":            ir.EncodingTextPlain,
				"text/html":                  ir.EncodingTextPlain,
				"application/octocat-stream": ir.EncodingTextPlain,
				// FIXME(tdakkota): multiple response types makes wrapper cry about
				// 	type name conflict.
				// "application/vnd.github.v3.star+json": ir.EncodingJSON,
				"application/vnd.github.v3.object": ir.EncodingJSON,
				"application/scim+json":            ir.EncodingJSON,
			},
			"k8s.json": {
				"application/jwk-set+json":               ir.EncodingJSON,
				"application/merge-patch+json":           ir.EncodingJSON,
				"application/strategic-merge-patch+json": ir.EncodingJSON,
			},
		},
		map[string][]string{
			"autorest/additionalProperties.json": {},
			"autorest/ApiManagementClient-openapi.json": {
				"oauth2 security",
			},
			"autorest/lro.json":                 {},
			"autorest/storage.json":             {},
			"autorest/xms-error-responses.json": {},
			"2ch.yml":                           {},
			"api.github.com.json": {
				"complex anyOf",
				"discriminator inference",
				"sum types with same names",
				"sum type parameter",
				"array defaults",
			},
			"manga.json":            {},
			"telegram_bot_api.json": {},
			"gotd_bot_api.json":     {},
			"k8s.json": {
				"unsupported content types",
			},
			"petstore-expanded.yml": {},
			"redoc/discriminator.json": {
				"unsupported content types",
			},
		}))
}

func TestNegative(t *testing.T) {
	walkTestdata(t, "_testdata/negative", func(t *testing.T, file string, data []byte) {
		a := require.New(t)
		_, name := path.Split(file)

		spec, err := ogen.Parse(data)
		a.NoError(err)

		f := location.NewFile(name, name, data)
		_, err = parser.Parse(spec, parser.Settings{
			InferTypes: true,
			File:       f,
		})
		a.NoError(err, "If the error is related to parser, move this test to parser package testdata")

		_, err = gen.NewGenerator(spec, gen.Options{
			InferSchemaType: true,
			File:            f,
		})
		a.Error(err)

		var buf strings.Builder
		if location.PrintPrettyError(&buf, true, err) {
			t.Log(buf.String())
		} else {
			t.Logf("%+v", err)
		}
	})
}
