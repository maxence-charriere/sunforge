package vi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
