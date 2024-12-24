package matcher

import (
	"bytes"
	"fmt"
)

type Text struct {
	Pattern   []byte
	Seperator byte
}

func (m *Text) Match(text []byte) []Result {
	switch m.Pattern[0] {
	case '*':
		return m.matchAny(m.Pattern[1:], text)
	case '?':
		return m.matchOne(m.Pattern[1:], text)
	default:
		return m.matchExact(text)
	}
}

func (m *Text) matchExact(text []byte) []Result {
	if len(text) < len(m.Pattern) {
		return []Result{}
	}
	patternLen := len(m.Pattern)
	if bytes.Equal(text[:patternLen], m.Pattern) {
		return []Result{{Fragment: text[:patternLen], Tail: text[patternLen:], Matcher: m}}
	}
	return []Result{}
}

func (m *Text) matchOne(pattern, text []byte) []Result {
	if len(text) < 1+len(pattern) {
		return []Result{}
	}
	patternLenPlus1 := len(pattern) + 1
	if bytes.Equal(pattern, text[1:patternLenPlus1]) {
		return []Result{{Fragment: text[:patternLenPlus1], Tail: text[patternLenPlus1:], Matcher: m}}
	}
	return []Result{}
}

func (m *Text) matchAny(pattern, text []byte) []Result {
	if len(text) < len(m.Pattern) {
		return []Result{}
	}
	seq := make([]Result, 0, 4)
	base := 0
	for base < len(text) {
		matchPos := bytes.Index(text[base:], pattern)
		if matchPos == -1 || bytes.IndexByte(text[base:base+matchPos], m.Seperator) != -1 {
			break
		}
		seq = append(seq, Result{Fragment: text[:base+matchPos], Tail: text[base+matchPos:], Matcher: m})
		base += matchPos + 1
	}
	return seq
}

func (m *Text) Split(other Interface) (common, left, right Interface) {
	if otherText, ok := other.(*Text); ok {
		if bytes.Equal(m.Pattern, otherText.Pattern) {
			return m, nil, nil
		}
		//find shortest common prefix
		commonText := make([]byte, 0, len(m.Pattern))
		for i := 0; i < len(m.Pattern) && i < len(otherText.Pattern); i++ {
			if m.Pattern[i] == otherText.Pattern[i] {
				commonText = append(commonText, m.Pattern[i])
			} else {
				break
			}
		}
		if len(commonText) > 0 {
			common = &Text{Pattern: commonText}
			remainder := len(m.Pattern) - len(commonText)
			if remainder > 0 {
				left = &Text{Pattern: m.Pattern[len(commonText):]}
			}
			remainder = len(otherText.Pattern) - len(commonText)
			if remainder > 0 {
				right = &Text{Pattern: otherText.Pattern[len(commonText):]}
			}
			return
		}
	}
	return nil, m, other
}

func (m *Text) String() string {
	return fmt.Sprintf("text %q", m.Pattern)
}
