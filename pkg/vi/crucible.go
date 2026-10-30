package vi

import "github.com/maxence-charriere/go-app/v11/pkg/app"

type Crucible interface {
	app.Composer

	ID(v string) Crucible
	Class(v ...string) Crucible
	Icon(string) Crucible
	Name(string) Crucible
	Version(string) Crucible
	Menu(...app.UI) Crucible
	Index(...app.UI) Crucible
	Nav(...app.UI) Crucible
	Content(...app.UI) Crucible
	Footer(...app.UI) Crucible
	Ads(...app.UI) Crucible
	Notifications(...app.UI) Crucible
	NotificationsSVGIcon(string) Crucible
	NotificationsMark(bool) Crucible
}

func NewCrucible() Crucible {
	return &crucible{}
}

type crucible struct {
	app.Compo

	IID                  string
	IClass               string
	IIcon                string
	IName                string
	IVersion             string
	IMenu                []app.UI
	IIndex               []app.UI
	INav                 []app.UI
	IContent             []app.UI
	IFooter              []app.UI
	IAds                 []app.UI
	INotifications       []app.UI
	INotificationSVGIcon string
	INotificationsMark   bool

	pop []app.UI
}

func (c *crucible) ID(v string) Crucible {
	c.IID = v
	return c
}

func (c *crucible) Class(v ...string) Crucible {
	c.IClass = app.AppendClass(c.IClass, v...)
	return c
}

func (c *crucible) Icon(v string) Crucible {
	c.IIcon = v
	return c
}

func (c *crucible) Name(v string) Crucible {
	c.IName = v
	return c
}

func (c *crucible) Version(v string) Crucible {
	c.IVersion = v
	return c
}

func (c *crucible) Menu(v ...app.UI) Crucible {
	c.IMenu = app.FilterUIElems(v...)
	return c
}

func (c *crucible) Index(v ...app.UI) Crucible {
	c.IIndex = app.FilterUIElems(v...)
	return c
}

func (c *crucible) Nav(v ...app.UI) Crucible {
	c.INav = app.FilterUIElems(v...)
	return c
}

func (c *crucible) Content(v ...app.UI) Crucible {
	c.IContent = app.FilterUIElems(v...)
	return c
}

func (c *crucible) Footer(v ...app.UI) Crucible {
	c.IFooter = app.FilterUIElems(v...)
	return c
}

func (c *crucible) Ads(v ...app.UI) Crucible {
	c.IAds = app.FilterUIElems(v...)
	return c
}

func (c *crucible) Notifications(v ...app.UI) Crucible {
	c.INotifications = app.FilterUIElems(v...)
	return c
}

func (c *crucible) NotificationsSVGIcon(v string) Crucible {
	c.INotificationSVGIcon = v
	return c
}

func (c *crucible) NotificationsMark(v bool) Crucible {
	c.INotificationsMark = v
	return c
}

func (c *crucible) Render() app.UI {
	return app.Div().
		ID(c.IID).
		Class("vi-crucible").
		Class(StyleFromBool("vi-crucible-without-ads", len(c.IAds) == 0)).
		Class(StyleFromBool("vi-crucible-with-ads", len(c.IAds) != 0)).
		Class(c.IClass).
		Body(
			app.Nav().
				Class("vi-crucible-pane").
				Class("vi-unselectable").
				Body(
					app.A().
						Class("vi-crucible-logo").
						Href("/").
						Body(
							app.Img().
								Altf("%s Logo", c.IName).
								Src(c.IIcon),
							app.Div().Text(c.IName),
							app.Div().Text(c.IVersion),
						),
				),
			app.Div().
				Class("vi-crucible-content").
				Body(c.IContent...),
			app.Div().Class("vi-crucible-nav"),
			app.Div().
				Class("vi-crucible-ads").
				Body(c.IAds...),
			app.Div().Class("vi-crucible-pop"),
		)
}
