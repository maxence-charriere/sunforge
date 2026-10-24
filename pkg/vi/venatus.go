package vi

import (
	"fmt"
	"sync"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/maxence-charriere/go-app/v11/pkg/errors"
	"github.com/maxence-charriere/go-app/v11/pkg/logs"
)

type VenatusConfig struct {
	Video       string
	Rectangle   string
	Leaderboard string
	Skycraper   string
}

// Venatus ad provider.
type Venatus struct {
	currentURLPath string
	loaded         map[string]bool
}

// Creates a new Venatus ad provider.
func NewVenatus() *Venatus {
	return &Venatus{
		loaded: map[string]bool{},
	}
}

func (p *Venatus) InitScript() app.HTMLScript {
	return nil
}

func (p *Venatus) OnNav(ctx app.Context) {
	if path := ctx.Page().URL().Path; app.IsClient && path != p.currentURLPath {
		if p.currentURLPath != "" {
			venatusMediaPush(func(scope app.Value) {
				app.Log(logs.New("creating a new venatus page session"))
				scope.Get("Instances").
					Get("pageManager").
					Call("newCleanPageSession")
				clear(p.loaded)
			})
		}
		p.currentURLPath = path
	}
}

func (p *Venatus) NewAdUnit(placement string) Ad {
	return newVenatusAdUnit(placement, p.displayAd)
}

func (p *Venatus) displayAd(tag, target string) {
	venatusMediaPush(func(scope app.Value) {
		if p.loaded[target] {
			return
		}
		p.loaded[target] = true

		switch tag {
		case "video":
			tag = "video"

		case "skyscraper":
			tag = "double_mpu"

		case "leaderboard":
			tag = "leaderboard"

		default:
			tag = "mpu"
		}

		config := scope.Get("Config").Call("get", tag)
		if !config.Truthy() {
			return
		}

		defer func() {
			if err := recover(); err != nil {
				app.Log(errors.New("loading venatus ad failed").
					WithTag("type", tag).
					WithTag("target", target).
					WithTag("reason", err))
			}
		}()

		app.Log(logs.New("loading venatus ad").
			WithTag("type", tag).
			WithTag("target", target))
		config.Call("display", target)
	})
}

type venatusAdUnit struct {
	app.Compo

	IID          string
	IClass       string
	IPlacement   string
	IMode        string
	IPlaceholder []app.UI

	loadOnce       sync.Once
	loaded         bool
	loadingIcon    bool
	currentURLPath string
	displayAd      func(tage, target string)
}

func newVenatusAdUnit(placement string, displayAd func(string, string)) *venatusAdUnit {
	return &venatusAdUnit{
		IPlacement: fmt.Sprintf("%s-slot", placement),
		IMode:      "rectangle",
		displayAd:  displayAd,
	}
}

func (u *venatusAdUnit) AsVideo() Ad {
	u.IMode = "video"
	return u
}

func (u *venatusAdUnit) AsRectangle() Ad {
	u.IMode = "rectangle"
	return u
}

func (u *venatusAdUnit) AsLeaderboard() Ad {
	u.IMode = "leaderboard"
	return u
}

func (u *venatusAdUnit) AsSkyscraper() Ad {
	u.IMode = "skyscraper"
	return u
}

func (u *venatusAdUnit) ID(v string) Ad {
	u.IID = v
	return u
}

func (u *venatusAdUnit) IDf(format string, v ...any) Ad {
	return u.ID(fmt.Sprintf(format, v...))
}

func (u *venatusAdUnit) Class(v ...string) Ad {
	u.IClass = app.AppendClass(u.IClass, v...)
	return u
}

func (u *venatusAdUnit) Placeholder(v ...app.UI) Ad {
	u.IPlaceholder = app.FilterUIElems(v...)
	return u
}

func (u *venatusAdUnit) OnNav(ctx app.Context) {
	if app.IsServer || !u.loaded {
		return
	}

	path := ctx.Page().URL().Path
	if path != u.currentURLPath {
		u.loadOnce = sync.Once{}
		u.load(ctx)
		u.currentURLPath = path
	}
}

func (u *venatusAdUnit) Render() app.UI {
	return NewObserver().
		AsAside().
		ID(u.IID).
		Class("vi-vm").
		Class("vi-vm-"+u.mode()).
		Class(u.IClass).
		OnBecomeVisible(u.onBecomeVisible).
		Content(
			app.Div().ID(u.IPlacement),
			app.Div().
				Class("vi-vm-placeholder").
				Body(
					app.IfSlice(u.loadingIcon, func() []app.UI {
						if len(u.IPlaceholder) != 0 {
							return u.IPlaceholder
						}
						return []app.UI{
							NewSpinner().
								Small(),
						}
					}),
				),
		)
}

func (u *venatusAdUnit) mode() string {
	if width, _ := app.Window().Size(); width < 970 && u.IMode != "video" {
		return "rectangle"
	}
	return u.IMode
}

func (u *venatusAdUnit) onBecomeVisible(ctx app.Context) {
	u.load(ctx)
}

func (u *venatusAdUnit) load(ctx app.Context) {
	u.loadOnce.Do(func() {
		u.loadingIcon = true
		ctx.After(time.Second*5, func(ctx app.Context) {
			u.loadingIcon = false
		})

		ctx.Defer(func(ctx app.Context) {
			u.displayAd(u.mode(), u.IPlacement)
			u.loaded = true
		})
	})
}

func venatusMedia() app.Value {
	if vm := app.Window().Get("__VM"); !vm.Truthy() {
		app.Window().Set("__VM", []any{})
	}
	return app.Window().Get("__VM")
}

func venatusMediaPush(v func(scope app.Value)) {
	var function app.Func
	function = app.FuncOf(func(this app.Value, args []app.Value) any {
		defer function.Release()
		v(args[1])
		return nil
	})
	venatusMedia().Call("push", function)
}
