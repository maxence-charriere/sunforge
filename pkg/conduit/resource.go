package conduit

import (
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// Resolve rewrites the given path to include the resolver's version when
// applicable.
//
// Paths under /web/ are mapped to /web/versions/ with the version inserted
// before the file extension. Special cases like app.wasm, app.js, and other
// core assets required to run a go-app PWA are also handled explicitly.
//
// The call is delegated to the base ResourceResolver when no rewrite is needed.
func ResourceResolverWithVersion(base app.ResourceResolver, version string, paths ...string) app.ResourceResolver {
	pathMap := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		if strings.HasPrefix(p, "/web/") {
			pathMap[p] = struct{}{}
		}
	}

	return versionedResourceResolver{
		ResourceResolver: base,
		version:          version,
		paths:            pathMap,
	}
}

type versionedResourceResolver struct {
	app.ResourceResolver

	version string
	paths   map[string]struct{}
}

func (r versionedResourceResolver) Resolve(path string) string {
	switch path {
	case "/web/app.wasm":
		path = "/web/versions/app-" + r.version + ".wasm"

	case "/app.js",
		"/wasm_exec.js",
		"/manifest.webmanifest",
		"/app.css",
		"/vi.css":
		path = "/" + r.version + path

	default:
		if _, resolvable := r.paths[path]; resolvable {
			extension := filepath.Ext(path)
			if extension == "" {
				return r.ResourceResolver.Resolve(path)
			}
			path = strings.Replace(path, "/web/", "/web/versions/", 1)
			path = strings.TrimSuffix(path, extension) + "-" + r.version + extension
		}
	}

	return r.ResourceResolver.Resolve(path)
}
