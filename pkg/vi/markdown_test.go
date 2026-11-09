package vi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseMarkdownMeta(t *testing.T) {
	tests := []struct {
		scenario string
		in       string
		expected MarkdownMeta
	}{
		{
			scenario: "empty input",
		},
		{
			scenario: "title from first H1",
			in:       "## Subtitle\n\n# First title\n\n# Second title",
			expected: MarkdownMeta{Title: "First title"},
		},
		{
			scenario: "H1 with tabs closing hashes and CRLF",
			in:       "  #\t你好 世界 ###\r\n",
			expected: MarkdownMeta{Title: "你好 世界"},
		},
		{
			scenario: "setext H1",
			in:       "A title\non two lines\n======\n",
			expected: MarkdownMeta{Title: "A title on two lines"},
		},
		{
			scenario: "H1 inline formatting",
			in:       "# A **bold** *title* with `code` and [a link](https://example.com)",
			expected: MarkdownMeta{Title: "A bold title with code and a link"},
		},
		{
			scenario: "H1 entities escapes images and autolinks",
			in:       "# A &amp; B \\*literal\\* ![sun](/sun.png) <https://example.com>",
			expected: MarkdownMeta{Title: "A & B *literal* sun https://example.com"},
		},
		{
			scenario: "explicit title after H1 takes precedence",
			in:       "# Heading\n\n<!-- vi:title Metadata title -->",
			expected: MarkdownMeta{Title: "Metadata title"},
		},
		{
			scenario: "ignore headings in code and HTML comments",
			in:       "```md\n# Fenced example\n```\n\n    # Indented example\n\n<!--\n# Comment\n-->\n\n# Real title",
			expected: MarkdownMeta{Title: "Real title"},
		},
		{
			scenario: "no H1",
			in:       "#hashtag\n\n## Subtitle\n\nSection\n-------",
		},
		{
			scenario: "empty first H1",
			in:       "#\n\n# Later title",
		},
		{
			scenario: "description comment",
			in:       "<!-- vi:description TBH -->",
			expected: MarkdownMeta{Description: "TBH"},
		},
		{
			scenario: "author with whitespace and unicode",
			in:       "<!-- vi:author\t  Renée Chen 陳 \t-->",
			expected: MarkdownMeta{Author: "Renée Chen 陳"},
		},
		{
			scenario: "later author wins",
			in:       "<!-- vi:author First Author -->\n<!-- vi:author Second Author -->",
			expected: MarkdownMeta{Author: "Second Author"},
		},
		{
			scenario: "all fields mixed with markdown",
			in: "<!-- vi:title A title -->\n# Heading\n\n" +
				"<!-- vi:author Jane Doe -->\n" +
				"<!-- vi:description A longer description. -->\nSome content.\n" +
				"<!-- vi:image /image.png -->\n" +
				"<!-- vi:og-image https://example.com/social.png -->\n" +
				"<!-- vi:modified-at 2026-11-09T12:30:00+08:00 -->",
			expected: MarkdownMeta{
				Title:       "A title",
				Author:      "Jane Doe",
				Description: "A longer description.",
				Image:       "/image.png",
				OGImage:     "https://example.com/social.png",
				ModifiedAt:  time.Date(2026, 11, 9, 12, 30, 0, 0, time.FixedZone("", 8*60*60)),
			},
		},
		{
			scenario: "tabs padding and CRLF",
			in:       " \t<!--\tvi:title\t  A  title\twith spaces \t--> \r\n<!--vi:description TBH-->",
			expected: MarkdownMeta{Title: "A  title\twith spaces", Description: "TBH"},
		},
		{
			scenario: "preserve punctuation and unicode",
			in:       "<!-- vi:description 你好: **world** & <b>hello</b> -->",
			expected: MarkdownMeta{Description: "你好: **world** & <b>hello</b>"},
		},
		{
			scenario: "later value wins",
			in:       "<!-- vi:title First -->\n<!-- vi:title Second -->",
			expected: MarkdownMeta{Title: "Second"},
		},
		{
			scenario: "date only",
			in:       "<!-- vi:modified-at 2026-11-09 -->",
			expected: MarkdownMeta{ModifiedAt: time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)},
		},
		{
			scenario: "invalid date is ignored",
			in:       "<!-- vi:modified-at yesterday -->\n<!-- vi:title Still parsed -->",
			expected: MarkdownMeta{Title: "Still parsed"},
		},
		{
			scenario: "invalid date preserves previous value",
			in:       "<!-- vi:modified-at 2026-11-09 -->\n<!-- vi:modified-at 2026-02-30 -->",
			expected: MarkdownMeta{ModifiedAt: time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)},
		},
		{
			scenario: "unknown keys and ordinary comments",
			in:       "<!-- A normal comment -->\n<!-- vi:unknown ignored -->\n<!-- title ignored -->",
		},
		{
			scenario: "missing values",
			in:       "<!-- -->\n<!-- vi:title -->\n<!--vi:description-->",
		},
		{
			scenario: "malformed and inline comments",
			in: "vi:title No comment\n<!-- vi:title Unclosed\nvi:title No opening -->\n" +
				"Text <!-- vi:title Inline -->\n<!-- vi:title Trailing --> text\n" +
				"<!-- vi:title First --><!-- vi:title Second -->",
		},
		{
			scenario: "multiline comments",
			in:       "<!--\nvi:title Multiline\n-->\n<!-- vi:description Multiple\nlines -->",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			meta := ParseMarkdownMeta(test.in)
			// Compare timestamps independently of the machine's local timezone.
			meta.ModifiedAt = meta.ModifiedAt.UTC()
			test.expected.ModifiedAt = test.expected.ModifiedAt.UTC()
			require.Equal(t, test.expected, meta)
		})
	}
}

func TestParseMarkdown(t *testing.T) {
	tests := []struct {
		scenario string
		in       string
		expected string
	}{
		{
			scenario: "empty input",
		},
		{
			scenario: "whitespace only",
			in:       " \t\n\n",
		},
		{
			scenario: "paragraphs",
			in:       "First paragraph.\n\nSecond paragraph.",
			expected: "<p>First paragraph.</p>\n<p>Second paragraph.</p>\n",
		},
		{
			scenario: "headings",
			in:       "# Title\n\n## Subtitle\n\nSection\n-------",
			expected: "<h1>Title</h1>\n<h2>Subtitle</h2>\n<h2>Section</h2>\n",
		},
		{
			scenario: "inline formatting",
			in:       "An *emphasis*, **strong text**, and `code`.",
			expected: "<p>An <em>emphasis</em>, <strong>strong text</strong>, and <code>code</code>.</p>\n",
		},
		{
			scenario: "unordered list",
			in:       "- First\n- Second",
			expected: "<ul>\n<li>First</li>\n<li>Second</li>\n</ul>\n",
		},
		{
			scenario: "ordered list with starting number",
			in:       "3. Third\n4. Fourth",
			expected: "<ol start=\"3\">\n<li>Third</li>\n<li>Fourth</li>\n</ol>\n",
		},
		{
			scenario: "blockquote",
			in:       "> A **quoted** paragraph.",
			expected: "<blockquote>\n<p>A <strong>quoted</strong> paragraph.</p>\n</blockquote>\n",
		},
		{
			scenario: "link with title",
			in:       "[Example](https://example.com \"Example site\")",
			expected: "<p><a href=\"https://example.com\" title=\"Example site\">Example</a></p>\n",
		},
		{
			scenario: "reference link",
			in:       "[Example][site]\n\n[site]: https://example.com",
			expected: "<p><a href=\"https://example.com\">Example</a></p>\n",
		},
		{
			scenario: "image",
			in:       "![A sunset](/sunset.png \"Sunset\")",
			expected: "<p><img src=\"/sunset.png\" alt=\"A sunset\" title=\"Sunset\"></p>\n",
		},
		{
			scenario: "fenced code with language and escaped HTML",
			in:       "```html\n<div>A & B</div>\n```",
			expected: "<pre><code class=\"language-html\">&lt;div&gt;A &amp; B&lt;/div&gt;\n</code></pre>\n",
		},
		{
			scenario: "indented code",
			in:       "    hello()\n",
			expected: "<pre><code>hello()\n</code></pre>\n",
		},
		{
			scenario: "soft and hard line breaks",
			in:       "First\nsecond  \nthird",
			expected: "<p>First\nsecond<br>\nthird</p>\n",
		},
		{
			scenario: "thematic break",
			in:       "---",
			expected: "<hr>\n",
		},
		{
			scenario: "escaped punctuation and HTML characters",
			in:       "\\*literal\\* & 1 < 2",
			expected: "<p>*literal* &amp; 1 &lt; 2</p>\n",
		},
		{
			scenario: "unicode text",
			in:       "你好 **世界** 🌅",
			expected: "<p>你好 <strong>世界</strong> 🌅</p>\n",
		},
		{
			scenario: "unclosed inline markup",
			in:       "An *unfinished emphasis and [link",
			expected: "<p>An *unfinished emphasis and [link</p>\n",
		},
		{
			scenario: "raw HTML is omitted",
			in:       "<script>alert('hello')</script>",
			expected: "<!-- raw HTML omitted -->\n",
		},
		{
			scenario: "javascript link destination is omitted",
			in:       "[Click](javascript:alert%281%29)",
			expected: "<p><a href=\"\">Click</a></p>\n",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			md, err := ParseMarkdown(test.in)
			require.NoError(t, err)
			require.Equal(t, test.expected, md)
			t.Log(md)
		})
	}
}
