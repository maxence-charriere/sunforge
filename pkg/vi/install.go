package vi

import "github.com/maxence-charriere/go-app/v11/pkg/app"

// InstallPWA is an event handler that shows the app installation prompt.
// It can be used with OnClick to let users install the progressive web app.
func InstallPWA(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}

func NewHowToInstallOnAppleDevices(appName string) app.UI {
	return NewBox().
		Class("vi-install-pwa").
		WithGap().
		Content(
			NewMediaObject().
				AsH2().
				WithHorizontalCentering().
				WithExtraLargeMedia().
				SVGIcon(AppleSVG).
				Textf("Install %s", appName),
			app.P().Textf("Add %s to your Home Screen or Dock for quick access, just like a native app.", appName),
			app.H3().Text("iPhone & iPad"),
			app.Ol().Body(
				app.Li().Body(
					app.Text("Tap the "),
					app.Span().
						Class("accent").
						Text("☰"),
					app.Text(" on the left side of the address bar."),
				),
				app.Li().Body(
					app.Text("Select "),
					app.Span().
						Class("accent").
						Text("Share"),
					app.Text("."),
				),
				app.Li().Body(
					app.Text("Scroll down and tap "),
					app.Span().
						Class("accent").
						Text("Add to Home Screen"),
					app.Text("."),
				),
			),
			app.H3().Text("MAC (SAFARI)"),
			app.Ol().Body(
				app.Li().Body(
					app.Text("Click the "),
					NewSVG().
						Class("accent").
						Data(ShareVariantSVG),
					app.Text(" on the right side of the toolbar."),
				),
				app.Li().Body(
					app.Text("Select "),
					app.Span().
						Class("accent").
						Text("Add to Dock"),
					app.Text("."),
				),
			),
		)
}
