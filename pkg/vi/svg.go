package vi

import (
	"fmt"
	"strings"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

var (
	// The default SVG used when a loaders get an error.
	DefaultErrorSVG = `<svg $id $class $help" viewBox="0 0 24 24">
    	<path fill="currentColor" d="M22 14H21C21 10.13 17.87 7 14 7H13V5.73C13.6 5.39 14 4.74 14 4C14 2.9 13.11 2 12 2S10 2.9 10 4C10 4.74 10.4 5.39 11 5.73V7H10C6.13 7 3 10.13 3 14H2C1.45 14 1 14.45 1 15V18C1 18.55 1.45 19 2 19H3V20C3 21.11 3.9 22 5 22H19C20.11 22 21 21.11 21 20V19H22C22.55 19 23 18.55 23 18V15C23 14.45 22.55 14 22 14M9.86 16.68L8.68 17.86L7.5 16.68L6.32 17.86L5.14 16.68L6.32 15.5L5.14 14.32L6.32 13.14L7.5 14.32L8.68 13.14L9.86 14.32L8.68 15.5L9.86 16.68M18.86 16.68L17.68 17.86L16.5 16.68L15.32 17.86L14.14 16.68L15.32 15.5L14.14 14.32L15.32 13.14L16.5 14.32L17.68 13.14L18.86 14.32L17.68 15.5L18.86 16.68Z" />
	</svg>`

	// The default SVG used for app install.
	DefaultInstallSVG = `<svg $id $class $help viewBox="0 0 24 24">
    	<path fill="currentColor" d="M21 16H3V4H21M21 2H3C1.9 2 1 2.9 1 4V16C1 17.1 1.9 18 3 18H10V20H8V22H16V20H14V18H21C22.1 18 23 17.1 23 16V4C23 2.9 22.1 2 21 2M16 10H13V6H11V10H8L12 14L16 10" />
	</svg>`
)

// The interface that describes a SVG.
type SVG interface {
	app.UI

	// Sets the ID.
	ID(v string) SVG

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) SVG

	// Sets the svg data.
	Data(v string) SVG

	// Sets the help with the given format and value.
	Help(v string) SVG
	Helpf(format string, v ...any) SVG
}

// Creates a new SVG.
func NewSVG() SVG {
	return &svg{}
}

type svg struct {
	app.Compo

	IID    string
	IClass string
	IData  string
	IHelp  string
}

func (s *svg) ID(v string) SVG {
	s.IID = v
	return s
}

func (s *svg) Class(v ...string) SVG {
	s.IClass = app.AppendClass(s.IClass, v...)
	return s
}

func (s *svg) Data(v string) SVG {
	s.IData = v
	return s
}

func (s *svg) Help(v string) SVG {
	s.IHelp = v
	return s
}

func (s *svg) Helpf(format string, v ...any) SVG {
	return s.Help(fmt.Sprintf(format, v...))
}

func (s *svg) Render() app.UI {
	data := s.IData

	if s.IID != "" {
		data = strings.ReplaceAll(data, "$id", fmt.Sprintf(`id="%s"`, s.IClass))
	} else {
		data = strings.ReplaceAll(data, "$id", "")
	}

	if s.IClass != "" {
		data = strings.ReplaceAll(data, "$class", fmt.Sprintf(`class="%s"`, s.IClass))
	} else {
		data = strings.ReplaceAll(data, "$class", "")
	}

	if s.IHelp != "" {
		data = strings.ReplaceAll(data, "$help", fmt.Sprintf(`title="%s"`, s.IClass))
	} else {
		data = strings.ReplaceAll(data, "$help", "")
	}

	return app.Raw(data)
}
