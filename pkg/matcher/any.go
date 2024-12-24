package matcher

// Match all remaining text
type Any struct {
}

func (m *Any) Match(text []byte) []Result {
	return []Result{{Fragment: text, Matcher: m}}
}
func (m *Any) Split(other Interface) (common, left, right Interface) {
	if _, ok := other.(*Any); ok {
		return m, nil, nil
	}
	return nil, m, other
}
func (m *Any) String() string {
	return "any"
}
