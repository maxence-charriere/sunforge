package vi

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Thumbnail interface {
	ID(string) Thumbnail
	IDf(string, ...any) Thumbnail
	Class(...string) Thumbnail
	Image(string) Thumbnail
	Imagef(string, ...any) Thumbnail
	FallbackImage(string) Thumbnail
	FallbackImagef(string, ...any) Thumbnail
	Title(string) Thumbnail
	Titlef(string, ...any) Thumbnail
	Href(string) Thumbnail
	Hreff(string, ...any) Thumbnail
	Time(time.Time) Thumbnail
	Highlight(bool) Thumbnail
	OnClick(app.EventHandler) Thumbnail

	app.Composer
}

func NewThumbnail() Thumbnail {
	return &thumbnail{IRoot: "div"}
}

type thumbnail struct {
	app.Compo

	IRoot          string
	IID            string
	IClass         string
	IImage         string
	IFallbackImage string
	ITitle         string
	IHref          string
	ITime          time.Time
	IHighlight     bool
	IOnClick       app.EventHandler
}

func (t *thumbnail) ID(v string) Thumbnail {
	t.IID = v
	return t
}

func (t *thumbnail) IDf(format string, v ...any) Thumbnail {
	return t.ID(fmt.Sprintf(format, v...))
}

func (t *thumbnail) Class(v ...string) Thumbnail {
	t.IClass = app.AppendClass(t.IClass, v...)
	return t
}

func (t *thumbnail) Image(v string) Thumbnail {
	t.IImage = v
	return t
}

func (t *thumbnail) Imagef(format string, v ...any) Thumbnail {
	return t.Image(fmt.Sprintf(format, v...))
}

func (t *thumbnail) FallbackImage(v string) Thumbnail {
	t.IFallbackImage = v
	return t
}

func (t *thumbnail) FallbackImagef(format string, v ...any) Thumbnail {
	return t.FallbackImage(fmt.Sprintf(format, v...))
}

func (t *thumbnail) Title(v string) Thumbnail {
	t.ITitle = v
	return t
}

func (t *thumbnail) Titlef(format string, v ...any) Thumbnail {
	return t.Title(fmt.Sprintf(format, v...))
}

func (t *thumbnail) Href(v string) Thumbnail {
	t.IHref = v
	t.IRoot = "a"
	return t
}

func (t *thumbnail) Hreff(format string, v ...any) Thumbnail {
	return t.Href(fmt.Sprintf(format, v...))
}

func (t *thumbnail) Time(v time.Time) Thumbnail {
	t.ITime = v
	return t
}

func (t *thumbnail) Highlight(v bool) Thumbnail {
	t.IHighlight = v
	return t
}

func (t *thumbnail) OnClick(e app.EventHandler) Thumbnail {
	t.IOnClick = e
	return t
}

func (t *thumbnail) Render() app.UI {
	image := t.IImage
	if image == "" {
		image = t.IFallbackImage
	}

	var highlight string
	if t.IHighlight {
		highlight = "vi-thumbnail-with-highlight"
	}

	root := NewBox()
	if t.IHref != "" {
		root = root.AsA()
	} else {
		root = root.LikeButton()
	}

	return root.
		WithoutPadding().
		ID(t.IID).
		Class("vi-thumbnail").
		Class(t.IClass).
		Href(t.IHref).
		OnClick(t.IOnClick).
		Content(
			app.Img().
				Src(image).
				Alt(t.ITitle),
			app.Div(),
			NewBox().
				WithoutBackground().
				Content(
					app.H3().
						Class("h2").
						Text(t.ITitle),
					NewMediaObject().
						SVGIcon(UpdateSVG).
						Content(
							app.Time().
								Class(highlight).
								DateTime(t.ITime.Format(time.RFC3339)).
								Textf("%s ago", TimeSince(t.ITime)),
						),
				),
		)
}
