package vi

import (
	"testing"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/stretchr/testify/require"
)

func TestCrucibleObserverActions(t *testing.T) {
	c := &crucibleObserverHarness{crucible: &crucible{}}
	engine := app.NewTestEngine()
	require.NoError(t, engine.Load(c))
	engine.ConsumeAll()

	target := &crucibleObserverValue{truthy: true}
	var calls []string
	c.observer = &crucibleObserverValue{
		call: func(method string, args ...any) app.Value {
			require.Len(t, args, 1)
			require.Same(t, target, args[0])
			calls = append(calls, method)
			return nil
		},
	}

	for _, action := range []string{observeElement, unobserveElement} {
		c.ctx.NewActionWithValue(action, target)
		engine.ConsumeAll()
	}
	require.Equal(t, []string{"observe", "unobserve"}, calls)

	for _, invalid := range []any{nil, "not an element", &crucibleObserverValue{}} {
		c.ctx.NewActionWithValue(observeElement, invalid)
		c.ctx.NewActionWithValue(unobserveElement, invalid)
	}
	engine.ConsumeAll()
	require.Len(t, calls, 2)

	// An unavailable or disconnected browser observer must be a no-op.
	c.observer = nil
	c.ctx.NewActionWithValue(observeElement, target)
	c.ctx.NewActionWithValue(unobserveElement, target)
	require.NotPanics(t, engine.ConsumeAll)
}

func TestCrucibleObserverDismount(t *testing.T) {
	c := &crucible{}
	var dismounter app.Dismounter = c
	require.NotPanics(t, dismounter.OnDismount)

	closed := 0
	c.observer = &crucibleObserverValue{}
	c.closeObserver = func() { closed++ }
	dismounter.OnDismount()
	dismounter.OnDismount()

	require.Equal(t, 1, closed)
	require.Nil(t, c.observer)
	require.Nil(t, c.closeObserver)
}

// Run mount registration through the native test engine's server lifecycle.
type crucibleObserverHarness struct {
	*crucible
	ctx app.Context
}

func (c *crucibleObserverHarness) OnPreRender(ctx app.Context) {
	c.ctx = ctx
	c.OnMount(ctx)
}

func (c *crucibleObserverHarness) Render() app.UI {
	return app.Div()
}

type crucibleObserverValue struct {
	app.Value
	truthy bool
	call   func(string, ...any) app.Value
}

func (v *crucibleObserverValue) Truthy() bool {
	return v.truthy
}

func (v *crucibleObserverValue) Call(method string, args ...any) app.Value {
	return v.call(method, args...)
}
