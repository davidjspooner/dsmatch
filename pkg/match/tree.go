package match

type Leaf[T any] struct {
	Result  Result
	Payload T
}

type Branch[T any] struct {
	children map[Interface]*Branch[T]
	Payload  T
}

func newBranch[T any]() *Branch[T] {
	return &Branch[T]{children: make(map[Interface]*Branch[T])}
}

func (b *Branch[T]) findLeaves(text []byte) []Leaf[T] {
	leaves := make([]Leaf[T], 0, 1)
	for matcher, child := range b.children {
		results := matcher.Match(text)
		for _, result := range results {
			if len(result.Tail) == 0 {
				leaves = append(leaves, Leaf[T]{Result: result, Payload: child.Payload})
			} else {
				leaves = append(leaves, child.findLeaves(result.Tail)...)
			}
		}
	}
	return leaves
}

func (b *Branch[T]) Add(matcher Interface) (*Branch[T], error) {
	for m, child := range b.children {
		common, left, right := m.Split(matcher)
		if common != nil {
			if left == nil && right == nil {
				return child, nil
			}
			delete(b.children, m)
			newChild := newBranch[T]()
			if left != nil {
				newChild.children[left] = child
			}
			if right != nil {
				child = newBranch[T]()
				newChild.children[right] = child
			}
			b.children[common] = newChild
			return child, nil
		}
	}
	child := newBranch[T]()
	b.children[matcher] = child
	return child, nil
}

type Tree[T any] struct {
	root *Branch[T]
}

func (t *Tree[T]) Empty() bool {
	return t.root == nil
}

func (t *Tree[T]) FindLeaves(text []byte) []Leaf[T] {
	return t.root.findLeaves(text)
}

func (t *Tree[T]) Add(matchers Sequence) (*Branch[T], error) {
	if t.root == nil {
		t.root = newBranch[T]()
	}
	current := t.root
	var err error
	for _, matcher := range matchers {
		current, err = current.Add(matcher)
		if err != nil {
			return nil, err
		}
	}
	return current, nil
}
