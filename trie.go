package etrie

import (
	"errors"
)

type Trie[T any] struct {
	root     *Node[T]
	splitter PatternSplitter
}

func NewTrie[T any](splitter PatternSplitter) *Trie[T] {
	if splitter == nil {
		splitter = TextSplitter{}
	}
	return &Trie[T]{
		splitter: splitter,
	}
}

func (t *Trie[T]) Insert(pattern string, value T) error {
	parts := t.splitter.Split(pattern)
	if t.root == nil {
		t.root = &Node[T]{}
	}
	return t.root.insert(parts, value)
}

func (t *Trie[T]) MustInsert(pattern string, value T) {
	err := t.Insert(pattern, value)
	if err != nil {
		panic(err)
	}
}

type Node[T any] struct {
	children       map[byte]*Node[T]
	parameterChild *Node[T]
	patternPart    PatternPart
	value          *T
}

func (n *Node[T]) insert(parts []PatternPart, value T) error {
	if n.isEmptyNode() {
		part := parts[0]
		n.patternPart = part
		if len(parts) == 1 {
			n.value = &value
			return nil
		}
		return n.insertChild(parts[1:], value)
	}
	return errors.New("not implemented")
}

func (n *Node[T]) insertChild(parts []PatternPart, value T) error {
	part := parts[0]
	if part.Parameter {
		return n.insertParameterChild(part, parts[1:], value)
	}
	firstByte := part.Value[0]
	childNode := n.children[firstByte]
	if childNode == nil {
		childNode = &Node[T]{
			patternPart: part,
		}
		n.children[firstByte] = childNode
		if len(parts) == 1 {
			childNode.value = &value
			return nil
		}
		return childNode.insertChild(parts[1:], value)
	}
	longestPrefix := findLongestCommonPrefix(part.Value, childNode.patternPart.Value)
	newParentNode := &Node[T]{
		patternPart: PatternPart{
			Value: longestPrefix,
		},
		children: map[byte]*Node[T]{},
	}
	n.children[firstByte] = newParentNode

	childNode.patternPart.Value = childNode.patternPart.Value[len(longestPrefix):]
	newParentNode.children[childNode.patternPart.Value[0]] = childNode

	childNode = &Node[T]{
		patternPart: PatternPart{
			Value: part.Value[len(longestPrefix):],
		},
	}
	newParentNode.children[childNode.patternPart.Value[0]] = childNode
	if len(parts) == 1 {
		childNode.value = &value
		return nil
	}
	return childNode.insertChild(parts[1:], value)
}

func (n *Node[T]) insertParameterChild(part PatternPart, parts []PatternPart, value T) error {
	if n.parameterChild != nil {
		return errors.New("conflict pattern: same place can only have one parameter part")
	}
	childNode := &Node[T]{
		patternPart: part,
	}
	n.parameterChild = childNode
	if len(parts) == 1 {
		childNode.value = &value
		return nil
	}
	return childNode.insertChild(parts[1:], value)
}

func findLongestCommonPrefix(s1, s2 string) string {
	i := 0
	for ; i < len(s1) && i < len(s2); i++ {
		if s1[i] != s2[i] {
			break
		}
	}
	return s1[:i]
}

func (n *Node[T]) isEmptyNode() bool {
	return len(n.children) == 0 && n.value == nil && n.parameterChild == nil
}

type PatternPart struct {
	Parameter bool
	Value     string
}

type PatternSplitter interface {
	Split(string) []PatternPart
	ConsumeParameter(pattern string, part PatternPart) string
}

type TextSplitter struct {
}

func (s TextSplitter) Split(text string) []PatternPart {
	return []PatternPart{
		{
			Parameter: false,
			Value:     text,
		},
	}
}

func (s TextSplitter) ConsumeParameter(pattern string, part PatternPart) string {
	return ""
}
