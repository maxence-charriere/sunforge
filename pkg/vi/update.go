package vi

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type updatePWA struct {
	app.Compo

	IClass   string
	IAppName string
	ILoading bool

	reloading bool
}

func newUpdatePWA(appName string) *updatePWA {
	return &updatePWA{
		IAppName: appName,
	}
}

func (u *updatePWA) Class(v ...string) *updatePWA {
	u.IClass = app.AppendClass(u.IClass, v...)
	return u
}

func (u *updatePWA) Loading(v bool) *updatePWA {
	u.ILoading = v
	return u
}

func (u *updatePWA) Render() app.UI {
	return app.Aside().
		Class("vi-update-pwa").
		Class(u.IClass).
		OnClick(u.updatePWA).
		Body(
			app.If(!u.reloading, func() app.UI {
				return NewSVG().Data(InstallSVG)
			}).Else(func() app.UI {
				return NewSpinner()
			}),
			app.Div().Body(
				app.H2().Textf("%s Update Ready", u.IAppName),
				app.If(!u.reloading, func() app.UI {
					return app.P().Text("Click/Tap to update")
				}).Else(func() app.UI {
					return app.P().Text("Updating...")
				}),
			),
		)
}

func (u *updatePWA) updatePWA(ctx app.Context, e app.Event) {
	u.reloading = true
	ctx.Defer(func(ctx app.Context) {
		ctx.Reload()
	})
}
