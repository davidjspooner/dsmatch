package match

// Match all remaining text
type All struct {
}

func (m *All) Match(text []byte) []Result {
	return []Result{{Fragment: text, Matcher: m}}
}
func (m *All) Split(other Interface) (common, left, right Interface) {
	if _, ok := other.(*All); ok {
		return m, nil, nil
	}
	return nil, m, other
}
