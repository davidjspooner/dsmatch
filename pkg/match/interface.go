package match

type KeyValue struct {
	Key, Value string
}

type Result struct {
	Fragment, Tail []byte
	KeyValues      []KeyValue
	Matcher        Interface
}

type Interface interface {
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
