package vi

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
)

type Markdown interface {
	app.Composer

	AsSection() Markdown
	AsAside() Markdown
	ID(v string) Markdown
	Class(v ...string) Markdown
}

func NewMarkdown(md string) Markdown {
	return &markdown{
		IRoot:     "div",
		IMarkdown: strings.TrimSpace(md),
	}
}

type markdown struct {
	app.Compo

	IRoot     string
	IID       string
	IClass    string
	IMarkdown string
}

func (m *markdown) AsSection() Markdown {
	m.IRoot = "section"
	return m
}

func (m *markdown) AsAside() Markdown {
	m.IRoot = "aside"
	return m
}

func (m *markdown) ID(v string) Markdown {
	m.IID = v
	return m
}

func (m *markdown) Class(v ...string) Markdown {
	m.IClass = app.AppendClass(m.IClass, v...)
	return m
}

func (m *markdown) Render() app.UI {
	md, err := ParseMarkdown(m.IMarkdown)
	if err != nil {
		md = fmt.Errorf(`<span class="error-color">Decoding markdown failed: %w.</span>`, err).Error()
	}

	return app.Raw(fmt.Sprintf(
		`<%s id=%q class="vi-markdown %s">%s</%s>`,
		m.IRoot,
		m.IID,
		m.IClass,
		md,
		m.IRoot,
	))
}

func ParseMarkdown(md string) (string, error) {
	var b bytes.Buffer
	var source = []byte(md)

	p := parser.New()
	r := html.New()
	doc := p.Parse(source)
	if err := r.Render(&b, source, doc); err != nil {
		return "", err
	}

	return b.String(), nil
}

type MarkdownMeta struct {
	Title       string
	Description string
	Image       string
	OGImage     string
	UpdatedAt   time.Time
}

func ParseMarkdownMeta(md string) MarkdownMeta {
	return MarkdownMeta{}
}
