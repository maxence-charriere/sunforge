package vi

import (
	"fmt"
	"sync"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

const (
	crucibleContentOnlyWidth = 1302
)

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
	NotificationsBadge(bool) Crucible
}

func NewCrucible() Crucible {
	return &crucible{}
}

type crucible struct {
	app.Compo

	IID                 string
	IClass              string
	IIcon               string
	IName               string
	IVersion            string
	IMenu               []app.UI
	IIndex              []app.UI
	INav                []app.UI
	IContent            []app.UI
	IFooter             []app.UI
	IAds                []app.UI
	INotifications      []app.UI
	INotificationsBadge bool

	panePopped           bool
	notificationsVisible bool
	installVisible       bool
	appleBrowser         bool
	updateVisible        bool

	pop             []app.UI
	popX            int
	popY            int
	fullscreenPop   bool
	initPopPosition sync.Once
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

func (c *crucible) NotificationsBadge(v bool) Crucible {
	c.INotificationsBadge = v
	return c
}

func (c *crucible) OnMount(ctx app.Context) {
	ctx.Handle(showPop, c.handleShowPop)
	ctx.Handle(hidePop, c.handleHidePop)
	ctx.Handle(closePop, c.handleClosePop)

	c.OnAppInstallChange(ctx)
	c.OnAppUpdate(ctx)
}

func (c *crucible) OnNav(ctx app.Context) {
}

func (c *crucible) OnAppInstallChange(ctx app.Context) {
	c.installVisible = ctx.IsAppInstallable()
	c.appleBrowser = ctx.IsAppleBrowser()
}

func (c *crucible) OnAppUpdate(ctx app.Context) {
	c.updateVisible = ctx.AppUpdateAvailable()
}

func (c *crucible) Render() app.UI {
	return app.Div().
		ID(c.IID).
		Class("vi-crucible").
		Class(StyleFromBool("vi-crucible-without-ads", len(c.IAds) == 0)).
		Class(StyleFromBool("vi-crucible-with-ads", len(c.IAds) != 0)).
		Class(StyleFromBool("vi-crucible-pane-popped", c.panePopped)).
		Class(StyleFromBool("vi-crucible-pop-visible", len(c.pop) != 0)).
		Class(StyleFromBool("vi-crucible-pop-fullscreen", c.fullscreenPop)).
		Class(StyleFromBool("vi-crucible-notifications-visible", c.notificationsVisible)).
		Class(c.IClass).
		Body(
			app.Nav().
				Class("vi-crucible-pane").
				Class("vi-unselectable").
				OnClick(c.closeOverlay).
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
					app.Div().
						Class("vi-crucible-menu").
						Body(c.IMenu...),
					app.Div().
						Class("vi-crucible-index").
						Body(c.IIndex...),
				),

			app.Div().
				Class("vi-crucible-content").
				Body(
					app.Main().
						ID(PageTopID).
						Body(c.IContent...),
					app.Footer().Body(c.IFooter...),
				),

			app.Nav().
				Class("vi-crucible-nav").
				Body(
					app.Ul().Body(
						NewMediaObject().
							InLi().
							AsButton().
							SVGIcon(MenuSVG).
							OnClick(c.popPane),
						app.Range(c.INav).
							Slice(func(i int) app.UI {
								return c.INav[i]
							}),
						NewMediaObject().
							InLi().
							AsButton().
							SVGIcon(ShareSVG).
							OnClick(ShareCurrentURL),
						app.If(c.installVisible, func() app.UI {
							if !c.appleBrowser {
								return NewMediaObject().
									InLi().
									AsButton().
									SVGIcon(InstallSVG).
									OnClick(InstallPWA)
							}
							return NewPoppable().
								WithoutHover().
								Content(
									NewMediaObject().
										InLi().
										AsButton().
										SVGIcon(InstallSVG),
								).
								Pop(NewHowToInstallOnAppleDevices(c.IName))
						}),
						NewMediaObject().
							InLi().
							AsButton().
							Class(StyleFromBool("vi-crucible-notifications-badge", c.INotificationsBadge)).
							SVGIcon(BellSVG).
							OnClick(c.showNotifications),
					),
				),

			app.Div().
				Class("vi-crucible-ads").
				Body(
					app.Aside().Body(c.IAds...),
				),

			app.Div().
				Class("vi-crucible-overlay").
				OnClick(c.closeOverlay),

			app.If(app.IsClient, func() app.UI {
				return app.Aside().
					Class("vi-crucible-pop").
					Class("vi-unselectable").
					Style("transform", fmt.Sprintf("translate(%vpx, %vpx)", c.popX, c.popY)).
					OnClick(c.closeOverlay).
					Body(
						app.Range(c.pop).Slice(func(i int) app.UI {
							return c.pop[i]
						}),
					)
			}),

			app.Div().
				Class("vi-crucible-notifications").
				Class("vi-unselectable").
				OnClick(c.closeOverlay).
				Body(
					app.Aside().Body(c.INotifications...),
				),

			app.If(c.updateVisible, func() app.UI {
				return newUpdatePWA(c.IName).
					Class("vi-crucible-update").
					Class("vi-unselectable")
			}),
		)
}

func (c *crucible) popPane(ctx app.Context, e app.Event) {
	c.panePopped = true
}

func (c *crucible) closeOverlay(ctx app.Context, e app.Event) {
	c.panePopped = false
	c.notificationsVisible = false
	c.fullscreenPop = false
	c.pop = nil
}

func (c *crucible) showNotifications(ctx app.Context, e app.Event) {
	c.notificationsVisible = true
}

func (c *crucible) handleShowPop(ctx app.Context, a app.Action) {
	if c.fullscreenPop {
		return
	}

	pop, ok := a.Value.([]app.UI)
	if !ok {
		return
	}

	c.fullscreenPop = a.Tags.Get("fullscreen") == "true"
	if c.fullscreenPop {
		c.pop = pop
		return
	}

	if windowWidth, _ := app.Window().Size(); windowWidth < crucibleContentOnlyWidth {
		return
	}

	c.pop = pop
	c.initPopPosition.Do(func() {
		c.popX, c.popY = app.Window().CursorPosition()
	})

	ctx.Defer(func(ctx app.Context) {
		x, y := app.Window().CursorPosition()

		popElements := app.Window().
			Get("document").
			Call("getElementsByClassName", "vi-crucible-pop")
		if popElements.Length() == 0 {
			return
		}

		popElement := popElements.Index(0)
		popWidth := popElement.Get("offsetWidth").Int()
		popHeight := popElement.Get("offsetHeight").Int()
		windowWidth, windowHeight := app.Window().Size()

		if x+popWidth > windowWidth {
			x = windowWidth - popWidth
		}
		x = max(0, x)

		if y+popHeight > windowHeight {
			y = windowHeight - popHeight
		}
		y = max(0, y)

		if x != c.popX || y != c.popY {
			ctx.Dispatch(func(ctx app.Context) {
				c.popX = x
				c.popY = y
			})
		}
	})
}

func (c *crucible) handleHidePop(ctx app.Context, a app.Action) {
	if !c.fullscreenPop {
		c.pop = nil
	}
}

func (c *crucible) handleClosePop(ctx app.Context, a app.Action) {
	if c.fullscreenPop {
		c.closeOverlay(ctx, app.Event{})
	}
}
