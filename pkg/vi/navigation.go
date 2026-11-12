package vi

import "github.com/maxence-charriere/go-app/v11/pkg/app"

// ShareCurrentURL shares the current page URL.
func ShareCurrentURL(ctx app.Context, e app.Event) {
	e.PreventDefault()
	currentURL := app.Window().Get("document").Get("URL").String()

	shareNotSupported := func() {
		navigator := app.Window().Get("navigator")
		if !navigator.Truthy() || !navigator.Get("clipboard").Truthy() {
			app.Window().Call("alert", "Sharing is not supported on your current device or browser.")
			return
		}

		navigator.Get("clipboard").Call("writeText", currentURL)
		app.Window().Call("alert", "Copied to Clipboard.")
	}

	navigator := app.Window().Get("navigator")
	if currentURL == "" || !navigator.Truthy() {
		shareNotSupported()
		return
	}

	canShare := navigator.Get("canShare")
	if !canShare.Truthy() {
		shareNotSupported()
		return
	}

	shareConfig := map[string]any{
		"url": currentURL,
	}

	if !navigator.Call("canShare", shareConfig).Truthy() {
		shareNotSupported()
		return
	}

	navigator.Call("share", shareConfig)
}

// ScrollToTop scrolls to the top of the page.
func ScrollToTop(ctx app.Context) {
	ctx.ScrollTo(PageTopID)
}

// ScrollToFragment scrolls to the element identified by the current page URL's
// fragment, or to the top of the page if the URL has no fragment.
func ScrollToFragment(ctx app.Context) {
	if fragment := ctx.Page().URL().Fragment; fragment != "" {
		ctx.ScrollTo(fragment)
		return
	}
	ScrollToTop(ctx)
}

// NavigateBack navigates to the previous page in the browser's history.
func NavigateBack(ctx app.Context, e app.Event) {
	e.PreventDefault()
	e.StopImmediatePropagation()
	app.Window().Get("history").Call("back")
}

// CanNavigateBack reports whether the browser can navigate to the previous page.
func CanNavigateBack() bool {
	navigation := app.Window().Get("navigation")
	if !navigation.Truthy() {
		return false
	}
	return navigation.Get("canGoBack").Truthy()
}
