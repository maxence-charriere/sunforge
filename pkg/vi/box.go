package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Box interface {
	app.Composer

	AsA() Box
	AsArticle() Box
	AsButton() Box
	AsSection() Box
	AsOl() Box
	AsUl() Box
	AsFigure() Box
	AsPre() Box
	LikeButton() Box
	WithHeader() Box
	WithGap() Box
	WithColumns() Box
	WithHorizontalCentering() Box
	WithBackgroundImage(v string) Box
	WithBackgroundImagef(format string, v ...any) Box
	WithMediaBackgroundImage(v string) Box
	WithMediaBackgroundImagef(format string, v ...any) Box
	WithMarkdown() Box
	WithoutBackground() Box
	WithoutPadding() Box
	WithoutTopPadding() Box
	WithoutVerticalCentering() Box
	WithoutSelection() Box
	ID(v string) Box
	Class(v string) Box
	Style(k, v string) Box
	Stylef(k, format string, v ...any) Box
	Help(v string) Box
	Helpf(format string, v ...any) Box
	Href(v string) Box
	Hreff(format string, v ...any) Box
	Rel(v string) Box
	Relf(format string, v ...any) Box
	Alt(v string) Box
	Altf(format string, v ...any) Box
	Loading(v string) Box
	Start(v int) Box
	Content(v ...app.UI) Box
	Text(v string) Box
	Textf(format string, v ...any) Box
	OnClick(v app.EventHandler) Box
}

func NewBox() Box {
	return &box{IRootType: "div"}
}

type box struct {
	app.Compo

	IRootType                 string
	IWithHeader               bool
	IWithGap                  bool
	IWithColumns              bool
	IWithHorizontalCentering  bool
	IWithMediaBackgroundImage bool
	IWithMarkdown             bool
	IWithoutBackground        bool
	IWithoutPadding           bool
	IWithoutTopPadding        bool
	IWithoutCentering         bool
	IWithoutSelection         bool
	ILikeButton               bool
	IID                       string
	IClass                    string
	IStyles                   map[string]string
	IHelp                     string
	IHref                     string
	IRel                      string
	IBackgroundImageURL       string
	IAlt                      string
	ILoading                  string
	IStart                    int
	IContent                  []app.UI
	IOnClick                  app.EventHandler
}

func (b *box) AsA() Box {
	b.IRootType = "a"
	return b
}

func (b *box) AsArticle() Box {
	b.IRootType = "article"
	return b
}

func (b *box) AsButton() Box {
	b.IRootType = "button"
	return b
}

func (b *box) AsSection() Box {
	b.IRootType = "section"
	return b
}

func (b *box) AsOl() Box {
	b.IRootType = "ol"
	return b
}

func (b *box) AsUl() Box {
	b.IRootType = "ul"
	return b
}

func (b *box) AsFigure() Box {
	b.IRootType = "figure"
	return b
}

func (b *box) AsPre() Box {
	b.IRootType = "pre"
	return b
}

func (b *box) LikeButton() Box {
	b.ILikeButton = true
	return b
}

func (b *box) WithHeader() Box {
	b.IWithHeader = true
	return b
}

func (b *box) WithGap() Box {
	b.IWithGap = true
	return b
}

func (b *box) WithColumns() Box {
	b.IWithColumns = true
	return b
}

func (b *box) WithHorizontalCentering() Box {
	b.IWithHorizontalCentering = true
	return b
}

func (b *box) WithoutBackground() Box {
	b.IWithoutBackground = true
	return b
}

func (b *box) WithoutPadding() Box {
	b.IWithoutPadding = true
	return b
}

func (b *box) WithoutTopPadding() Box {
	b.IWithoutTopPadding = true
	return b
}

func (b *box) WithoutVerticalCentering() Box {
	b.IWithoutCentering = true
	return b
}

func (b *box) WithBackgroundImage(v string) Box {
	b.IBackgroundImageURL = v
	return b
}

func (b *box) WithBackgroundImagef(format string, v ...any) Box {
	return b.WithBackgroundImage(fmt.Sprintf(format, v...))
}

func (b *box) WithMediaBackgroundImage(v string) Box {
	b.IWithMediaBackgroundImage = true
	return b.WithBackgroundImage(v)
}

func (b *box) WithMediaBackgroundImagef(format string, v ...any) Box {
	return b.WithMediaBackgroundImage(fmt.Sprintf(format, v...))
}

func (b *box) WithMarkdown() Box {
	b.IWithMarkdown = true
	return b
}

func (b *box) WithoutSelection() Box {
	b.IWithoutSelection = true
	return b
}

func (b *box) ID(v string) Box {
	b.IID = v
	return b
}

func (b *box) Class(v string) Box {
	b.IClass = app.AppendClass(b.IClass, v)
	return b
}

func (b *box) Style(k, v string) Box {
	if b.IStyles == nil {
		b.IStyles = make(map[string]string)
	}
	b.IStyles[k] = v
	return b
}

func (b *box) Stylef(k, format string, v ...any) Box {
	return b.Style(k, fmt.Sprintf(format, v...))
}

func (b *box) Help(v string) Box {
	b.IHelp = v
	return b
}

func (b *box) Helpf(format string, v ...any) Box {
	return b.Help(fmt.Sprintf(format, v...))
}

func (b *box) Href(v string) Box {
	b.IHref = v
	return b
}

func (b *box) Hreff(format string, v ...any) Box {
	return b.Href(fmt.Sprintf(format, v...))
}

func (b *box) Rel(v string) Box {
	b.IRel = v
	return b
}

func (b *box) Relf(format string, v ...any) Box {
	return b.Rel(fmt.Sprintf(format, v...))
}

func (b *box) Alt(v string) Box {
	b.IAlt = v
	return b
}

func (b *box) Altf(format string, v ...any) Box {
	return b.Alt(fmt.Sprintf(format, v...))
}

func (b *box) Loading(v string) Box {
	b.ILoading = v
	return b
}

func (b *box) Start(v int) Box {
	b.IStart = v
	return b
}

func (b *box) Content(v ...app.UI) Box {
	b.IContent = app.FilterUIElems(v...)
	return b
}

func (b *box) Text(v string) Box {
	return b.Content(app.Text(v))
}

func (b *box) Textf(format string, v ...any) Box {
	return b.Text(fmt.Sprintf(format, v...))
}

func (b *box) OnClick(v app.EventHandler) Box {
	b.IOnClick = v
	return b
}

func (b *box) Render() app.UI {
	root := app.Elem(b.IRootType).
		ID(b.IID).
		Class("vi-box").
		Class(b.IClass).
		Class(IfTrue(b.IWithHeader, "vi-box-with-header")).
		Class(IfTrue(b.IWithGap, "vi-box-with-gap")).
		Class(IfTrue(b.IWithColumns, "vi-box-with-columns")).
		Class(IfTrue(b.IWithHorizontalCentering, "vi-box-with-horizontal-centering")).
		Class(IfTrue(b.IWithMarkdown, "vi-box-with-markdown")).
		Class(IfTrue(b.IBackgroundImageURL != "", "vi-box-with-background-img")).
		Class(IfTrue(b.IWithMediaBackgroundImage, "vi-box-with-media-background-img")).
		Class(IfTrue(b.IWithoutBackground, "vi-box-without-background")).
		Class(IfTrue(b.IWithoutPadding, "vi-box-without-padding")).
		Class(IfTrue(b.IWithoutTopPadding, "vi-box-without-padding-top")).
		Class(IfTrue(b.IWithoutCentering, "vi-box-without-centering")).
		Class(IfTrue(b.IWithoutSelection, "vi-box-without-selection")).
		Class(IfTrue(b.ILikeButton, "vi-box-like-button")).
		Styles(b.IStyles).
		Title(b.IHelp).
		Body(
			app.Range(b.IContent).Slice(func(i int) app.UI {
				return b.IContent[i]
			}),
			app.If(b.IBackgroundImageURL != "", func() app.UI {
				return app.Img().
					Alt(b.IAlt).
					Src(b.IBackgroundImageURL).
					Loading(b.ILoading)
			}),
		)

	if b.IRootType == "a" {
		root = root.
			Attr("href", b.IHref).
			Attr("rel", b.IRel)
	}
	if b.IRootType == "ol" {
		root = root.Attr("start", b.IStart)
	}
	if b.IOnClick != nil {
		root = root.OnClick(b.IOnClick)
	}
	return root
}
