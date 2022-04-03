// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comment

import (
	"bytes"
	"fmt"
	"strings"
)

// A Printer is a doc comment printer.
// The fields in the struct can be filled in before calling
// any of the printing methods
// in order to customize the details of the printing process.
type Printer struct {
	// HeadingLevel is the nesting level used for
	// HTML and Markdown headings.
	// If HeadingLevel is zero, it defaults to level 3,
	// meaning to use <h3> and ###.
	HeadingLevel int

	// HeadingID is a function that computes the heading ID
	// (anchor tag) to use for the heading h when generating
	// HTML and Markdown. If HeadingID returns an empty string,
	// then the heading ID is omitted.
	// If HeadingID is nil, h.DefaultID is used.
	HeadingID func(h *Heading) string

	// DocLinkURL is a function that computes the URL for the given DocLink.
	// If DocLinkURL is nil, then link.DefaultURL(p.DocLinkBaseURL) is used.
	DocLinkURL func(link *DocLink) string

	// DocLinkBaseURL is used when DocLinkURL is nil,
	// passed to [DocLink.DefaultURL] to construct a DocLink's URL.
	// See that method's documentation for details.
	DocLinkBaseURL string

	// TextPrefix is a prefix to print at the start of every line
	// when generating text output using the Text method.
	TextPrefix string

	// TextCodePrefix is the prefix to print at the start of each
	// preformatted (code block) line when generating text output,
	// instead of (not in addition to) TextPrefix.
	// If TextCodePrefix is the empty string, it defaults to TextPrefix+"\t".
	TextCodePrefix string

	// TextWidth is the maximum width text line to generate,
	// measured in Unicode code points,
	// excluding TextPrefix and the newline character.
	// If TextWidth is zero, it defaults to 80 minus the number of code points in TextPrefix.
	// If TextWidth is negative, there is no limit.
	TextWidth int
}

type commentPrinter struct {
	*Printer
	headingPrefix string
	needDoc       map[string]bool
}

// Comment returns the standard Go formatting of the Doc,
// without any comment markers.
func (p *Printer) Comment(d *Doc) []byte {
	cp := &commentPrinter{Printer: p}
	var out bytes.Buffer
	for i, x := range d.Content {
		if i > 0 && blankBefore(x) {
			out.WriteString("\n")
		}
		cp.block(&out, x)
	}

	return out.Bytes()
}

// blankBefore reports whether the block x requires a blank line before it.
// All blocks do, except for Lists that return false from x.BlankBefore().
func blankBefore(x Block) bool {
	if x, ok := x.(*List); ok {
		return x.BlankBefore()
	}
	return true
}

// block prints the block x to out.
func (p *commentPrinter) block(out *bytes.Buffer, x Block) {
	switch x := x.(type) {
	default:
		fmt.Fprintf(out, "?%T", x)

	case *Paragraph:
		p.text(out, "", x.Text)
		out.WriteString("\n")
	}
}

// text prints the text sequence x to out.
func (p *commentPrinter) text(out *bytes.Buffer, indent string, x []Text) {
	for _, t := range x {
		switch t := t.(type) {
		case Plain:
			p.indent(out, indent, string(t))
		}
	}
}

// indent prints s to out, indenting with the indent string
// after each newline in s.
func (p *commentPrinter) indent(out *bytes.Buffer, indent, s string) {
	for s != "" {
		line, rest, ok := strings.Cut(s, "\n")
		out.WriteString(line)
		if ok {
			out.WriteString("\n")
			out.WriteString(indent)
		}
		s = rest
	}
}