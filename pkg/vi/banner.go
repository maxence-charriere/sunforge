package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Banner interface {
	app.Composer

	WithFullHeight() Banner
	WithImageAlignedTop() Banner
	ID(v string) Banner
	Class(v ...string) Banner
	Alt(v string) Banner
	Altf(format string, v ...any) Banner
	Image(v string) Banner
	Imagef(format string, v ...any) Banner
	Media(v app.UI) Banner
	Subtitle(v string) Banner
	Subtitlef(format string, v ...any) Banner
	Content(v ...app.UI) Banner
	Nav(v ...app.UI) Banner
}

func NewBanner(title string) Banner {
	return &banner{ITitle: title}
}

func NewBannerf(format string, v ...any) Banner {
	return NewBanner(fmt.Sprintf(format, v...))
}

type banner struct {
	app.Compo

	IWithFullHeight      bool
	IWithImageAlignedTop bool
	IID                  string
	IClass               string
	ITitle               string
	IAlt                 string
	IMedia               app.UI
	IContent             []app.UI
	INav                 []app.UI
}

func (b *banner) WithFullHeight() Banner {
	b.IWithFullHeight = true
	return b
}

func (b *banner) WithImageAlignedTop() Banner {
	b.IWithImageAlignedTop = true
	return b
}

func (b *banner) ID(v string) Banner {
	b.IID = v
	return b
}

func (b *banner) Class(v ...string) Banner {
	b.IClass = app.AppendClass(b.IClass, v...)
	return b
}

func (b *banner) Alt(v string) Banner {
	b.IAlt = v
	if img, ok := b.IMedia.(app.HTMLImg); ok {
		b.IMedia = img.Alt(b.IAlt)
	}
	return b
}

func (b *banner) Altf(format string, v ...any) Banner {
	return b.Alt(fmt.Sprintf(format, v...))
}

func (b *banner) Image(v string) Banner {
	return b.Media(app.Img().
		Alt(b.IAlt).
		Src(v))
}

func (b *banner) Imagef(format string, v ...any) Banner {
	return b.Image(fmt.Sprintf(format, v...))
}

func (b *banner) Media(v app.UI) Banner {
	if m := app.FilterUIElems(v); len(m) != 0 {
		b.IMedia = m[0]
	}
	return b
}

func (b *banner) Subtitle(v string) Banner {
	return b.Content(app.Text(v))
}

func (b *banner) Subtitlef(format string, v ...any) Banner {
	return b.Subtitle(fmt.Sprintf(format, v...))
}

func (b *banner) Content(v ...app.UI) Banner {
	b.IContent = app.FilterUIElems(v...)
	return b
}

func (b *banner) Nav(v ...app.UI) Banner {
	b.INav = app.FilterUIElems(v...)
	return b
}

func (b *banner) Render() app.UI {
	return app.Header().
		ID(b.IID).
		Class("vi-banner").
		Class(IfTrue(b.IWithFullHeight, "vi-banner-full")).
		Class(IfTrue(b.IWithImageAlignedTop, "vi-banner-image-top")).
		Class(b.IClass).
		Body(
			app.Div().Body(b.IMedia),
			app.Div().Body(
				app.H1().Text(b.ITitle),
				app.If(len(b.IContent) != 0, func() app.UI {
					return app.Div().Body(b.IContent...)
				}),
				app.If(len(b.INav) != 0, func() app.UI {
					return NewBox().
						AsUl().
						WithColumns().
						WithoutBackground().
						WithoutPadding().
						Class("p").
						Content(b.INav...)
				}),
			),
		)
}
