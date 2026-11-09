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

// Markdown is a component that renders Markdown as HTML with vi-markdown styling.
// Its container can be a div, section, or aside with a custom ID and CSS classes.
type Markdown interface {
	app.Composer

	AsSection() Markdown
	AsAside() Markdown
	ID(v string) Markdown
	Class(v ...string) Markdown
}

// NewMarkdown creates a Markdown component from md, trimming leading and trailing
// whitespace. The component uses a div container by default.
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

// MarkdownDocument holds document metadata extracted from Markdown source by a
// [MarkdownReader]. Its field comments describe the supported metadata keys.
type MarkdownDocument struct {
	// Title is the value of vi:title.
	Title string

	// Author is the value of vi:author.
	Author string

	// Description is the value of vi:description.
	Description string

	// Image is the value of vi:image.
	Image string

	// OGImage is the value of vi:og-image.
	OGImage string

	// ModifiedAt is the value of vi:modified-at, parsed as RFC3339 or YYYY-MM-DD.
	ModifiedAt time.Time
}

// MarkdownReader collects metadata from Markdown source into a [MarkdownDocument].
// Metadata comments use a key and value separated by a space or tab, such as
// <!-- vi:description A short description -->. See [MarkdownDocument] for
// supported keys and value formats.
type MarkdownReader struct {
	markdown string
	document MarkdownDocument
}

func (r *MarkdownReader) Document() MarkdownDocument {
	return r.document
}

func (r *MarkdownReader) Read(md string) {
	for line := range strings.SplitSeq(md, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "<!--"):

		case strings.HasPrefix(line, "# "):
		}
	}
}

func (r *MarkdownReader) readComment(comment string) {
	comment, ok := strings.CutPrefix(comment, "<!--")
	if !ok {
		return
	}
	comment, ok = strings.CutSuffix(comment, "-->")
	if !ok ||
		strings.Contains(comment, "<!--") ||
		strings.Contains(comment, "-->") {
		return
	}
	comment = strings.TrimSpace(comment)

	separator := strings.IndexAny(comment, " \t")
	if separator == -1 {
		return
	}

	key := comment[:separator]
	value := strings.TrimSpace(comment[separator+1:])
	switch key {
	case "vi:title":
		r.document.Title = value

	case "vi:author":
		r.document.Author = value

	case "vi:description":
		r.document.Description = value

	case "vi:image":
		r.document.Image = value

	case "vi:og-image":
		r.document.OGImage = value

	case "vi:modified-at":
		for _, layout := range []string{
			time.RFC3339,
			time.DateOnly,
		} {
			if date, err := time.Parse(layout, value); err == nil {
				r.document.ModifiedAt = date
				break
			}
		}
	}
}

func (r *MarkdownReader) readH1(v string) {

}
