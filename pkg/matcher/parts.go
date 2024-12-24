package matcher

import (
	"bytes"
	"fmt"
)

// Match any number of elements
type PathParts struct {
	Seperator byte
	Min, Max  int
}

func (m *PathParts) Match(text []byte) []Result {
	matches := make([]Result, 0, 10)
	pos := 0

	matches = append(matches, Result{Tail: text, Matcher: m})
	for {
		if len(matches) > m.Max {
			break
		}
		slashPos := bytes.IndexByte(text[pos:], m.Seperator)
		if slashPos == -1 {
			matches = append(matches, Result{Fragment: text, Matcher: m})
			break
		} else {
			pos += slashPos
			matches = append(matches, Result{Fragment: text[:pos], Tail: text[pos:], Matcher: m})
			pos += 1
		}
	}
	if len(matches) < m.Min {
		return nil
	}
	return matches[m.Min:]
}

func (m *PathParts) Split(other Interface) (common, left, right Interface) {
	if otherParts, ok := other.(*PathParts); ok {
		if m.Seperator == otherParts.Seperator && m.Min == otherParts.Min && m.Max == otherParts.Max {
			return m, nil, nil
		}
	}
	return nil, m, other
}

func (m *PathParts) String() string {
	return fmt.Sprintf("path parts %q: %d..%d", m.Seperator, m.Min, m.Max)
}
