package match

type EndOfText struct {
}

func (m *EndOfText) Match(text []byte) []Result {
	if len(text) == 0 {
		return []Result{{Matcher: m}}
	}
	return []Result{}
}

func (m *EndOfText) Split(other Interface) (common, left, right Interface) {
	if _, ok := other.(*EndOfText); ok {
		return m, nil, nil
	}
	return nil, m, other
}
