package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Logo interface {
	app.Composer

	ID(string) Logo
	Class(...string) Logo
	Icon(string) Logo
	Iconf(string, ...any) Logo
	Title(string) Logo
	Titlef(string, ...any) Logo
	Version(string) Logo
	Versionf(string, ...any) Logo
}

func NewLogo() Logo {
	return &logo{}
}

type logo struct {
	app.Compo

	IID      string
	IClass   string
	IIcon    string
	ITitle   string
	IVersion string
}

func (l *logo) ID(v string) Logo {
	l.IID = v
	return l
}

func (l *logo) Class(v ...string) Logo {
	l.IClass = app.AppendClass(l.IClass, v...)
	return l
}

func (l *logo) Icon(v string) Logo {
	l.IIcon = v
	return l
}

func (l *logo) Iconf(format string, v ...any) Logo {
	return l.Icon(fmt.Sprintf(format, v...))
}

func (l *logo) Title(v string) Logo {
	l.ITitle = v
	return l
}

func (l *logo) Titlef(format string, v ...any) Logo {
	return l.Title(fmt.Sprintf(format, v...))
}

func (l *logo) Version(v string) Logo {
	l.IVersion = v
	return l
}

func (l *logo) Versionf(format string, v ...any) Logo {
	return l.Version(fmt.Sprintf(format, v...))
}

func (l *logo) Render() app.UI {
	var withIcon string
	if l.IIcon != "" {
		withIcon = "vi-logo-with-icon"
	}

	return app.A().
		ID(l.IID).
		Class("vi-logo").
		Class(withIcon).
		Class(l.IClass).
		Href("/").
		Body(
			app.If(l.IIcon != "", func() app.UI {
				return app.Img().
					Alt(fmt.Sprintf("%s Logo", l.ITitle)).
					Src(l.IIcon)
			}),
			app.Div().Text(l.ITitle),
			app.Div().Text(l.IVersion),
		)
}
