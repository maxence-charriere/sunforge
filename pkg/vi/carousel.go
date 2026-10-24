package vi

import (
	"fmt"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

const (
	DefaultCarouselButtonSVG = `<svg $id $class $help viewBox="0 0 24 24">
		<path fill="currentColor" d="M12,8A4,4 0 0,0 8,12A4,4 0 0,0 12,16A4,4 0 0,0 16,12A4,4 0 0,0 12,8Z" />
	</svg>`
)

type Carousel interface {
	app.Composer

	AsSection() Carousel
	ID(v string) Carousel
	Class(v string) Carousel
	QueryKey(v string) Carousel
	Content(v ...app.UI) Carousel
}

func NewCarousel() Carousel {
	return &carousel{IRoot: "div"}
}

type carousel struct {
	app.Compo

	IRoot     string
	IID       string
	IClass    string
	IQueryKey string
	IContent  []app.UI

	index int
}

func (c *carousel) AsSection() Carousel {
	c.IRoot = "section"
	return c
}

func (c *carousel) ID(v string) Carousel {
	c.IID = v
	return c
}

func (c *carousel) Class(v string) Carousel {
	c.IClass = app.AppendClass(c.IClass, v)
	return c
}

func (c *carousel) QueryKey(v string) Carousel {
	c.IQueryKey = v
	return c
}

func (c *carousel) Content(v ...app.UI) Carousel {
	c.IContent = app.FilterUIElems(v...)
	return c
}

func (c *carousel) OnMount(ctx app.Context) {
	c.initIndex(ctx)
}

func (c *carousel) OnNav(ctx app.Context) {
	c.initIndex(ctx)
}

func (c *carousel) initIndex(ctx app.Context) {
	if len(c.IContent) == 0 || c.IQueryKey == "" {
		return
	}
	index := ctx.Page().URL().Query().Get(c.IQueryKey)
	c.setIndex(index)
}

func (c *carousel) setIndex(v string) {
	index, _ := strconv.Atoi(v)
	if index < 0 {
		index = 0
	}
	if index >= len(c.IContent) {
		index = len(c.IContent) - 1
	}
	c.index = index
}

func (c *carousel) Render() app.UI {
	return app.Elem(c.IRoot).
		ID(c.IID).
		Class("vi-carousel").
		Class(c.IClass).
		Aria("label", "Carousel").
		Body(
			app.Range(c.IContent).Slice(func(i int) app.UI {
				var selected string
				if i == c.index {
					selected = "vi-carousel-item-show"
				}

				return app.Div().
					Class(selected).
					Aria("label", fmt.Sprintf("Slide %v", i+1)).
					Body(c.IContent[i])
			}),
			app.Div().
				Role("tablist").
				Aria("label", "Carousel Pagination").
				Body(
					app.Range(c.IContent).Slice(func(i int) app.UI {
						var selected string
						if i == c.index {
							selected = "vi-carousel-item-selected"
						}

						return app.Button().
							DataSet("index", i).
							Role("tab").
							Aria("label", fmt.Sprintf("Go to Slide %v", i+1)).
							Aria("controls", "carousel").
							Aria("selected", i == c.index).
							Class(selected).
							OnClick(c.showElement).Body(
							NewSVG().Data(DefaultCarouselButtonSVG),
						)
					}),
				),
		)
}

func (c *carousel) showElement(ctx app.Context, e app.Event) {
	index := ctx.JSSrc().
		Get("dataset").
		Get("index").
		String()
	c.setIndex(index)

	if c.IQueryKey == "" {
		return
	}

	u := *ctx.Page().URL()
	q := u.Query()
	if c.index == 0 {
		q.Del(c.IQueryKey)
	} else {
		q.Set(c.IQueryKey, strconv.Itoa(c.index))
	}
	u.RawQuery = q.Encode()
	ctx.Page().ReplaceURL(&u)
}
