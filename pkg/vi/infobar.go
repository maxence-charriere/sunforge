package vi

import (
	"fmt"
	"sync"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type InfoBar interface {
	app.Composer

	WithBorder() InfoBar
	ID(v string) InfoBar
	Class(v ...string) InfoBar
	Help(v string) InfoBar
	Helpf(format string, v ...any) InfoBar
	Value(v float64) InfoBar
	Max(v float64) InfoBar
	Tags(v ...app.UI) InfoBar
}

func NewInfoBar() InfoBar {
	return &infoBar{}
}

type infoBar struct {
	app.Compo

	IID         string
	IClass      string
	IHelp       string
	IValue      float64
	IMax        float64
	IWithBorder bool
	ITags       []app.UI

	initOnce sync.Once
	fill     bool
}

func (b *infoBar) WithBorder() InfoBar {
	b.IWithBorder = true
	return b
}

func (b *infoBar) ID(v string) InfoBar {
	b.IID = v
	return b
}

func (b *infoBar) Class(v ...string) InfoBar {
	b.IClass = app.AppendClass(b.IClass, v...)
	return b
}

func (b *infoBar) Help(v string) InfoBar {
	b.IHelp = v
	return b
}

func (b *infoBar) Helpf(format string, v ...any) InfoBar {
	return b.Help(fmt.Sprintf(format, v...))
}

func (b *infoBar) Value(v float64) InfoBar {
	b.IValue = v
	return b
}

func (b *infoBar) Max(v float64) InfoBar {
	b.IMax = v
	return b
}

func (b *infoBar) Tags(v ...app.UI) InfoBar {
	b.ITags = app.FilterUIElems(v...)
	return b
}

func (b *infoBar) Render() app.UI {
	return NewObserver().
		OnBecomeVisible(b.whenVisible).
		Content(
			app.Div().
				ID(b.IID).
				Class("vi-infobar").
				Class(StyleFromBool("vi-infobar-with-border", b.IWithBorder)).
				Class(b.IClass).
				Title(b.IHelp).
				Body(
					app.Div().
						Role("img").
						Body(
							app.If(!b.fill, func() app.UI {
								return app.Div()
							}).Else(func() app.UI {
								return app.Div().
									Style("transform", fmt.Sprintf("translateX(-%.3f%%)", 100-(b.IValue*100/b.IMax)))
							}),
						),
					app.If(len(b.ITags) != 0, func() app.UI {
						return NewBox().
							AsUl().
							WithColumns().
							WithoutPadding().
							WithoutBackground().
							Content(b.ITags...)
					}),
				),
		)
}

func (b *infoBar) whenVisible(ctx app.Context) {
	b.initOnce.Do(func() {
		b.fill = true
	})
}
