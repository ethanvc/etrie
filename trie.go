package etrie

import (
	"errors"
	"fmt"
	path2 "path"
	"strings"
)

type Trie[T any] struct {
	root     *Node[T]
	splitter PatternSplitter
}

func NewTrie[T any](splitter PatternSplitter) *Trie[T] {
	if splitter == nil {
		splitter = GinPathSplitter{}
	}
	return &Trie[T]{
		root:     &Node[T]{},
		splitter: splitter,
	}
}

func (t *Trie[T]) Insert(pattern string, value T) error {
	parts, err := t.splitter.Split(pattern)
	if err != nil {
		return err
	}
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
			n.pattern = pattern
			n.value = &value
			return nil
		}
		return t.insertStaticChild(n, pattern, parts[1:], value)
	}
	if parts[0] == n.patternPart {
		return t.insertStaticChild(n, pattern, parts, value)
	}
	if n.patternPart.Parameter || parts[0].Parameter {
		return fmt.Errorf("both pattern use parameter but with different placeholder. the new pattern is %s, already exist placeholder is %s", pattern, n.patternPart.Value)
	}
	commonPrefix := findLongestCommonPrefix(parts[0].Value, n.patternPart.Value)
	childNode := &Node[T]{}
	*childNode = *n
	childNode.patternPart.Value = n.patternPart.Value[len(commonPrefix):]
	n.reset()
	n.children = map[byte]*Node[T]{}
	n.patternPart.Value = commonPrefix
	n.children[childNode.patternPart.Value[0]] = childNode

	childNode = &Node[T]{}
	parts[0].Value = parts[0].Value[len(commonPrefix):]
	n.children[parts[0].Value[0]] = childNode

	return t.insert(childNode, pattern, parts, value)
}

func (t *Trie[T]) insertStaticChild(n *Node[T], pattern string, parts []PatternPart, value T) error {
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
		return t.insertStaticChild(childNode, pattern, parts[1:], value)
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
	return t.insertStaticChild(childNode, pattern, parts[1:], value)
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
	return t.insertStaticChild(childNode, pattern, parts[1:], value)
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
		param, _ := t.splitter.ConsumeParameter(path, n.patternPart)
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

func (n *Node[T]) reset() {
	n.children = nil
	n.parameterChild = nil
	n.patternPart = PatternPart{}
	n.value = nil
	n.pattern = ""
}

func (n *Node[T]) GetValue() T {
	return *n.value
}

func (n *Node[T]) GetPattern() string {
	return n.pattern
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
	Split(string) ([]PatternPart, error)
	ConsumeParameter(path string, part PatternPart) (Param, error)
}

type GinPathSplitter struct {
}

func (s GinPathSplitter) Split(path string) ([]PatternPart, error) {
	path = path2.Clean(path)
	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("path must begin with /")
	}
	var parts []PatternPart
	start := 0
	for i := 0; i < len(path); i++ {
		ch := path[i]
		if ch != '/' {
			continue
		}
		if i >= len(path)-1 {
			break
		}
		i++
		if path[i] != ':' && path[i] != '*' {
			continue
		}
		paramStart := i
		for i < len(path) && path[i] != '/' {
			i++
		}
		if i-paramStart <= 0 {
			return nil, errors.New("parameter must have a valid name")
		}
		parts = append(parts, PatternPart{
			Value: path[start:paramStart],
		})
		parts = append(parts, PatternPart{
			Parameter: true,
			Value:     path[paramStart:i],
		})
		start = i
	}
	if start < len(path) {
		parts = append(parts, PatternPart{
			Value: path[start:],
		})
	}
	return parts, nil
}

func (s GinPathSplitter) ConsumeParameter(path string, part PatternPart) (Param, error) {
	if !part.Parameter {
		return Param{}, errors.New("only parameter part can call ConsumeParameter")
	}
	if part.Value[0] == '*' {
		return Param{
			Key:   part.Value[1:],
			Value: path,
		}, nil
	}
	if part.Value[0] == ':' {
		index := strings.IndexByte(path, '/')
		if index == -1 {
			index = len(path)
		}
		return Param{
			Key:   part.Value[1:],
			Value: path[:index],
		}, nil
	}
	return Param{}, errors.New("parameter pattern not support: " + part.Value)
}

type Param struct {
	Key   string
	Value string
}
