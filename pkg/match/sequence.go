package match

type Sequence []Interface

//var _ Interface = Sequence{}

func (m Sequence) Match(text []byte) []Result {
	if len(m) == 0 {
		return []Result{}
	}
	bases := m[0].Match(text)
	remainingSequence := m[1:]
	results := []Result{}
	for _, base := range bases {
		interim := remainingSequence.Match(base.Tail)
		for i := range interim {
			newKeyValues := make([]KeyValue, len(base.KeyValues)+len(interim[i].KeyValues))
			copy(newKeyValues, base.KeyValues)
			copy(newKeyValues[len(base.KeyValues):], interim[i].KeyValues)
			interim[i].KeyValues = newKeyValues
			bytes := make([]byte, len(base.Fragment)+len(interim[i].Fragment))
			bytes = append(bytes, base.Fragment...)
			bytes = append(bytes, interim[i].Fragment...)
			interim[i].Fragment = bytes
		}
		results = append(results, interim...)
	}
	return results
}
