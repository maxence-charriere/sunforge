package vi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMarkdownReader(t *testing.T) {
	tests := []struct {
		scenario string
		in       string
		expected MarkdownDocument
	}{
		{
			scenario: "empty input",
		},
		{
			scenario: "title from last H1",
			in:       "## Subtitle\n\n# First title\n\n# Second title",
			expected: MarkdownDocument{Title: "Second title"},
		},
		{
			scenario: "H1 with tabs closing hashes and CRLF",
			in:       "  #\t你好 世界 ###\r\n",
			expected: MarkdownDocument{Title: "你好 世界"},
		},
		{
			scenario: "setext headings are not read as H1",
			in:       "A title\non two lines\n======\n",
		},
		{
			scenario: "H1 inline formatting",
			in:       "# A **bold** *title* with `code` and [a link](https://example.com)",
			expected: MarkdownDocument{Title: "A **bold** *title* with `code` and [a link](https://example.com)"},
		},
		{
			scenario: "H1 entities escapes images and autolinks",
			in:       "# A &amp; B \\*literal\\* ![sun](/sun.png) <https://example.com>",
			expected: MarkdownDocument{Title: "A &amp; B \\*literal\\* ![sun](/sun.png) <https://example.com>"},
		},
		{
			scenario: "later metadata title wins",
			in:       "# Heading\n\n<!-- vi:title Metadata title -->",
			expected: MarkdownDocument{Title: "Metadata title"},
		},
		{
			scenario: "later H1 title wins",
			in:       "<!-- vi:title Metadata title -->\n\n# Heading",
			expected: MarkdownDocument{Title: "Heading"},
		},
		{
			scenario: "hash attached to title is preserved",
			in:       "# Learn C#",
			expected: MarkdownDocument{Title: "Learn C#"},
		},
		{
			scenario: "no H1",
			in:       "#hashtag\n\n## Subtitle\n\nSection\n-------",
		},
		{
			scenario: "empty H1 clears title",
			in:       "# Earlier title\n\n# ###",
		},
		{
			scenario: "bare hashes and other headings preserve title",
			in:       "# Earlier title\n\n#\n# \n#\t\n#hashtag\n## Subtitle\n##\tSubtitle",
			expected: MarkdownDocument{Title: "Earlier title"},
		},
		{
			scenario: "description comment",
			in:       "<!-- vi:description TBH -->",
			expected: MarkdownDocument{Description: "TBH"},
		},
		{
			scenario: "author with whitespace and unicode",
			in:       "<!-- vi:author\t  Renée Chen 陳 \t-->",
			expected: MarkdownDocument{Author: "Renée Chen 陳"},
		},
		{
			scenario: "later author wins",
			in:       "<!-- vi:author First Author -->\n<!-- vi:author Second Author -->",
			expected: MarkdownDocument{Author: "Second Author"},
		},
		{
			scenario: "all fields mixed with markdown",
			in: "<!-- vi:title A title -->\n# Heading\n\n" +
				"<!-- vi:author Jane Doe -->\n" +
				"<!-- vi:description A longer description. -->\nSome content.\n" +
				"<!-- vi:image /image.png -->\n" +
				"<!-- vi:og-image https://example.com/social.png -->\n" +
				"<!-- vi:modified-at 2026-11-09T12:30:00+08:00 -->",
			expected: MarkdownDocument{
				Title:       "Heading",
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
			expected: MarkdownDocument{Title: "A  title\twith spaces", Description: "TBH"},
		},
		{
			scenario: "preserve punctuation and unicode",
			in:       "<!-- vi:description 你好: **world** & <b>hello</b> -->",
			expected: MarkdownDocument{Description: "你好: **world** & <b>hello</b>"},
		},
		{
			scenario: "later value wins",
			in:       "<!-- vi:title First -->\n<!-- vi:title Second -->",
			expected: MarkdownDocument{Title: "Second"},
		},
		{
			scenario: "date only",
			in:       "<!-- vi:modified-at 2026-11-09 -->",
			expected: MarkdownDocument{ModifiedAt: time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)},
		},
		{
			scenario: "invalid date is ignored",
			in:       "<!-- vi:modified-at yesterday -->\n<!-- vi:title Still parsed -->",
			expected: MarkdownDocument{Title: "Still parsed"},
		},
		{
			scenario: "invalid date preserves previous value",
			in:       "<!-- vi:modified-at 2026-11-09 -->\n<!-- vi:modified-at 2026-02-30 -->",
			expected: MarkdownDocument{ModifiedAt: time.Date(2026, 11, 9, 0, 0, 0, 0, time.UTC)},
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
			var reader MarkdownReader
			reader.Read(test.in)
			meta := reader.Document()
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
