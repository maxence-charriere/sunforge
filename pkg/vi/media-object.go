package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// MediaObject is the interface that describes a displaying media
// (e.g., an icon) along with related text content.
type MediaObject interface {
	app.Composer

	AsH1() MediaObject
	AsH2() MediaObject
	AsH3() MediaObject
	AsHeader() MediaObject
	AsA() MediaObject
	AsLi() MediaObject
	InLi() MediaObject
	AsButton() MediaObject
	AsFigCaption() MediaObject
	WithHorizontalCentering() MediaObject
	WithMinSize() MediaObject
	WithBackground() MediaObject
	WithMediaRounded() MediaObject
	WithMediaBorder() MediaObject
	WithSmallMedia() MediaObject
	WithNormalMedia() MediaObject
	WithLargeMedia() MediaObject
	WithExtraLargeMedia() MediaObject
	WithControls() MediaObject
	WithoutMediaBorder() MediaObject
	Vertical() MediaObject
	ID(v string) MediaObject
	Class(v string) MediaObject
	SVGIcon(v string) MediaObject
	Image(v string) MediaObject
	Imagef(format string, v ...any) MediaObject
	Media(v app.UI) MediaObject
	Help(v string) MediaObject
	Helpf(format string, v ...any) MediaObject
	Href(v string) MediaObject
	Hreff(format string, v ...any) MediaObject
	Rel(v string) MediaObject
	Relf(format string, v ...any) MediaObject
	Alt(v string) MediaObject
	Altf(format string, v ...any) MediaObject
	Lazy(v bool) MediaObject
	Order(v int) MediaObject
	Text(v string) MediaObject
	Textf(format string, v ...any) MediaObject
	Content(v ...app.UI) MediaObject

	OnClick(v app.EventHandler) MediaObject
}

// Creates a media object.
func NewMediaObject() MediaObject {
	return &mediaObject{
		IRoot: "div",
	}
}

type mediaObject struct {
	app.Compo

	IRoot                    string
	IOuterElem               string
	IWithHorizontalCentering bool
	IWithMinSize             bool
	IWithBackground          bool
	IWithMediaBorder         bool
	IWithMediaRounded        bool
	IWithControls            bool
	IWithoutMediaBorder      bool
	IVertical                bool
	IMediaSize               string
	IID                      string
	IClass                   string
	IMedia                   app.UI
	IHelp                    string
	IHref                    string
	IRel                     string
	IAlt                     string
	ILazy                    bool
	IOrder                   int
	IContent                 []app.UI
	IOnClick                 app.EventHandler
}

func (o *mediaObject) AsH1() MediaObject {
	o.IRoot = "h1"
	return o
}

func (o *mediaObject) AsH2() MediaObject {
	o.IRoot = "h2"
	return o
}

func (o *mediaObject) AsH3() MediaObject {
	o.IRoot = "h3"
	return o
}

func (o *mediaObject) AsHeader() MediaObject {
	o.IRoot = "header"
	return o
}

func (o *mediaObject) AsA() MediaObject {
	o.IRoot = "a"
	return o
}

func (o *mediaObject) AsLi() MediaObject {
	o.IRoot = "li"
	return o
}

func (o *mediaObject) InLi() MediaObject {
	o.IOuterElem = "li"
	return o
}

func (o *mediaObject) AsButton() MediaObject {
	o.IRoot = "button"
	return o
}

func (o *mediaObject) AsFigCaption() MediaObject {
	o.IRoot = "figcaption"
	return o
}

func (o *mediaObject) WithHorizontalCentering() MediaObject {
	o.IWithHorizontalCentering = true
	return o
}

func (o *mediaObject) WithMinSize() MediaObject {
	o.IWithMinSize = true
	return o
}

func (o *mediaObject) WithBackground() MediaObject {
	o.IWithBackground = true
	return o
}

func (o *mediaObject) WithMediaRounded() MediaObject {
	o.IWithMediaRounded = true
	return o
}

func (o *mediaObject) WithMediaBorder() MediaObject {
	o.IWithMediaBorder = true
	return o
}

func (o *mediaObject) WithSmallMedia() MediaObject {
	o.IMediaSize = "vi-media-object-with-media-small"
	return o
}

func (o *mediaObject) WithNormalMedia() MediaObject {
	o.IMediaSize = "vi-media-object-with-media-normal"
	return o
}

func (o *mediaObject) WithLargeMedia() MediaObject {
	o.IMediaSize = "vi-media-object-with-media-large"
	return o
}

func (o *mediaObject) WithExtraLargeMedia() MediaObject {
	o.IMediaSize = "vi-media-object-with-media-xlarge"
	return o
}

func (o *mediaObject) WithControls() MediaObject {
	o.IWithControls = true
	return o
}

func (o *mediaObject) WithoutMediaBorder() MediaObject {
	o.IWithoutMediaBorder = true
	return o
}

func (o *mediaObject) Vertical() MediaObject {
	o.IVertical = true
	return o
}

func (o *mediaObject) ID(v string) MediaObject {
	o.IID = v
	return o
}

func (o *mediaObject) Class(v string) MediaObject {
	o.IClass = app.AppendClass(o.IClass, v)
	return o
}

func (o *mediaObject) SVGIcon(v string) MediaObject {
	o.IMedia = NewSVG().Data(v)
	return o
}

func (o *mediaObject) Image(v string) MediaObject {
	img := app.Img().
		Alt(o.IAlt).
		Src(v)
	if o.ILazy {
		img = img.Loading("lazy")
	}
	o.IMedia = img
	return o
}

func (o *mediaObject) Imagef(format string, v ...any) MediaObject {
	return o.Image(fmt.Sprintf(format, v...))
}

func (o *mediaObject) Media(v app.UI) MediaObject {
	if media := app.FilterUIElems(v); len(media) != 0 {
		o.IMedia = media[0]
	}
	return o
}

func (o *mediaObject) Alt(v string) MediaObject {
	o.IAlt = v
	if img, ok := o.IMedia.(app.HTMLImg); ok {
		o.IMedia = img.Alt(o.IAlt)
	}
	return o
}

func (o *mediaObject) Altf(format string, v ...any) MediaObject {
	return o.Alt(fmt.Sprintf(format, v...))
}

func (o *mediaObject) Lazy(v bool) MediaObject {
	o.ILazy = v
	if img, ok := o.IMedia.(app.HTMLImg); ok && o.ILazy {
		o.IMedia = img.Loading("lazy")
	}
	return o
}

func (o *mediaObject) Help(v string) MediaObject {
	o.IHelp = v
	return o
}

func (o *mediaObject) Helpf(format string, v ...any) MediaObject {
	return o.Help(fmt.Sprintf(format, v...))
}

func (o *mediaObject) Href(v string) MediaObject {
	o.IHref = v
	return o
}

func (o *mediaObject) Hreff(format string, v ...any) MediaObject {
	return o.Href(fmt.Sprintf(format, v...))
}

func (o *mediaObject) Rel(v string) MediaObject {
	o.IRel = v
	return o
}

func (o *mediaObject) Relf(format string, v ...any) MediaObject {
	return o.Rel(fmt.Sprintf(format, v...))
}

func (o *mediaObject) Order(v int) MediaObject {
	o.IOrder = v
	return o
}

func (o *mediaObject) Text(v string) MediaObject {
	content := app.Text(v)
	if o.IMedia == nil {
		return o.SVGIcon(SVGIconByTitle(v)).
			Content(content)
	}
	return o.Content(content)
}

func (o *mediaObject) Textf(format string, v ...any) MediaObject {
	return o.Text(fmt.Sprintf(format, v...))
}

func (o *mediaObject) Content(v ...app.UI) MediaObject {
	o.IContent = app.FilterUIElems(v...)
	return o
}

func (o *mediaObject) OnClick(v app.EventHandler) MediaObject {
	o.IOnClick = v
	return o
}

func (o *mediaObject) Render() app.UI {
	root := app.Elem(o.IRoot).
		ID(o.IID).
		Class("vi-media-object").
		Class(StyleFromBool("vi-media-object-only-media", len(o.IContent) == 0)).
		Class(StyleFromBool("vi-media-object-with-horizontal-centering", o.IWithHorizontalCentering)).
		Class(StyleFromBool("vi-media-object-with-min-size", o.IWithMinSize)).
		Class(StyleFromBool("vi-media-object-with-background", o.IWithBackground)).
		Class(StyleFromBool("vi-media-object-with-media-rounded", o.IWithMediaRounded)).
		Class(StyleFromBool("vi-media-object-with-media-border", o.IWithMediaBorder)).
		Class(StyleFromBool("vi-media-object-without-media-border", o.IWithoutMediaBorder)).
		Class(StyleFromBool("vi-media-object-with-controls", o.IWithControls)).
		Class(o.IMediaSize).
		Class(StyleFromBool("vi-media-object-with-ordering", o.IOrder > 0)).
		Class(StyleFromBool("vi-media-object-vertical", o.IVertical)).
		Class(o.IClass).
		Title(o.IHelp).
		Body(
			app.If(o.IOrder > 0, func() app.UI {
				return app.Div().Text(o.IOrder)
			}),
			o.IMedia,
			app.Range(o.IContent).Slice(func(i int) app.UI {
				return o.IContent[i]
			}),
		)

	switch o.IRoot {
	case "a":
		root = root.
			Attr("href", o.IHref).
			Attr("rel", o.IRel)
	}
	if o.IOnClick != nil {
		root = root.OnClick(o.IOnClick)
	}
	if o.IOuterElem != "" {
		root = app.Elem(o.IOuterElem).
			Body(root)
	}
	return root
}
