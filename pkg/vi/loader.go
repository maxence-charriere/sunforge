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

	// Set whether the the loader is displayed.
	Loading(v bool) Loader

	// Set the description.
	Description(v string) Loader
	Descriptionf(format string, v ...any) Loader

	// Set the loading error.
	Err(v error) Loader

	// Sets the content.
	Content(v app.UI) Loader

	// Sets the icon. Default is a spinner.
	Icon(v app.UI) Loader

	// Sets the error icon. Default is a spinner.
	ErrorIcon(v app.UI) Loader

	// Set the title.
	Title(v app.UI) Loader

	// Expands to all available space
	Expand() Loader
}

func NewLoader() Loader {
	return &loader{ILoading: true}
}

type loader struct {
	app.Compo

	IID          string
	IClass       string
	IDescription string
	ILoading     bool
	IErr         error
	IContent     app.UI
	IIcon        app.UI
	IErrorIcon   app.UI
	ITitle       app.UI
	IExpand      bool

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

func (l *loader) Loading(v bool) Loader {
	l.ILoading = v
	return l
}

func (l *loader) Description(v string) Loader {
	l.IDescription = v
	return l
}

func (l *loader) Descriptionf(format string, v ...any) Loader {
	return l.Description(fmt.Sprintf(format, v...))
}

func (l *loader) Err(err error) Loader {
	l.IErr = err
	return l
}

func (l *loader) Content(v app.UI) Loader {
	content := app.FilterUIElems(v)
	if len(content) != 0 {
		l.IContent = content[0]
	}
	return l
}

func (l *loader) Icon(v app.UI) Loader {
	icon := app.FilterUIElems(v)
	if len(icon) != 0 {
		l.IIcon = icon[0]
	}
	return l
}

func (l *loader) ErrorIcon(v app.UI) Loader {
	icon := app.FilterUIElems(v)
	if len(icon) != 0 {
		l.IErrorIcon = icon[0]
	}
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
	if l.ILoading || l.IErr != nil || l.reloading {
		var expand string
		if l.IExpand {
			expand = "vi-loader-expand"
		}

		return app.Div().
			DataSet("vi", "loader").
			Class("vi-loader").
			Class(expand).
			Class(l.IClass).
			Body(
				app.If(l.IErr == nil || l.reloading, func() app.UI {
					if l.IIcon != nil {
						return l.IIcon
					}
					return NewSpinner()
				}).Else(func() app.UI {
					if l.IErrorIcon != nil {
						return l.IErrorIcon
					}
					return NewSVG().Data(DefaultErrorSVG)
				}),

				app.Div().Body(
					l.ITitle,
					app.If(l.IErr == nil, func() app.UI {
						return app.P().Text(l.IDescription)
					}).ElseSlice(func() []app.UI {
						return []app.UI{
							app.P().
								Class("vi-loader-content-error").
								Text(errors.UIError(l.IErr)),
							NewBox().
								WithoutBackground().
								WithColumns().
								Class("vi-loader-error-controls").
								Class(Disabled(l.reloading)).
								Content(
									NewMediaObject().
										AsButton().
										SVGIcon(ArrowLeftSVG).
										Text("Back").
										OnClick(NavigateBack),
									NewMediaObject().
										AsButton().
										SVGIcon(ReloadSVG).
										Text("Reload").
										OnClick(l.reload),
								),
						}
					}),
				),
			)
	}
	return l.IContent
}

func (l *loader) reload(ctx app.Context, e app.Event) {
	l.reloading = true
	ctx.Defer(func(ctx app.Context) {
		ctx.Reload()
	})
}
