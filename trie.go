package etrie

import (
	"errors"
	"strings"
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
		root:     &Node[T]{},
		splitter: splitter,
	}
}

func (t *Trie[T]) Insert(pattern string, value T) error {
	parts := t.splitter.Split(pattern)
	return t.insert(t.root, pattern, parts, value)
}

func (t *Trie[T]) MustInsert(pattern string, value T) {
	err := t.Insert(pattern, value)
	if err != nil {
		panic(err)
	}
}

func (t *Trie[T]) insert(n *Node[T], pattern string, parts []PatternPart, value T) error {
	if n.isEmptyNode() {
		part := parts[0]
		n.patternPart = part
		if len(parts) == 1 {
			n.value = &value
			return nil
		}
		return t.insertChild(n, pattern, parts[1:], value)
	}
	return errors.New("not implemented")
}

func (t *Trie[T]) insertChild(n *Node[T], pattern string, parts []PatternPart, value T) error {
	part := parts[0]
	if part.Parameter {
		return t.insertParameterChild(n, pattern, part, parts[1:], value)
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
		return t.insertChild(childNode, pattern, parts[1:], value)
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
	return t.insertChild(childNode, pattern, parts[1:], value)
}

func (t *Trie[T]) insertParameterChild(n *Node[T], pattern string, part PatternPart, parts []PatternPart, value T) error {
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
	return t.insertChild(childNode, pattern, parts[1:], value)
}

func (t *Trie[T]) Search(path string, params *[]Param) *Node[T] {
	return t.search(t.root, path, params)
}

func (t *Trie[T]) search(n *Node[T], path string, params *[]Param) (resultNode *Node[T]) {
	originParamsLen := len(*params)
	defer func() {
		if resultNode == nil {
			*params = (*params)[:originParamsLen]
		}
	}()

	if n.patternPart.Parameter {
		param := t.splitter.ConsumeParameter(path, n.patternPart)
		*params = append(*params, param)
		path = path[len(param.Value):]
	} else {
		if !strings.HasPrefix(path, n.patternPart.Value) {
			return nil
		}
		path = path[len(n.patternPart.Value):]
	}
	if path == "" {
		if n.value != nil {
			return n
		}
		return nil
	}
	childNode := n.children[path[0]]
	if childNode != nil {
		resultNode = t.search(childNode, path, params)
		if resultNode != nil {
			return resultNode
		}
	}
	if n.parameterChild != nil {
		resultNode = t.search(n.parameterChild, path, params)
		if resultNode != nil {
			return resultNode
		}
	}
	return nil
}

type Node[T any] struct {
	children       map[byte]*Node[T]
	parameterChild *Node[T]
	patternPart    PatternPart
	value          *T
	pattern        string
}

func (n *Node[T]) GetValue() T {
	return *n.value
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
	ConsumeParameter(path string, part PatternPart) Param
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

func (s TextSplitter) ConsumeParameter(path string, part PatternPart) Param {
	return Param{}
}

type Param struct {
	Key   string
	Value string
}
