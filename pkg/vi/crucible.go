package vi

import (
	"fmt"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

const (
	crucibleContentOnlyWidth = 1302
)

// Crucible is a responsive application layout with branding, menus, navigation,
// main content, a footer, optional ads, and notifications. It supports popovers,
// intersection observers, and progressive web app installation and update prompts.
type Crucible interface {
	app.Composer

	ID(v string) Crucible
	Class(v ...string) Crucible
	Icon(string) Crucible
	Version(string) Crucible
	Menu(...app.UI) Crucible
	Index(...app.UI) Crucible
	Nav(...app.UI) Crucible
	Content(...app.UI) Crucible
	Footer(...app.UI) Crucible
	Ads(...app.UI) Crucible
	Notifications(...app.UI) Crucible
	NotificationsBadge(bool) Crucible
	WithFollowMeMenuSection() Crucible
}

// NewCrucible creates a responsive application layout for the named application.
// The name is used in the branding, installation and update prompts, and footer
// copyright notice, which uses the year in which the layout is created.
func NewCrucible(appName string) Crucible {
	return &crucible{
		IAppName:    appName,
		currentYear: time.Now().Year(),
	}
}

type crucible struct {
	app.Compo

	IID                      string
	IClass                   string
	IIcon                    string
	IAppName                 string
	IVersion                 string
	IMenu                    []app.UI
	IIndex                   []app.UI
	INav                     []app.UI
	IContent                 []app.UI
	IFooter                  []app.UI
	IAds                     []app.UI
	INotifications           []app.UI
	INotificationsBadge      bool
	IWithFollowMeMenuSection bool

	panePopped           bool
	notificationsVisible bool
	installVisible       bool
	appleBrowser         bool
	updateVisible        bool
	currentYear          int

	pop             []app.UI
	popX            int
	popY            int
	fullscreenPop   bool
	initPopPosition sync.Once

	observer      app.Value
	closeObserver func()
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

func (c *crucible) WithFollowMeMenuSection() Crucible {
	c.IWithFollowMeMenuSection = true
	return c
}

func (c *crucible) OnMount(ctx app.Context) {
	ctx.Handle(showPop, c.handleShowPop)
	ctx.Handle(hidePop, c.handleHidePop)
	ctx.Handle(closePop, c.handleClosePop)
	ctx.Handle(observeElement, c.handleObserveElement)
	ctx.Handle(unobserveElement, c.handleUnobserveElement)

	c.OnAppInstallChange(ctx)
	c.OnAppUpdate(ctx)

	c.initObservables(ctx)
}

func (c *crucible) OnNav(ctx app.Context) {
	c.closeOverlay(ctx, app.Event{})

	if ctx.Page().URL().Fragment == "" {
		ScrollToTop(ctx)
	}
}

func (c *crucible) OnAppInstallChange(ctx app.Context) {
	c.installVisible = ctx.IsAppInstallable()
	c.appleBrowser = ctx.IsAppleBrowser()
}

func (c *crucible) OnAppUpdate(ctx app.Context) {
	c.updateVisible = ctx.AppUpdateAvailable()
}

func (c *crucible) OnDismount() {
	c.closeObservables()
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
								Altf("%s Logo", c.IAppName).
								Src(c.IIcon),
							app.Div().Text(c.IAppName),
							app.Div().Text(c.IVersion),
						),
					app.Div().
						Class("vi-crucible-menu").
						Class(StyleFromBool("vi-crucible-menu-follow-me-section", c.IWithFollowMeMenuSection)).
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
					app.Footer().
						Class("vi-unselectable").
						Body(
							app.Ul().Body(
								NewMediaObject().
									InLi().
									AsA().
									Class(StyleFromURLPath("accent", "/terms")).
									Href("/terms").
									Text("Terms of Service"),
								NewMediaObject().
									InLi().
									AsA().
									Class(StyleFromURLPath("accent", "/privacy")).
									Href("/privacy").
									Text("Privacy Policy"),
								app.Range(c.IFooter).Slice(func(i int) app.UI {
									return app.Li().Body(c.IFooter[i])
								}),
								NewMediaObject().
									InLi().
									SVGIcon(CopyrightSVG).
									Textf("%v %s",
										c.currentYear,
										c.IAppName,
									),
							),
						),
				),

			app.Nav().
				Class("vi-crucible-nav").
				Class("vi-unselectable").
				Body(
					app.Ul().Body(
						NewMediaObject().
							InLi().
							AsButton().
							WithMediaBorder().
							WithMediaRounded().
							Image(c.IIcon).
							OnClick(c.popPane),
						NewMediaObject().
							InLi().
							AsButton().
							Class(Disabled(!CanNavigateBack())).
							SVGIcon(ArrowLeftSVG).
							OnClick(NavigateBack),
						app.Range(c.INav).
							Slice(func(i int) app.UI {
								return app.Li().Body(c.INav[i])
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
								Pop(NewHowToInstallOnAppleDevices(c.IAppName))
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
				return newUpdatePWA(c.IAppName).
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

func (c *crucible) initObservables(ctx app.Context) {
	intersectionObserver := app.Window().Get("IntersectionObserver")
	if !intersectionObserver.Truthy() {
		return
	}

	intersectionObserverForeach := app.FuncOf(func(this app.Value, args []app.Value) any {
		entry := args[0]
		if !entry.Get("isIntersecting").Bool() {
			return nil
		}
		ctx.NewActionWithValue(notifyObserverIntersection, entry.Get("target"))
		return nil
	})

	intersectionObserverCallback := app.FuncOf(func(this app.Value, args []app.Value) any {
		entries := args[0]
		entries.Call("forEach", intersectionObserverForeach)
		return nil
	})

	c.observer = intersectionObserver.New(intersectionObserverCallback)
	c.closeObserver = func() {
		c.observer.Call("disconnect")
		intersectionObserverCallback.Release()
		intersectionObserverForeach.Release()
	}
}

func (c *crucible) closeObservables() {
	if c.closeObserver == nil {
		return
	}
	c.closeObserver()
	c.closeObserver = nil
	c.observer = nil
}

func (c *crucible) handleObserveElement(ctx app.Context, a app.Action) {
	elem, ok := a.Value.(app.Value)
	if !ok || !elem.Truthy() || c.observer == nil {
		return
	}
	c.observer.Call("observe", elem)
}

func (c *crucible) handleUnobserveElement(ctx app.Context, a app.Action) {
	elem, ok := a.Value.(app.Value)
	if !ok || !elem.Truthy() || c.observer == nil {
		return
	}
	c.observer.Call("unobserve", elem)
}
