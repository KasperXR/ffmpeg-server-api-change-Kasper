package omniglyph

import (
	"strings"
)

type WordWrap interface {
	ReplaceGlyphs() string // Replace glyphs in the text
	Wrap() string          // Wrap the text into lines of the given width
}

type WordWrapper struct {
	Width        int    // The width of each line
	Text         string // The text to wrap
	NewLine      string // The new line character
	Prefix       string // The prefix to add to the text
	PrefixStart  bool   // Whether to add the prefix to the first line
	PrefixAll    bool   // Should the prefix be added to the text
	IndentAmount int    // Number of spaces to indent the text
	IndentStart  bool   // Should the indent be added to the first line
	IndentAll    bool   // Should the indent be added to all lines
	IndentGlyph  string // The glyph to use for indentation
	Separator    string // String for separating the words,
	Joiner       string // Joiner is the string used to join the words
}

// Escape special characters in the text
func (ww *WordWrapper) ReplaceGlyphs(glyphs, replacements []string) {
	if len(glyphs) != len(replacements) {
		panic("Glyphs and replacements must be of the same length")
	}

	for i, glyph := range glyphs {
		ww.Text = strings.Replace(ww.Text, glyph, replacements[i], -1)
	}
}

// Wrap the text into lines of the given width
// `prefixText` is the prefix to add to each line, if it's an empty string, no prefix is added
// If `indentStart` is true, the indent is added to the first line
func (ww *WordWrapper) Wrap() string {
	var result string
	var currentLine string
	var count int

	// Split the string into words, recognised and separated by spaces
	words := strings.Split(ww.Text, ww.Separator)

	// Indent the first line if required
	if ww.IndentStart {
		currentLine += strings.Repeat(ww.IndentGlyph, ww.IndentAmount)
	}

	// Prefix the text
	if ww.PrefixStart {
		currentLine += ww.Prefix
		count += len(ww.Prefix)
	}

	for idx, word := range words {

		// Add the word to the current line if it fits
		if count+len(word) < int(ww.Width) {
			currentLine += word + ww.Joiner
			count += len(word) + len(ww.Joiner)
		} else {
			// Word doesn't fit, make new line
			result += currentLine + ww.NewLine

			// We indent before prefixing the new line
			if ww.IndentAll {
				result += strings.Repeat(ww.IndentGlyph, ww.IndentAmount)
			}

			// Prefix the new line
			if ww.PrefixAll && idx != 0 {
				result += ww.Prefix
				count += len(ww.Prefix)
			}

			currentLine = word + ww.Joiner
			// Reset the count, include the word and joiner character
			count = len(word) + len(ww.Joiner)
		}

	}
	//	fmt.Println(result)

	result += currentLine

	ww.Text = result

	return ww.Text
}
