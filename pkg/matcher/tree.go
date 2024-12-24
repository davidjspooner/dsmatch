package matcher

import (
	"bytes"
	"fmt"
	"io"
)

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

func (b *Branch[T]) findLeaves(baseResult *Result, tail []byte) []Leaf[T] {
	leaves := make([]Leaf[T], 0, 1)
	for matcher, child := range b.children {
		results := matcher.Match(tail)
		for _, result := range results {
			subBaseResult := CombineResults(baseResult, &result)
			if len(result.Tail) == 0 {
				leaves = append(leaves, Leaf[T]{Result: subBaseResult, Payload: child.Payload})
			} else {
				childLeaves := child.findLeaves(&subBaseResult, result.Tail)
				leaves = append(leaves, childLeaves...)
			}
		}
	}
	return leaves
}

func (b *Branch[T]) set(matcher Interface, child *Branch[T]) {
	if b.children == nil {
		b.children = make(map[Interface]*Branch[T])
	}
	if matcher == nil {
		panic("matcher is nil")
	}
	if child == nil {
		panic("child is nil")
	}
	b.children[matcher] = child
}
func (b *Branch[T]) Add(matcher Interface) (*Branch[T], error) {
	n := 0
	for m, child := range b.children {
		n++
		fmt.Printf("n: %d, m: %v, matcher: %v\n", n, m, matcher)
		common, left, right := m.Split(matcher)
		if common != nil {
			if left == nil && right == nil {
				return child, nil
			}

			if left != nil {
				delete(b.children, m)
				newChild := newBranch[T]()
				b.set(common, newChild)
				newChild.set(left, child)
				child = newChild
			}
			if right == nil {
				return child, nil
			}
			return child.Add(right)
		}
	}
	child := newBranch[T]()
	b.set(matcher, child)
	return child, nil
}

func (b *Branch[T]) writeDebug(sb io.Writer, depth int) {
	for matcher, child := range b.children {
		for i := 0; i < depth; i++ {
			io.WriteString(sb, "  ")
		}
		io.WriteString(sb, "└ ")
		io.WriteString(sb, matcher.String())
		io.WriteString(sb, "\n")
		child.writeDebug(sb, depth+1)
	}
}

type Tree[T any] struct {
	root *Branch[T]
}

func (t *Tree[T]) Empty() bool {
	return t.root == nil
}

func (t *Tree[T]) FindLeaves(text []byte) []Leaf[T] {
	baseResult := Result{Fragment: text}
	return t.root.findLeaves(&baseResult, text)
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

func (t *Tree[T]) String() string {
	sb := bytes.Buffer{}
	t.root.writeDebug(&sb, 0)
	return sb.String()
}

func (t *Tree[T]) WriteDebug(w io.Writer, indent int) {
	t.root.writeDebug(w, indent)
}
