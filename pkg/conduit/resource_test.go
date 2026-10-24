package conduit

import (
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/stretchr/testify/require"
)

func TestResourceResolverWithVersion(t *testing.T) {
	version := "v0.0.0"
	tests := []struct {
		scenario string
		in       string
		out      string
	}{
		{
			scenario: "not a resource",
			in:       "/hello",
			out:      "/hello",
		},
		{
			scenario: "app.wasm",
			in:       "/web/app.wasm",
			out:      "/web/versions/app-" + version + ".wasm",
		},
		{
			scenario: "go-app required file",
			in:       "/wasm_exec.js",
			out:      "/" + version + "/wasm_exec.js",
		},
		{
			scenario: "vi required file",
			in:       "/vi.css",
			out:      "/" + version + "/vi.css",
		},
		{
			scenario: "specified file resource",
			in:       "/web/murlok.css",
			out:      "/web/versions/murlok-" + version + ".css",
		},
		{
			scenario: "specified non-file resource",
			in:       "/web/murlok",
			out:      "/web/murlok",
		},
	}

	for _, u := range tests {
		t.Run(u.scenario, func(t *testing.T) {
			r := ResourceResolverWithVersion(app.LocalDir(""), version,
				"/web/murlok.css",
				"/web/murlok",
			)
			require.Equal(t, u.out, r.Resolve(u.in))
		})
	}
}
