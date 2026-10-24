package vi

import (
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// Lockable is the interface that describes content that can be wrapped
// with a locked state.
type Lockable interface {
	app.Composer

	// Sets the ID.
	ID(v string) Lockable

	// Sets the class. Multiple classes can be defined by successive calls.
	Class(v ...string) Lockable

	// Sets whether the content is locked.
	Locked(v bool) Lockable

	// Sets the main content.
	Content(v app.UI) Lockable

	// Sets the content displayed when locked.
	LockedContent(v ...app.UI) Lockable
}

// Creates a lockable content wrapper.
func NewLockable() Lockable {
	return &lockable{}
}

type lockable struct {
	app.Compo

	IID            string
	IClass         string
	ILocked        bool
	IContent       app.UI
	ILockedContent []app.UI
}

func (l *lockable) ID(v string) Lockable {
	l.IID = v
	return l
}

func (l *lockable) Class(v ...string) Lockable {
	l.IClass = app.AppendClass(l.IClass, v...)
	return l
}

func (l *lockable) Locked(v bool) Lockable {
	l.ILocked = v
	return l
}

func (l *lockable) Content(v app.UI) Lockable {
	if elems := app.FilterUIElems(v); len(elems) != 0 {
		l.IContent = elems[0]
	}
	return l
}

func (l *lockable) LockedContent(v ...app.UI) Lockable {
	l.ILockedContent = app.FilterUIElems(v...)
	return l
}

func (l *lockable) Render() app.UI {
	if !l.ILocked {
		return l.IContent
	}

	return NewPoppable().
		WithoutHover().
		ID(l.IID).
		Class("vi-lockable").
		Class(l.IClass).
		Content(l.IContent).
		Pop(l.ILockedContent...)
}
