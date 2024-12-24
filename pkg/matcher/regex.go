package matcher

import (
	"fmt"
	"regexp"
)

type Regex struct {
	Regexp *regexp.Regexp
}

func (m *Regex) Match(text []byte) []Result {
	pos := m.Regexp.FindIndex(text)
	if pos == nil {
		return nil
	}
	if pos[0] != 0 {
		return nil
	}
	match := text[0:pos[1]]
	tail := text[pos[1]:]
	return []Result{{Fragment: match, Tail: tail, Matcher: m}}
}

func (m *Regex) Split(other Interface) (common, left, right Interface) {
	if otherRegex, ok := other.(*Regex); ok {
		if m.Regexp.String() == otherRegex.Regexp.String() {
			return m, nil, nil
		}
	}
	return nil, nil, nil
}

func (m *Regex) String() string {
	return fmt.Sprintf("regex %q", m.Regexp.String())
}
