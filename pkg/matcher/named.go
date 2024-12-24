package matcher

import "fmt"

type Named struct {
	Name    string
	Matcher Interface
}

func (m *Named) Match(text []byte) []Result {
	results := m.Matcher.Match(text)
	for i := range results {
		results[i].KeyValues = append(results[i].KeyValues, KeyValue{Key: m.Name, Value: string(results[i].Fragment)})
	}
	return results
}

func (m *Named) Split(other Interface) (common, left, right Interface) {
	if otherNamed, ok := other.(*Named); ok {
		if otherNamed.Name == m.Name {
			_, subLeft, subRight := m.Matcher.Split(otherNamed.Matcher)
			if subLeft == nil && subRight == nil {
				return m, nil, nil
			}
		}
	}
	return nil, m, other
}

func (m *Named) String() string {
	return fmt.Sprintf("named %q: %s", m.Name, m.Matcher)
}
