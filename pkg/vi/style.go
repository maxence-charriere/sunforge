package vi

import (
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// Disabled returns the "disabled" class name if v is true.
func Disabled(v bool) string {
	return StyleFromBool("disabled", v)
}

// StyleFromBool returns className if v is true, otherwise an empty string.
func StyleFromBool(className string, v bool) string {
	return IfTrue(v, className)
}

// StyleFromURLPath returns className when the current URL path matches
// or starts with the given path.
func StyleFromURLPath(className string, path string) string {
	pagePath := app.Window().URL().Path
	if pagePath == "/" {
		return StyleFromBool(className, pagePath == path)
	}
	return StyleFromBool(className, strings.HasPrefix(pagePath, path))
}

// StyleFromURLFragment returns className when the current URL fragment
// matches the given fragment.
func StyleFromURLFragment(className string, fragment string) string {
	return StyleFromBool(className, app.Window().URL().Fragment == fragment)
}
