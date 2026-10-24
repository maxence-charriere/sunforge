package vi

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Poppable interface {
	app.UI

	// Sets the ID.
	ID(v string) Poppable

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Poppable

	// Sets a style. Multiple styles can be defined by successive calls.
	Style(k, v string) Poppable

	// Sets root element to li.
	AsLi() Poppable

	// Sets root element to article.
	AsArticle() Poppable

	// Sets the content.
	Content(v app.UI) Poppable

	// Sets the content that pops when hovering, clicking, or taping.
	Pop(v ...app.UI) Poppable

	// Prevents pop when mouse hover.
	WithoutHover() Poppable
}

// Creates an element that can pop content when clicked.
func NewPoppable() Poppable {
	return &poppable{
		IRootType: "div",
	}
}

type poppable struct {
	app.Compo

	IID           string
	IClass        string
	IRootType     string
	IContent      app.UI
	IStyles       map[string]string
	IPop          []app.UI
	IWithoutHover bool
}

func (p *poppable) ID(v string) Poppable {
	p.IID = v
	return p
}

func (p *poppable) Class(v ...string) Poppable {
	p.IClass = app.AppendClass(p.IClass, v...)
	return p
}

func (p *poppable) Style(k, v string) Poppable {
	if p.IStyles == nil {
		p.IStyles = make(map[string]string)
	}
	p.IStyles[k] = v
	return p
}

func (p *poppable) AsLi() Poppable {
	p.IRootType = "li"
	return p
}

func (p *poppable) AsArticle() Poppable {
	p.IRootType = "article"
	return p
}

func (p *poppable) Content(v app.UI) Poppable {
	if content := app.FilterUIElems(v); len(content) != 0 {
		p.IContent = content[0]
	}
	return p
}

func (p *poppable) Pop(v ...app.UI) Poppable {
	p.IPop = app.FilterUIElems(v...)
	return p
}

func (p *poppable) WithoutHover() Poppable {
	p.IWithoutHover = true
	return p
}

func (p *poppable) Render() app.UI {
	return app.Elem(p.IRootType).
		ID(p.IID).
		Class("vi-poppable").
		Class(p.IClass).
		Styles(p.IStyles).
		OnMouseMove(p.onMouseHover).
		OnMouseLeave(p.onMouseLeave).
		OnClick(p.onClick).
		Body(p.IContent)
}

func (p *poppable) onMouseHover(ctx app.Context, e app.Event) {
	if p.IWithoutHover {
		return
	}
	ctx.NewActionWithValue(showPoppable, p.IPop)
}

func (p *poppable) onMouseLeave(ctx app.Context, e app.Event) {
	if p.IWithoutHover {
		return
	}
	ctx.NewAction(hidePoppable)
}

func (p *poppable) onClick(ctx app.Context, e app.Event) {
	e.PreventDefault()
	e.StopImmediatePropagation()
	ctx.NewActionWithValue(showPoppable, p.IPop, app.T("fullscreen", true))
}
