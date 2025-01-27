package matcher

import (
	"bytes"
	"fmt"
)

type Pattern struct {
	Text      []byte
	EndOfText bool
}

func (m *Pattern) Match(text []byte) []Result {
	if len(m.Text) == 0 {
		return m.matchExact(text)
	}
	switch m.Text[0] {
	case '*':
		return m.matchAny(m.Text[1:], text)
	case '?':
		return m.matchOne(m.Text[1:], text)
	default:
		return m.matchExact(text)
	}
}

func (m *Pattern) matchExact(text []byte) []Result {
	if len(text) < len(m.Text) {
		return []Result{}
	}
	patternLen := len(m.Text)
	if bytes.Equal(text[:patternLen], m.Text) {
		if m.EndOfText && len(text) > patternLen {
			return []Result{}
		}
		return []Result{{Fragment: text[:patternLen], Tail: text[patternLen:], Matcher: m}}
	}
	return []Result{}
}

func (m *Pattern) matchOne(pattern, text []byte) []Result {
	if len(text) < 1+len(pattern) {
		return []Result{}
	}
	patternLenPlus1 := len(pattern) + 1
	if bytes.Equal(pattern, text[1:patternLenPlus1]) {
		if m.EndOfText && len(text) > patternLenPlus1 {
			return []Result{}
		}
		return []Result{{Fragment: text[:patternLenPlus1], Tail: text[patternLenPlus1:], Matcher: m}}
	}
	return []Result{}
}

func (m *Pattern) matchAny(pattern, text []byte) []Result {
	if len(text) < len(m.Text) {
		return []Result{}
	}
	seq := make([]Result, 0, 4)
	if len(m.Text) == 1 && m.EndOfText {
		seq = append(seq, Result{Fragment: text, Tail: nil, Matcher: m})
		return seq
	}
	base := 0
	for base <= len(text) {
		matchPos := bytes.Index(text[base:], pattern)
		//		if matchPos == -1 || bytes.IndexByte(text[base:base+matchPos], m.Seperator) != -1 {
		//			break
		//		}
		seq = append(seq, Result{Fragment: text[:base+matchPos], Tail: text[base+matchPos:], Matcher: m})
		base += matchPos + 1
	}
	return seq
}

func (m *Pattern) Split(other Interface) (common, left, right Interface) {
	if otherText, ok := other.(*Pattern); ok {
		if bytes.Equal(m.Text, otherText.Text) {
			return m, nil, nil
		}
		//find shortest common prefix
		commonText := make([]byte, 0, len(m.Text))
		for i := 0; i < len(m.Text) && i < len(otherText.Text); i++ {
			if m.Text[i] == otherText.Text[i] {
				commonText = append(commonText, m.Text[i])
			} else {
				break
			}
		}
		if len(commonText) > 0 {
			common = &Pattern{Text: commonText}
			remainder := len(m.Text) - len(commonText)
			if remainder > 0 {
				left = &Pattern{Text: m.Text[len(commonText):]}
			}
			remainder = len(otherText.Text) - len(commonText)
			if remainder > 0 {
				right = &Pattern{Text: otherText.Text[len(commonText):]}
			}
			return
		}
	}
	return nil, m, other
}

func (m *Pattern) String() string {
	if m.EndOfText {
		return fmt.Sprintf("text %q (EOT)", m.Text)
	}
	return fmt.Sprintf("text %q", m.Text)
}
