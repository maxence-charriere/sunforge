package vi

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// A spinner.
type Spinner interface {
	app.UI

	// Sets the ID.
	ID(v string) Spinner

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Spinner

	// Sets the spinner to be small.
	Small() Spinner
}

// Creates a spinner.
func NewSpinner() Spinner {
	return &spinner{}
}

type spinner struct {
	app.Compo

	IID      string
	IClass   string
	IDisplay string
}

func (s *spinner) ID(v string) Spinner {
	s.IID = v
	return s
}

func (s *spinner) Class(v ...string) Spinner {
	s.IClass = app.AppendClass(s.IClass, v...)
	return s
}

func (s *spinner) Small() Spinner {
	s.IDisplay = "vi-small"
	return s
}

func (s *spinner) Render() app.UI {
	return app.Div().
		DataSet("vi", "spinner").
		ID(s.IID).
		Class("vi-spinner").
		Class(s.IDisplay).
		Class(s.IClass).
		Body(
			app.Div(),
			app.Div(),
			app.Div(),
		)
}
