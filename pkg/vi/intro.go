package vi

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type IntroOverview interface {
	app.Composer

	ID(v string) IntroOverview
	IDf(format string, v ...any) IntroOverview
	Class(v ...string) IntroOverview
	Content(v ...app.UI) IntroOverview
}

func NewIntroOverview() IntroOverview {
	return &introOverview{}
}

type introOverview struct {
	app.Compo

	IID      string
	IClass   string
	IContent []app.UI
}

func (o *introOverview) ID(v string) IntroOverview {
	o.IID = v
	return o
}

func (o *introOverview) IDf(format string, v ...any) IntroOverview {
	return o.ID(fmt.Sprintf(format, v...))
}

func (o *introOverview) Class(v ...string) IntroOverview {
	o.IClass = app.AppendClass(o.IClass, v...)
	return o
}

func (o *introOverview) Content(v ...app.UI) IntroOverview {
	o.IContent = app.FilterUIElems(v...)
	return o
}

func (o *introOverview) Render() app.UI {
	return app.Ul().
		ID(o.IID).
		Class("vi-intro-overview").
		Class(o.IClass).
		Aria("label", "Overview Icons").
		Body(
			app.Range(o.IContent).Slice(func(i int) app.UI {
				return app.Li().Body(o.IContent[i])
			}),
		)
}

type IntroIcon interface {
	app.Composer

	WithIconBorder() IntroIcon
	WithRoundedIcon() IntroIcon
	ID(v string) IntroIcon
	IDf(format string, v ...any) IntroIcon
	Class(v ...string) IntroIcon
	Icon(v app.UI) IntroIcon
	SVGIcon(v string) IntroIcon
	Label(v string) IntroIcon
	Labelf(format string, v ...any) IntroIcon
	Time(t time.Time, v string) IntroIcon
	Timef(t time.Time, format string, v ...any) IntroIcon
	Href(v string) IntroIcon
	Hreff(format string, v ...any) IntroIcon
	OnClick(h app.EventHandler) IntroIcon
}

func NewIntroIcon() IntroIcon {
	return &introIcon{IRoot: "div"}
}

type introIcon struct {
	app.Compo

	IRoot        string
	IID          string
	IClass       string
	IIcon        app.UI
	ILabel       app.UI
	IHref        string
	IIconBorder  bool
	IRoundedIcon bool
	IOnClick     app.EventHandler
}

func (i *introIcon) WithIconBorder() IntroIcon {
	i.IIconBorder = true
	return i
}

func (i *introIcon) WithRoundedIcon() IntroIcon {
	i.IRoundedIcon = true
	return i
}

func (i *introIcon) ID(v string) IntroIcon {
	i.IID = v
	return i
}

func (i *introIcon) IDf(format string, v ...any) IntroIcon {
	return i.ID(fmt.Sprintf(format, v...))
}

func (i *introIcon) Class(v ...string) IntroIcon {
	i.IClass = app.AppendClass(i.IClass, v...)
	return i
}

func (i *introIcon) Icon(v app.UI) IntroIcon {
	if icon := app.FilterUIElems(v); len(icon) > 0 {
		i.IIcon = icon[0]
	}
	return i
}

func (i *introIcon) SVGIcon(v string) IntroIcon {
	return i.Icon(NewSVG().Data(v))
}

func (i *introIcon) Label(v string) IntroIcon {
	i.ILabel = app.Text(v)
	if i.IIcon == nil {
		return i.SVGIcon(SVGIconByTitle(v))
	}
	return i
}

func (i *introIcon) Labelf(format string, v ...any) IntroIcon {
	return i.Label(fmt.Sprintf(format, v...))
}

func (i *introIcon) Time(t time.Time, v string) IntroIcon {
	i.ILabel = app.Time().
		DateTime(t.Format(time.RFC3339)).
		Text(v)
	if i.IIcon == nil {
		return i.SVGIcon(SVGIconByTitle(v))
	}
	return i
}

func (i *introIcon) Timef(t time.Time, format string, v ...any) IntroIcon {
	return i.Time(t, fmt.Sprintf(format, v...))
}

func (i *introIcon) Href(v string) IntroIcon {
	i.IRoot = "a"
	i.IHref = v
	return i
}

func (i *introIcon) Hreff(format string, v ...any) IntroIcon {
	return i.Href(fmt.Sprintf(format, v...))
}

func (i *introIcon) OnClick(v app.EventHandler) IntroIcon {
	i.IOnClick = v
	return i
}

func (i *introIcon) Render() app.UI {
	icon := app.Elem(i.IRoot).
		ID(i.IID).
		Class("vi-intro-icon").
		Class(StyleFromBool("vi-intro-icon-button", i.IHref != "" || i.IOnClick != nil)).
		Class(StyleFromBool("vi-intro-icon-with-icon-border", i.IIconBorder)).
		Class(StyleFromBool("vi-intro-icon-with-rounded-icon", i.IRoundedIcon)).
		Class(i.IClass).
		Body(
			i.IIcon,
			app.Span().Body(i.ILabel),
		)

	if i.IHref != "" {
		icon = icon.Attr("href", i.IHref)
	}
	if i.IOnClick != nil {
		icon = icon.OnClick(i.IOnClick)
	}
	return icon
}
