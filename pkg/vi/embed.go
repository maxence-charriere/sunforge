//go:build !wasm

package vi

import _ "embed"

//go:embed vi.css
var viCSS string

func init() {
	styles = viCSS
}
