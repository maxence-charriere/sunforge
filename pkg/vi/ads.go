package vi

import "github.com/maxence-charriere/go-app/v11/pkg/app"

type AdProvider interface {
	// InitScript returns the initialization script that loads the ad provider's
	// code.
	InitScript() app.HTMLScript

	// OnNav is a method to perform initializations that occur when a page is
	// navigated.
	OnNav(app.Context)

	// NewAdUnit creates a UI element from the given ad unit.
	NewAdUnit(placement string) Ad
}

// An ad UI element.
type Ad interface {
	app.UI

	// Sets the ID.
	ID(v string) Ad
	IDf(format string, v ...any) Ad

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Ad

	// Sets the ad as a video.
	AsVideo() Ad

	// Sets the ad as a rectangle.
	AsRectangle() Ad

	// Sets the ad as a leaderboard.
	AsLeaderboard() Ad

	// Sets the ad as skycraper.
	AsSkyscraper() Ad

	// Sets the placeholder displays when the ad is unfilled.
	Placeholder(v ...app.UI) Ad
}

// An ad manager.
type AdManager struct {
	providers []AdProvider
}

// Creates a new ad manager that uses the given providers. Ad providers are
// ordered by priority..
func NewAdManager(providers ...AdProvider) *AdManager {
	return &AdManager{providers: providers}
}

// InitScripts returns the initialization scripts required by each ad provider.
func (m *AdManager) InitScripts() []app.UI {
	scripts := make([]app.UI, 0, len(m.providers))
	for _, p := range m.providers {
		if s := p.InitScript(); s != nil {
			scripts = append(scripts, s)
		}
	}
	return scripts
}

// OnNav calls each provider's OnNav method.
func (m *AdManager) OnNav(ctx app.Context) {
	if app.IsServer {
		return
	}

	for _, p := range m.providers {
		p.OnNav(ctx)
	}
}

// NewAdUnit creates a UI element from the given ad unit. It used the first
// non-nil ad unit returned by the given ad providers.
func (m *AdManager) NewAdUnit(placement string) Ad {
	for _, p := range m.providers {
		if adUnit := p.NewAdUnit(placement); adUnit != nil {
			return adUnit
		}
	}
	return nil
}
