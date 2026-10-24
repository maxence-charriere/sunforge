package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Switcher interface {
	app.UI

	ID(v string) Switcher
	Class(v ...string) Switcher
	Text(v string) Switcher
	Textf(format string, v ...any) Switcher
	Help(v string) Switcher
	Helpf(format string, v ...any) Switcher
	LeftIcon(v app.UI) Switcher
	LeftSVGIcon(v string) Switcher
	LeftLabel(v string) Switcher
	LeftLabelf(format string, v ...any) Switcher
	RightIcon(v app.UI) Switcher
	RightSVGIcon(v string) Switcher
	RightLabel(v string) Switcher
	RightLabelf(format string, v ...any) Switcher
	Index(v int) Switcher
	OnChange(v func(ctx app.Context, idx int)) Switcher
}

func NewSwitcher() Switcher {
	return &switcher{}
}

type switcher struct {
	app.Compo

	IID         string
	IClass      string
	IText       string
	IHelp       string
	ILeftIcon   app.UI
	ILeftLabel  string
	IRightIcon  app.UI
	IRightLabel string
	IIndex      int
	IOnChange   func(app.Context, int)
}

func (s *switcher) ID(v string) Switcher {
	s.IID = v
	return s
}

func (s *switcher) Class(v ...string) Switcher {
	s.IClass = app.AppendClass(s.IClass, v...)
	return s
}

func (s *switcher) Text(v string) Switcher {
	s.IText = v
	return s
}

func (s *switcher) Textf(format string, v ...any) Switcher {
	return s.Text(fmt.Sprintf(format, v...))
}

func (s *switcher) Help(v string) Switcher {
	s.IHelp = v
	return s
}

func (s *switcher) Helpf(format string, v ...any) Switcher {
	return s.Help(fmt.Sprintf(format, v...))
}

func (s *switcher) LeftIcon(v app.UI) Switcher {
	s.ILeftIcon = v
	return s
}

func (s *switcher) LeftSVGIcon(v string) Switcher {
	return s.LeftIcon(NewSVG().Data(v))
}

func (s *switcher) LeftLabel(v string) Switcher {
	s.ILeftLabel = v
	return s
}

func (s *switcher) LeftLabelf(format string, v ...any) Switcher {
	return s.LeftLabel(fmt.Sprintf(format, v...))
}

func (s *switcher) RightIcon(v app.UI) Switcher {
	s.IRightIcon = v
	return s
}

func (s *switcher) RightSVGIcon(v string) Switcher {
	return s.RightIcon(NewSVG().Data(v))
}

func (s *switcher) RightLabel(v string) Switcher {
	s.IRightLabel = v
	return s
}

func (s *switcher) RightLabelf(format string, v ...any) Switcher {
	return s.RightLabel(fmt.Sprintf(format, v...))
}

func (s *switcher) Index(v int) Switcher {
	s.IIndex = v
	return s
}

func (s *switcher) OnChange(v func(app.Context, int)) Switcher {
	if v != nil {
		s.IOnChange = v
	}
	return s
}

func (s *switcher) Render() app.UI {
	selectedLabel := s.ILeftLabel
	if s.IIndex != 0 {
		selectedLabel = s.IRightLabel
	}

	return NewBox().
		WithGap().
		Class("vi-switcher").
		Class(s.IClass).
		OnClick(s.whenClicked).
		Content(
			app.Div().
				Class("h3").
				Body(
					app.Textf("%s - ", s.IText),
					app.Span().
						Class("accent").
						Text(selectedLabel),
				),
			NewBox().
				WithColumns().
				WithoutBackground().
				WithoutPadding().
				WithGap().
				Content(
					NewMediaObject().
						AsButton().
						WithMinSize().
						WithLargeMedia().
						Media(s.ILeftIcon).
						Class(StyleFromBool("accent", s.IIndex == 0)).
						OnClick(s.onLeftClick),
					NewMediaObject().
						AsButton().
						WithMinSize().
						WithLargeMedia().
						Media(s.IRightIcon).
						Class(StyleFromBool("accent", s.IIndex != 0)).
						OnClick(s.onRightClick),
				),
		)
}

func (s *switcher) whenClicked(ctx app.Context, e app.Event) {
	e.PreventDefault()
	e.StopImmediatePropagation()
}

func (s *switcher) onLeftClick(ctx app.Context, e app.Event) {
	if s.IOnChange != nil {
		s.IOnChange(ctx, 0)
	}
}

func (s *switcher) onRightClick(ctx app.Context, e app.Event) {
	if s.IOnChange != nil {
		s.IOnChange(ctx, 1)
	}
}
