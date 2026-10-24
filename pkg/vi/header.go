package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// A header that contains the page main title.
type Header interface {
	app.UI

	// Sets the ID.
	ID(v string) Header

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Header

	// Sets the alt.
	Alt(v string) Header
	Altf(format string, v ...any) Header

	// Sets the media element to an image with the given URL.
	Image(v string) Header
	Imagef(format string, v ...any) Header

	// Sets the media element.
	Media(v app.UI) Header

	// Sets the title with the given text.
	Title(v string) Header
	Titlef(format string, v ...any) Header

	// Sets the subtitle with the given text.
	Subtitle(v string) Header
	Subtitlef(format string, v ...any) Header

	// Sets the sub content with the given ui element.
	SubContent(v app.UI) Header
}

// Creates a header.
func NewHeader() Header {
	return &header{}
}

type header struct {
	app.Compo

	IID         string
	IClass      string
	ITitle      string
	ISubContent app.UI
	IAlt        string
	IMedia      app.UI
}

func (h *header) ID(v string) Header {
	h.IID = v
	return h
}

func (h *header) Class(v ...string) Header {
	h.IClass = app.AppendClass(h.IClass, v...)
	return h
}

func (h *header) Alt(v string) Header {
	h.IAlt = v
	if img, ok := h.IMedia.(app.HTMLImg); ok && h.IAlt != "" {
		h.IMedia = img.Alt(h.IAlt)
	}
	return h
}

func (h *header) Altf(format string, v ...any) Header {
	return h.Alt(fmt.Sprintf(format, v...))
}

func (h *header) Image(v string) Header {
	h.IMedia = app.Img().
		Alt(h.IAlt).
		Height(261).
		Src(v)
	return h
}

func (h *header) Imagef(format string, v ...any) Header {
	return h.Image(fmt.Sprintf(format, v...))
}

func (h *header) Media(v app.UI) Header {
	h.IMedia = v
	return h
}

func (h *header) Title(v string) Header {
	h.ITitle = v
	return h
}

func (h *header) Titlef(format string, v ...any) Header {
	return h.Title(fmt.Sprintf(format, v...))
}

func (h *header) Subtitle(v string) Header {
	h.ISubContent = app.P().Text(v)
	return h
}

func (h *header) Subtitlef(format string, v ...any) Header {
	return h.Subtitle(fmt.Sprintf(format, v...))
}

func (h *header) SubContent(v app.UI) Header {
	if subcontent := app.FilterUIElems(v); len(subcontent) != 0 {
		h.ISubContent = subcontent[0]
	}
	return h
}

func (h *header) Render() app.UI {
	return app.Header().
		ID(h.IID).
		Class("vi-header").
		Class(h.IClass).
		Body(
			h.IMedia,
			app.Div().Body(
				app.H1().Text(h.ITitle),
				h.ISubContent,
			),
		)
}
