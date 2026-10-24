package vi

import (
	"fmt"
	"sync"

	"github.com/maxence-charriere/go-app/v11/pkg/analytics"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

var (
	// The default SVG used for app menu.
	DefaultAppMenuSVG = `<svg $id $class $help viewBox="0 0 24 24">
		<path fill="currentColor" d="M3,6H21V8H3V6M3,11H21V13H3V11M3,16H21V18H3V16Z" />
	</svg>`

	// The default SVG used for notifications.
	DefaultNotificationSVG = `<svg $id $class $help viewBox="0 0 24 24">
    	<path fill="currentColor" d="M21,19V20H3V19L5,17V11C5,7.9 7.03,5.17 10,4.29C10,4.19 10,4.1 10,4A2,2 0 0,1 12,2A2,2 0 0,1 14,4C14,4.1 14,4.19 14,4.29C16.97,5.17 19,7.9 19,11V17L21,19M14,21A2,2 0 0,1 12,23A2,2 0 0,1 10,21" />
	</svg>`

	// The default SVG used for app updates.
	DefaultDownloadSVG = `<svg $id $class $help viewBox="0 0 24 24">
		<path fill="currentColor" d="M5,20H19V18H5M19,9H15V3H9V9H5L12,16L19,9Z" />
	</svg>`

	// The default SVG used for poppable close button.
	DefaultCloseSVG = `<svg $id $class $help viewBox="0 0 24 24">
		<path fill="currentColor" d="M12,2C17.53,2 22,6.47 22,12C22,17.53 17.53,22 12,22C6.47,22 2,17.53 2,12C2,6.47 6.47,2 12,2M15.59,7L12,10.59L8.41,7L7,8.41L10.59,12L7,15.59L8.41,17L12,13.41L15.59,17L17,15.59L13.41,12L17,8.41L15.59,7Z" />
	</svg>`
)

const (
	// The id of the top of the page.
	PageTopID = "vi-page-top"

	showPoppable               = "/actions/vi/layout/poppable/show"
	hidePoppable               = "/actions/vi/layout/poppable/hide"
	closePoppable              = "/actions/vi/layout/poppable/close"
	notifyObserverIntersection = "/actions/vi/layout/observer/intersection"
	observeElement             = "/actions/vi/observer/observe"
	unobserveElement           = "/actions/vi/observer/unobserve"
)

type Layout interface {
	app.Composer

	ID(v string) Layout
	Class(v string) Layout
	AppName(v string) Layout
	Menu(v ...app.UI) Layout
	MenuOverlay(v ...app.UI) Layout
	Index(v ...app.UI) Layout
	NavBar(v ...app.UI) Layout
	Content(v ...app.UI) Layout
	Footer(v ...app.UI) Layout
	Ads(v ...app.UI) Layout
	Notifications(v ...app.UI) Layout
	NotificationsHighlight(v bool) Layout
	NotificationsIcon(v string) Layout
}

func NewLayout() *layout {
	return &layout{
		IAppName:           "App",
		INotificationsIcon: DefaultNotificationSVG,
	}
}

type layout struct {
	app.Compo

	IID                     string
	IClass                  string
	IAppName                string
	IMenu                   []app.UI
	IMenuOverlay            []app.UI
	IIndex                  []app.UI
	INavBar                 []app.UI
	IContent                []app.UI
	IFooter                 []app.UI
	IAds                    []app.UI
	INotifications          []app.UI
	INotificationsHighlight bool
	INotificationsIcon      string

	showOverlayMenu           bool
	showNotifications         bool
	showInstallButton         bool
	isAppleBrowser            bool
	showAppUpdated            bool
	reloading                 bool
	poppable                  []app.UI
	fullscreenPoppable        bool
	initPoppablePosition      sync.Once
	poppableX                 int
	poppableY                 int
	intersectionObserver      app.Value
	closeIntersectionObserver func()
}

func (l *layout) ID(v string) Layout {
	l.IID = v
	return l
}

func (l *layout) Class(v string) Layout {
	l.IClass = app.AppendClass(l.IClass, v)
	return l
}

func (l *layout) AppName(v string) Layout {
	l.IAppName = v
	return l
}

func (l *layout) Menu(v ...app.UI) Layout {
	l.IMenu = app.FilterUIElems(v...)
	return l
}

func (l *layout) MenuOverlay(v ...app.UI) Layout {
	l.IMenuOverlay = app.FilterUIElems(v...)
	return l
}

func (l *layout) Index(v ...app.UI) Layout {
	l.IIndex = app.FilterUIElems(v...)
	return l
}

func (l *layout) NavBar(v ...app.UI) Layout {
	l.INavBar = app.FilterUIElems(v...)
	return l
}

func (l *layout) Content(v ...app.UI) Layout {
	l.IContent = app.FilterUIElems(v...)
	return l
}

func (l *layout) Footer(v ...app.UI) Layout {
	l.IFooter = app.FilterUIElems(v...)
	return l
}

func (l *layout) Ads(v ...app.UI) Layout {
	l.IAds = app.FilterUIElems(v...)
	return l
}

func (l *layout) Notifications(v ...app.UI) Layout {
	l.INotifications = app.FilterUIElems(v...)
	return l
}

func (l *layout) NotificationsHighlight(v bool) Layout {
	l.INotificationsHighlight = v
	return l
}

func (l *layout) NotificationsIcon(v string) Layout {
	l.INotificationsIcon = v
	return l
}

func (l *layout) OnMount(ctx app.Context) {
	ctx.Handle(showPoppable, l.handleShowPoppable)
	ctx.Handle(hidePoppable, l.handleHidePoppable)
	ctx.Handle(closePoppable, l.handleClosePoppable)
	ctx.Handle(observeElement, l.handleObserveElement)
	ctx.Handle(unobserveElement, l.handleUnobserveElement)

	intersectionObserver := app.Window().Get("IntersectionObserver")
	if !intersectionObserver.Truthy() {
		return
	}
	onIntectionObserverForeach := app.FuncOf(func(this app.Value, args []app.Value) any {
		entry := args[0]
		if !entry.Get("isIntersecting").Bool() {
			return nil
		}
		ctx.NewActionWithValue(notifyObserverIntersection, entry.Get("target"))
		return nil
	})
	onIntersectionObserver := app.FuncOf(func(this app.Value, args []app.Value) any {
		entries := args[0]
		entries.Call("forEach", onIntectionObserverForeach)
		return nil
	})
	l.intersectionObserver = intersectionObserver.New(onIntersectionObserver)
	l.closeIntersectionObserver = func() {
		l.intersectionObserver.Call("disconnect")
		onIntersectionObserver.Release()
		onIntectionObserverForeach.Release()
	}
}

func (l *layout) OnDismount() {
	if l.closeIntersectionObserver != nil {
		l.closeIntersectionObserver()
	}
	l.intersectionObserver = nil
}

func (l *layout) OnNav(ctx app.Context) {
	l.showInstallButton = ctx.IsAppInstallable()
	l.isAppleBrowser = ctx.IsAppleBrowser()
	l.OnAppUpdate(ctx)
	l.poppable = nil
	l.fullscreenPoppable = false
	l.showNotifications = false
}

func (l *layout) OnAppInstallChange(ctx app.Context) {
	appInstallable := ctx.IsAppInstallable()
	if !appInstallable {
		analytics.Track("pwa_install", map[string]interface{}{
			"path": ctx.Page().URL().Path,
		})
	}
	l.showInstallButton = appInstallable
}

func (l *layout) OnAppUpdate(ctx app.Context) {
	l.showAppUpdated = ctx.AppUpdateAvailable()
}

func (l *layout) handleShowPoppable(ctx app.Context, a app.Action) {
	if l.fullscreenPoppable {
		return
	}

	poppable, ok := a.Value.([]app.UI)
	if !ok {
		return
	}

	l.fullscreenPoppable = a.Tags.Get("fullscreen") == "true"
	if l.fullscreenPoppable {
		l.poppable = poppable
		return
	}

	if windowWidth, _ := app.Window().Size(); windowWidth < 1042 {
		return
	}

	l.poppable = poppable
	l.initPoppablePosition.Do(func() {
		l.poppableX, l.poppableY = app.Window().CursorPosition()
	})

	ctx.Defer(func(ctx app.Context) {
		x, y := app.Window().CursorPosition()

		poppableElements := app.Window().
			Get("document").
			Call("getElementsByClassName", "vi-layout-poppable")
		if poppableElements.Length() == 0 {
			return
		}

		poppableElement := poppableElements.Index(0)
		poppableWidth := poppableElement.Get("offsetWidth").Int()
		poppableHeight := poppableElement.Get("offsetHeight").Int()
		windowWidth, windowHeight := app.Window().Size()

		if x+poppableWidth > windowWidth {
			x = windowWidth - poppableWidth
		}
		x = max(0, x)

		if y+poppableHeight > windowHeight {
			y = windowHeight - poppableHeight
		}
		y = max(0, y)

		if x != l.poppableX || y != l.poppableY {
			ctx.Dispatch(func(ctx app.Context) {
				l.poppableX = x
				l.poppableY = y
			})
		}
	})
}

func (l *layout) handleHidePoppable(ctx app.Context, a app.Action) {
	if !l.fullscreenPoppable {
		l.poppable = nil
	}
}

func (l *layout) handleClosePoppable(ctx app.Context, a app.Action) {
	if l.fullscreenPoppable {
		l.poppable = nil
		l.fullscreenPoppable = false
	}
}

func (l *layout) handleObserveElement(ctx app.Context, a app.Action) {
	elem, ok := a.Value.(app.Value)
	if !ok || !elem.Truthy() || l.intersectionObserver == nil {
		return
	}
	l.intersectionObserver.Call("observe", elem)
}

func (l *layout) handleUnobserveElement(ctx app.Context, a app.Action) {
	elem, ok := a.Value.(app.Value)
	if !ok || !elem.Truthy() || l.intersectionObserver == nil {
		return
	}
	l.intersectionObserver.Call("unobserve", elem)
}

func (l *layout) Render() app.UI {
	var withAds string
	if len(l.IAds) != 0 {
		withAds = "vi-layout-with-ads"

	}

	var notificationHighlight string
	if l.INotificationsHighlight {
		notificationHighlight = "vi-layout-navbar-notifications-highlight"
	}

	return app.Div().
		ID(l.IID).
		Class("vi-layout").
		Class(withAds).
		Body(
			// Content
			app.Div().
				Class("vi-layout-menu").
				Body(
					app.Nav().
						Aria("label", "Primary navigation").
						Body(l.IMenu...),
				),
			app.Div().
				Class("vi-layout-index").
				Body(
					app.Nav().
						Aria("label", "Secondary navigation").
						Body(l.IIndex...),
				),
			app.Div().
				Class("vi-layout-content").
				Body(
					app.Div().ID(PageTopID),
					app.Nav().
						Aria("label", "Content navigation").
						Class("vi-layout-navbar").
						Body(
							app.Range(l.INavBar).Slice(func(i int) app.UI {
								return l.INavBar[i]
							}),
							app.Ul().Body(
								NewMediaObject().
									InLi().
									AsButton().
									SVGIcon(SendSVG).
									OnClick(ShareCurrentURL).
									Help("Share"),
								NewMediaObject().
									InLi().
									AsButton().
									Class(notificationHighlight).
									SVGIcon(l.INotificationsIcon).
									OnClick(l.toggleNotifications).
									Help("Notifications"),
								app.If(l.showInstallButton, func() app.UI {
									if !l.isAppleBrowser {
										return NewMediaObject().
											InLi().
											AsButton().
											SVGIcon(DefaultDownloadSVG).
											OnClick(l.installApp).
											Helpf("Install %s", l.IAppName)
									}
									return NewPoppable().
										WithoutHover().
										Content(NewMediaObject().
											InLi().
											AsButton().
											SVGIcon(DefaultDownloadSVG).
											Helpf("Install %s", l.IAppName)).
										Pop(NewBox().
											WithoutSelection().
											WithGap().
											Class("vi-layout-apple-poppable-box").
											Content(
												NewMediaObject().
													AsH2().
													WithHorizontalCentering().
													WithExtraLargeMedia().
													SVGIcon(AppleSVG).
													Textf("Install %s on IPhone, IPad & Mac", l.IAppName),
												app.P().Textf("You can install %s as a Progressive Web App (PWA) on your Apple device to open it directly from your Home Screen or Dock in its own window, just like a native app.", l.IAppName),
												app.H3().Text("iPhone or iPad"),
												app.Ol().Body(
													app.Li().Body(
														app.Text("On the right of the address bar, tap the "),
														app.Span().
															Class("accent").
															Text("…"),
														app.Text(" button"),
													),
													app.Li().Body(
														app.Text("Choose "),
														app.Span().
															Class("accent").
															Text("Share"),
													),
													app.Li().Body(
														app.Text("Select "),
														app.Span().
															Class("accent").
															Text("Add to Home Screen"),
													),
													app.Li().Body(
														app.Text("Confirm with "),
														app.Span().
															Class("accent").
															Text("Add"),
														app.Text(" or "),
														app.Span().
															Class("accent").
															Text("✓"),
													),
												),
												app.H3().Text("MacOS (Safari)"),
												app.Ol().Body(
													app.Li().Body(
														app.Text("On the right of the address bar, click on the "),
														app.Span().
															Class("accent").
															Text("Share"),
														app.Text(" button"),
													),
													app.Li().Body(
														app.Text("Choose "),
														app.Span().
															Class("accent").
															Text("Add to Dock"),
														app.Text(" from the menu"),
													),
													app.Li().Body(
														app.Text("In the dialog that appears, click "),
														app.Span().
															Class("accent").
															Text("Add"),
													),
												),
											))
								}),
							),
						),
					app.Main().Body(l.IContent...),
					app.Footer().Body(l.IFooter...),
				),
			app.If(len(l.IAds) != 0, func() app.UI {
				return app.Aside().
					Class("vi-layout-ads").
					Body(l.IAds...)
			}),

			// Poppable
			app.If(app.IsClient, func() app.UI {
				var showPoppable string
				if len(l.poppable) != 0 {
					showPoppable = "vi-layout-poppable-show"
				}

				var fullscreenPoppable string
				if l.fullscreenPoppable {
					fullscreenPoppable = "vi-layout-poppable-fullscreen"
				}

				return app.Aside().
					Class("vi-layout-poppable").
					Class(showPoppable).
					Class(fullscreenPoppable).
					Style("transform", fmt.Sprintf("translate(%vpx, %vpx)", l.poppableX, l.poppableY)).
					OnClick(l.closeFullscreenPoppable).
					Body(
						app.Div().Body(
							app.Range(l.poppable).Slice(func(i int) app.UI {
								return l.poppable[i]
							}),
							NewMediaObject().
								AsButton().
								Class("vi-layout-poppable-close-button").
								SVGIcon(DefaultCloseSVG).
								Help("Close"),
						),
					)
			}),

			// Overlay Menu
			app.Nav().
				Class("vi-layout-menu-hamburger").
				Body(
					NewMediaObject().
						AsButton().
						SVGIcon(DefaultAppMenuSVG).
						OnClick(l.toggleOverlayMenu).
						Help("Menu"),
				),
			app.If(app.IsClient, func() app.UI {
				var showOverlayMenu string
				if l.showOverlayMenu {
					showOverlayMenu = "vi-layout-menu-overlay-show"
				}

				return app.Div().
					Class("vi-layout-menu-overlay").
					Class(showOverlayMenu).
					OnClick(l.toggleOverlayMenu).
					Body(
						app.Div().Body(
							app.Nav().
								Aria("label", "Small screen primary navigation").
								Body(l.IMenuOverlay...),
						),
					)
			}),

			// Notifications
			app.If(app.IsClient, func() app.UI {
				var showNotifications string
				if l.showNotifications {
					showNotifications = "vi-layout-notifications-show"
				}

				return app.Div().
					Class("vi-layout-notifications").
					Class(showNotifications).
					OnClick(l.toggleNotifications).
					Body(
						app.Div().Body(
							app.Aside().Body(l.INotifications...),
						),
					)
			}),

			// App Update
			app.If(l.showAppUpdated, func() app.UI {
				return app.Aside().
					Class("vi-layout-app-update").
					OnClick(l.updateApp).
					Body(
						app.If(!l.reloading, func() app.UI {
							return NewSVG().Data(DefaultDownloadSVG)
						}).Else(func() app.UI {
							return NewSpinner()
						}),
						app.Div().Body(
							app.H2().Textf("%s has been updated!", l.IAppName),
							app.If(!l.reloading, func() app.UI {
								return app.P().Text("Click/Tap to reload")
							}).Else(func() app.UI {
								return app.P().Text("Reloading...")
							}),
						),
					)
			}),
		)
}

func (l *layout) toggleOverlayMenu(ctx app.Context, e app.Event) {
	l.showOverlayMenu = !l.showOverlayMenu
}

func (l *layout) toggleNotifications(ctx app.Context, e app.Event) {
	l.showNotifications = !l.showNotifications
}

func (l *layout) installApp(ctx app.Context, e app.Event) {
	ctx.ShowAppInstallPrompt()
}

func (l *layout) updateApp(ctx app.Context, e app.Event) {
	l.reloading = true
	ctx.Defer(func(ctx app.Context) {
		ctx.Reload()
	})
}

func (l *layout) closeFullscreenPoppable(ctx app.Context, e app.Event) {
	ctx.NewAction(closePoppable)
}

func ShareCurrentURL(ctx app.Context, e app.Event) {
	e.PreventDefault()
	currentURL := app.Window().Get("document").Get("URL").String()

	shareNotSupported := func() {
		navigator := app.Window().Get("navigator")
		if !navigator.Truthy() || !navigator.Get("clipboard").Truthy() {
			app.Window().Call("alert", "Sharing is not supported on your current device or browser.")
			return
		}

		navigator.Get("clipboard").Call("writeText", currentURL)
		app.Window().Call("alert", "Copied to Clipboard.")
	}

	navigator := app.Window().Get("navigator")
	if currentURL == "" || !navigator.Truthy() {
		shareNotSupported()
		return
	}

	canShare := navigator.Get("canShare")
	if !canShare.Truthy() {
		shareNotSupported()
		return
	}

	shareConfig := map[string]any{
		"url": currentURL,
	}

	if !navigator.Call("canShare", shareConfig).Truthy() {
		shareNotSupported()
		return
	}

	navigator.Call("share", shareConfig)
}

func ScrollToTop(ctx app.Context) {
	ctx.ScrollTo(PageTopID)
}

func ScrollToFragment(ctx app.Context) {
	if fragment := ctx.Page().URL().Fragment; fragment != "" {
		ctx.ScrollTo(fragment)
		return
	}
	ScrollToTop(ctx)
}

func NavigateBack(ctx app.Context, e app.Event) {
	e.PreventDefault()
	e.StopImmediatePropagation()
	app.Window().Get("history").Call("back")
}
