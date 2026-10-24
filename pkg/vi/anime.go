package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Anime interface {
	app.Composer

	ID(string) Anime
	IDf(string, ...any) Anime
	Class(...string) Anime
}

func NewAnime(srcs ...string) Anime {
	return &anime{ISrcs: srcs}
}

type anime struct {
	app.Compo

	IID    string
	IClass string
	ISrcs  []string
}

func (a *anime) ID(v string) Anime {
	a.IID = v
	return a
}

func (a *anime) IDf(format string, v ...any) Anime {
	return a.ID(fmt.Sprintf(format, v...))
}

func (a *anime) Class(v ...string) Anime {
	a.IClass = app.AppendClass(a.IClass, v...)
	return a
}

func (a *anime) Render() app.UI {
	return app.Div().
		ID(a.IID).
		Class("vi-anime").
		Class(a.IClass).
		Role("img").
		Body(
			app.Range(a.ISrcs).Slice(func(i int) app.UI {
				return app.Img().
					Src(a.ISrcs[i]).
					Alt("")
			}),
		)
}
