package vi

import (
	"fmt"
	"strconv"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// Stepper is the interface that describes a numeric control with
// increment and decrement actions.
type Stepper interface {
	app.Composer

	// Sets the ID.
	ID(v string) Stepper

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Stepper

	// Sets the icon from the given SVG data.
	SVGIcon(v string) Stepper

	// Sets the icon.
	Icon(v app.UI) Stepper

	// Sets the minimum value, maximum value, and increment step.
	Range(min, max, step float64) Stepper

	// Sets the current value.
	Value(v float64) Stepper

	// Sets the text label.
	Text(v string) Stepper

	// Sets the text label with the given format and value.
	Textf(format string, v ...any) Stepper

	// Sets the displayed unit.
	Unit(v string) Stepper

	// Sets the displayed unit with the given format and value.
	Unitf(format string, v ...any) Stepper

	// Sets the help text.
	Help(v string) Stepper

	// Sets the help text with the given format and value.
	Helpf(format string, v ...any) Stepper

	// Sets the function used to convert the input value.
	ConvertValue(v func(string) float64) Stepper

	// Sets the handler called when the value changes.
	OnChange(v func(app.Context, float64)) Stepper
}

// Creates a stepper.
func NewStepper() Stepper {
	return &stepper{
		IMin:  0,
		IMax:  100,
		IStep: 1,
	}
}

type stepper struct {
	app.Compo

	IID           string
	IClass        string
	IText         string
	IUnit         string
	IHelp         string
	IMin          float64
	IMax          float64
	IStep         float64
	IValue        float64
	IIcon         app.UI
	IConvertValue func(string) float64
	IOnChange     func(app.Context, float64)
}

func (s *stepper) ID(v string) Stepper {
	s.IID = v
	return s
}

func (s *stepper) Class(v ...string) Stepper {
	s.IClass = app.AppendClass(s.IClass, v...)
	return s
}

func (s *stepper) SVGIcon(v string) Stepper {
	return s.Icon(NewSVG().Data(v))
}

func (s *stepper) Icon(v app.UI) Stepper {
	if elems := app.FilterUIElems(v); len(elems) != 0 {
		s.IIcon = elems[0]
	}
	return s
}

func (s *stepper) Range(min, max, step float64) Stepper {
	if step == 0 {
		step = 0.1
	}

	if min <= 0 {
		min = 0
	}

	if max <= min {
		max = min
	}

	s.IMin = min
	s.IMax = max
	s.IStep = step
	return s
}

func (s *stepper) Value(v float64) Stepper {
	s.IValue = v
	return s
}

func (s *stepper) Text(v string) Stepper {
	s.IText = v
	return s
}

func (s *stepper) Textf(format string, v ...any) Stepper {
	return s.Text(fmt.Sprintf(format, v...))
}

func (s *stepper) Unit(v string) Stepper {
	s.IUnit = v
	return s
}

func (s *stepper) Unitf(format string, v ...any) Stepper {
	return s.Unit(fmt.Sprintf(format, v...))
}

func (s *stepper) Help(v string) Stepper {
	s.IHelp = v
	return s
}

func (s *stepper) Helpf(format string, v ...any) Stepper {
	return s.Help(fmt.Sprintf(format, v...))
}

func (s *stepper) ConvertValue(v func(string) float64) Stepper {
	s.IConvertValue = v
	return s
}

func (s *stepper) OnChange(v func(app.Context, float64)) Stepper {
	s.IOnChange = v
	return s
}

func (s *stepper) Render() app.UI {
	return app.Div().
		ID(s.IID).
		Class("vi-stepper").
		Class(s.IClass).
		Role("group").
		OnClick(s.whenClicked).
		Body(
			app.Div().
				Class("vi-stepper-content").
				Body(
					app.If(s.IIcon == nil, func() app.UI {
						return NewSVG().Data(SVGIconByTitle(s.IText))
					}).Else(func() app.UI {
						return s.IIcon
					}),
					app.Div().
						Stylef("height", "%.3f%%", s.IValue*100/s.IMax).
						Role("img"),
					app.Div().
						Class("h3").
						Text(s.IText),
					app.Input().
						Class("h3").
						Type("text").
						Pattern("[0-9]*[a-zA-Z]*").
						AutoComplete(false).
						Value(s.IUnit).
						OnFocus(s.whenUnitFocus).
						OnChange(s.editValue),
				),
			app.Div().
				Class(Disabled(s.IValue >= s.IMax)).
				OnClick(s.increase).
				Body(
					NewMediaObject().
						AsButton().
						SVGIcon(MenuUpSVG),
				),
			app.Div().
				Class(Disabled(s.IValue <= s.IMin)).
				OnClick(s.decrease).
				Body(
					NewMediaObject().
						AsButton().
						SVGIcon(MenuDownSVG),
				),
			app.Div(),
			app.Div(),
		)
}

func (s *stepper) increase(ctx app.Context, e app.Event) {
	updatedValue := s.IValue + s.IStep
	if updatedValue > s.IMax {
		updatedValue = s.IMax
	}
	if s.IOnChange != nil {
		s.IOnChange(ctx, updatedValue)
	}
}

func (s *stepper) decrease(ctx app.Context, e app.Event) {
	updatedValue := s.IValue - s.IStep
	if updatedValue < s.IMin {
		updatedValue = s.IMin
	}
	if s.IOnChange != nil {
		s.IOnChange(ctx, updatedValue)
	}
}

func (s *stepper) whenUnitFocus(ctx app.Context, e app.Event) {
	target := e.Get("target")
	if !target.Truthy() {
		return
	}
	target.Call("select")
}

func (s *stepper) editValue(ctx app.Context, e app.Event) {
	var updatedValue float64

	value := ctx.JSSrc().Get("value").String()
	if s.IConvertValue == nil {
		updatedValue, _ = strconv.ParseFloat(value, 64)
	} else {
		updatedValue = s.IConvertValue(value)
	}

	if updatedValue < s.IMin {
		updatedValue = s.IMin
	} else if updatedValue > s.IMax {
		updatedValue = s.IMax
	}

	if s.IOnChange != nil {
		s.IOnChange(ctx, updatedValue)
	}
}

func (s *stepper) whenClicked(ctx app.Context, e app.Event) {
	e.PreventDefault()
	e.StopImmediatePropagation()
}
