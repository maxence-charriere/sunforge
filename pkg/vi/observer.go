package vi

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

type Observer interface {
	app.UI

	// Sets the ID.
	ID(v string) Observer

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Observer

	// Sets root data set.
	DataSet(k string, v any) Observer

	// Sets root element to aside.
	AsAside() Observer

	// Sets the content.
	Content(v ...app.UI) Observer

	// The function called when the component becomes visible.
	OnBecomeVisible(v func(app.Context)) Observer
}

// Create a element that receives notification when it become visible.
func NewObserver() Observer {
	return &observer{
		IRootType: "div",
	}
}

type observer struct {
	app.Compo

	IID              string
	IClass           string
	IRootType        string
	IContent         []app.UI
	IOnBecomeVisible func(app.Context)

	dataSets map[string]any
	root     app.Value
	release  func()
}

func (o *observer) ID(v string) Observer {
	o.IID = v
	return o
}

func (o *observer) Class(v ...string) Observer {
	o.IClass = app.AppendClass(o.IClass, v...)
	return o
}

func (o *observer) DataSet(k string, v any) Observer {
	if o.dataSets == nil {
		o.dataSets = map[string]any{}
	}
	o.dataSets[k] = v
	return o
}

func (o *observer) AsAside() Observer {
	o.IRootType = "aside"
	return o
}

func (o *observer) Content(v ...app.UI) Observer {
	o.IContent = app.FilterUIElems(v...)
	return o
}

func (o *observer) OnBecomeVisible(v func(app.Context)) Observer {
	o.IOnBecomeVisible = v
	return o
}

func (o *observer) OnMount(ctx app.Context) {
	ctx.Handle(notifyObserverIntersection, o.handleIntersect)

	o.root = ctx.JSSrc()
	ctx.NewActionWithValue(observeElement, o.root)
	o.release = func() {
		ctx.NewActionWithValue(unobserveElement, o.root)
	}
}

func (o *observer) OnDismount() {
	if o.release != nil {
		o.release()
	}
}

func (o *observer) Render() app.UI {
	return app.Elem(o.IRootType).
		DataSets(o.dataSets).
		ID(o.IID).
		Class("vi-observer").
		Class(o.IClass).
		Body(o.IContent...)
}

func (o *observer) handleIntersect(ctx app.Context, a app.Action) {
	source, ok := a.Value.(app.Value)
	if !ok || !source.Truthy() || !source.Equal(o.root) || o.IOnBecomeVisible == nil {
		return
	}

	o.IOnBecomeVisible(ctx)
}
