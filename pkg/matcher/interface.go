package matcher

import "fmt"

type KeyValue struct {
	Key, Value string
}

type Result struct {
	Fragment, Tail []byte
	KeyValues      []KeyValue
	Matcher        Interface
}

func CombineResults(results ...*Result) Result {
	combined :=Result{
		Fragment:  make([]byte, 0, 1024),
		KeyValues: make([]KeyValue, 0, 10),
	}
	for _, result := range results {
		combined.Fragment = append(combined.Fragment, result.Fragment...)
		combined.KeyValues = append(combined.KeyValues, result.KeyValues...)
	}
	if len(results) > 0 {
		combined.Tail = results[len(results)-1].Tail
	}
	return combined
}

type Interface interface {
	fmt.Stringer
	Match(text []byte) []Result
	Split(other Interface) (common, left, right Interface)
}

func Equal(a, b Interface) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	common, left, right := a.Split(b)
	if common == nil || left != nil || right != nil {
		return false
	}
	return true
}
