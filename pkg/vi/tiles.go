package vi

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Tiles interface {
	app.UI

	AsSection() Tiles
	AsAside() Tiles
	Smaller() Tiles
	WithLargeItems() Tiles
	ID(v string) Tiles
	Class(v ...string) Tiles
	Aria(k, v string) Tiles
	Ariaf(k, format string, v ...any) Tiles
	Content(v ...app.UI) Tiles
}

func NewTiles() Tiles {
	return &tiles{IRootType: "div"}
}

type tiles struct {
	app.Compo

	IID       string
	IClass    string
	IRootType string
	IItemSize string
	IAria     map[string]string
	IContent  []app.UI
}

func (t *tiles) AsSection() Tiles {
	t.IRootType = "section"
	return t
}

func (t *tiles) AsAside() Tiles {
	t.IRootType = "aside"
	return t
}

func (t *tiles) Smaller() Tiles {
	t.IItemSize = "vi-tiles-smaller"
	return t
}

func (t *tiles) WithLargeItems() Tiles {
	t.IItemSize = "vi-tiles-with-large-items"
	return t
}

func (t *tiles) ID(v string) Tiles {
	t.IID = v
	return t
}

func (t *tiles) Class(v ...string) Tiles {
	t.IClass = app.AppendClass(t.IClass, v...)
	return t
}

func (t *tiles) Aria(k, v string) Tiles {
	if t.IAria == nil {
		t.IAria = map[string]string{}
	}
	t.IAria[k] = v
	return t
}

func (t *tiles) Ariaf(k, format string, v ...any) Tiles {
	return t.Aria(k, fmt.Sprintf(format, v...))
}

func (t *tiles) Content(v ...app.UI) Tiles {
	t.IContent = app.FilterUIElems(v...)
	return t
}

func (t *tiles) Render() app.UI {
	tiles := app.Elem(t.IRootType).
		ID(t.IID).
		Class("vi-tiles").
		Class(t.IItemSize).
		Class(t.IClass).
		Body(
			app.Range(t.IContent).Slice(func(i int) app.UI {
				return t.IContent[i]
			}),
		)

	for k, v := range t.IAria {
		tiles = tiles.Aria(k, v)
	}
	return tiles
}
