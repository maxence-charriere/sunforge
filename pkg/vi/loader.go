package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/maxence-charriere/go-app/v11/pkg/errors"
)

// A container that displays a loader in place of its content when a loading
// state is set.
type Loader interface {
	app.UI

	// Sets the ID.
	ID(v string) Loader

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Loader

	// Set the title.
	Title(v app.UI) Loader

	// Set the description.
	Description(v string) Loader
	Descriptionf(format string, v ...any) Loader

	// Set the loading error.
	Err(v error) Loader

	// Sets the error svg icon.
	ErrorIconSVG(v string) Loader

	// Expands to all available space
	Expand() Loader
}

func NewLoader() Loader {
	return &loader{
		IErrorIconSVG: RobotDeadSVG,
	}
}

type loader struct {
	app.Compo

	IID           string
	IClass        string
	IDescription  string
	IErrorIconSVG string
	ITitle        app.UI
	IErr          error
	IExpand       bool

	reloading bool
}

func (l *loader) ID(v string) Loader {
	l.IID = v
	return l
}

func (l *loader) Class(v ...string) Loader {
	l.IClass = app.AppendClass(l.IClass, v...)
	return l
}

func (l *loader) Description(v string) Loader {
	l.IDescription = v
	return l
}

func (l *loader) Descriptionf(format string, v ...any) Loader {
	return l.Description(fmt.Sprintf(format, v...))
}

func (l *loader) ErrorIconSVG(v string) Loader {
	if v != "" {
		l.IErrorIconSVG = v
	}
	return l
}

func (l *loader) Err(err error) Loader {
	l.IErr = err
	return l
}

func (l *loader) Title(v app.UI) Loader {
	title := app.FilterUIElems(v)
	if len(title) != 0 {
		l.ITitle = title[0]
	}
	return l
}

func (l *loader) Expand() Loader {
	l.IExpand = true
	return l
}

func (l *loader) Render() app.UI {
	description := l.IDescription
	if l.IErr != nil {
		description = errors.UIError(l.IErr)
	}

	return app.Div().
		Class("vi-loader").
		Class(StyleFromBool("vi-loader-expand", l.IExpand)).
		Class(StyleFromBool("vi-loader-error", l.IErr != nil)).
		Class("vi-unselectable").
		Class(l.IClass).
		OnWheel(l.scroll).
		Body(
			app.Div().
				Class("vi-loader-content").
				OnClick(l.reload).
				Body(
					app.If(l.IErr == nil || l.reloading, func() app.UI {
						return NewSpinner()
					}).Else(func() app.UI {
						return NewSVG().Data(l.IErrorIconSVG)
					}),
					app.Div().Body(
						l.ITitle,
						app.P().Text(description),
						NewMediaObject().
							AsButton().
							Class("vi-loader-reload").
							SVGIcon(ReloadSVG).
							Text("Reload").
							OnClick(l.reload),
					),
				),
		)
}

func (l *loader) scroll(ctx app.Context, e app.Event) {
	if l.IExpand {
		e.PreventDefault()
	}
}

func (l *loader) reload(ctx app.Context, e app.Event) {
	l.reloading = true
	ctx.Defer(func(ctx app.Context) {
		ctx.Reload()
	})
}
