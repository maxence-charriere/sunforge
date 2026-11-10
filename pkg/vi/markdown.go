package vi

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/maxence-charriere/go-app/v11/pkg/app"
	"github.com/yuin/goldmark/v2/ast"
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

// ParseMarkdown renders Markdown as HTML with automatically generated heading IDs.
func ParseMarkdown(md string) (string, error) {
	var b bytes.Buffer
	var source = []byte(md)

	p := parser.New(parser.WithAutoHeadingID())
	r := html.New()
	doc := p.Parse(source)
	if err := r.Render(&b, source, doc); err != nil {
		return "", err
	}

	return b.String(), nil
}

// MarkdownDocument holds metadata and content extracted from Markdown source by a
// [MarkdownReader]. Its field comments describe the supported metadata keys.
type MarkdownDocument struct {
	// Title is the text of a top-level H1 or the value of vi:title, whichever occurs
	// last. Inline Markdown is preserved.
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

	// Intro is the original Markdown after the first top-level H1 and before
	// the first H2, or EOF. It is empty when no H1 precedes the first H2.
	Intro string

	// Sections contains each top-level H2 title and its Markdown body.
	Sections []MarkdownDocumentSection
}

// MarkdownDocumentSection separates a section heading from its Markdown body.
type MarkdownDocumentSection struct {
	// Title is the H2 text without heading markers. Inline Markdown is preserved.
	Title string

	// Content is the Markdown after the H2 and before the next top-level H2 or EOF.
	// Boundary newlines are trimmed; indentation and other whitespace are preserved.
	Content string
}

// MarkdownReader uses Goldmark to collect top-level headings, metadata comments,
// and content into a [MarkdownDocument] returned by Read. Both ATX and Setext
// headings are supported. Each read starts a new document.
// Metadata comments use a key and value separated by a space or tab, such as
// <!-- vi:description A short description -->. See [MarkdownDocument] for
// supported keys and value formats.
type MarkdownReader struct {
	document MarkdownDocument
}

func (r *MarkdownReader) Read(md string) MarkdownDocument {
	r.document = MarkdownDocument{}
	source := []byte(md)
	doc := parser.New().Parse(source)
	introStart := -1
	introEnd := len(md)
	sectionStart := -1
	var sectionHeading *ast.Heading

	// Only document-level blocks define metadata and section boundaries.
	// Goldmark keeps code, HTML, lists, and blockquotes in their own nodes.
	for node := range doc.Children() {
		switch node := node.(type) {
		case *ast.HTMLBlock:
			r.readComment(node, source)

		case *ast.Heading:
			switch node.Level {
			case 1:
				r.readH1(node, source)
				if introStart < 0 && sectionStart < 0 {
					introStart = markdownHeadingEnd(md, node)
				}

			case 2:
				// Include the heading's original indentation in the source slice.
				start := strings.LastIndexByte(md[:node.Pos()], '\n') + 1
				if sectionStart < 0 {
					introEnd = start
				} else {
					r.readSection(sectionHeading, source, md[sectionStart:start])
				}
				sectionHeading = node
				sectionStart = markdownHeadingEnd(md, node)
			}
		}
	}

	if introStart >= 0 && introStart < introEnd {
		r.document.Intro = md[introStart:introEnd]
	}
	if sectionStart >= 0 {
		r.readSection(sectionHeading, source, md[sectionStart:])
	}
	return r.document
}

// markdownHeadingEnd returns the position after the heading's final source line.
func markdownHeadingEnd(md string, heading *ast.Heading) int {
	position := heading.Pos()
	if segments := heading.Source(); len(segments) != 0 {
		position = segments[len(segments)-1].Start
	}
	lines := 1
	if heading.HeadingKind == ast.HeadingKindSetext {
		lines++ // Include the underline after the last line of heading text.
	}
	for range lines {
		end := strings.IndexByte(md[position:], '\n')
		if end < 0 {
			return len(md)
		}
		position += end + 1
	}
	return position
}

func (r *MarkdownReader) readComment(node *ast.HTMLBlock, source []byte) {
	comment := strings.TrimSpace(node.Value.Str(source))
	if strings.ContainsAny(comment, "\r\n") {
		return
	}
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

func (r *MarkdownReader) readH1(heading *ast.Heading, source []byte) {
	r.document.Title = markdownHeadingTitle(heading, source)
}

func markdownHeadingTitle(heading *ast.Heading, source []byte) string {
	var title strings.Builder
	for _, segment := range heading.Source() {
		title.Write(segment.Bytes(source))
	}
	return strings.TrimSpace(title.String())
}

func (r *MarkdownReader) readSection(heading *ast.Heading, source []byte, content string) {
	r.document.Sections = append(r.document.Sections, MarkdownDocumentSection{
		Title:   markdownHeadingTitle(heading, source),
		Content: strings.Trim(content, "\r\n"),
	})
}
