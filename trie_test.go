package etrie

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEmptyPath(t *testing.T) {
	var resultNode *Node[int]
	trie := NewTrie[int](nil)
	trie.MustInsert("", 3)
	var params []Param
	resultNode = trie.Search("", &params)
	require.Equal(t, 3, resultNode.GetValue())
	require.Equal(t, "", resultNode.GetPattern())
}

func TestStaticPath(t *testing.T) {
	trie := NewTrie[int](nil)
	trie.MustInsert("/abc/bcd", 1)
	n := trie.Search("/abc/bcd", &[]Param{})
	require.Equal(t, 1, n.GetValue())
	require.Equal(t, "/abc/bcd", n.GetPattern())
}

func TestStaticPath2(t *testing.T) {
	trie := NewTrie[int](nil)
	trie.MustInsert("/abc/bcd", 1)
	trie.MustInsert("/abcd/bcd", 2)
	n := trie.Search("/abcd/bcd", &[]Param{})
	require.Equal(t, 2, n.GetValue())
	require.Equal(t, "/abcd/bcd", n.GetPattern())
}

func TestSplitPath(t *testing.T) {
	s := GinPathSplitter{}
	parts, err := s.Split("")
	require.Error(t, err)
	parts, _ = s.Split("/")
	equalParts(t, parts, "/")
}

func equalParts(t *testing.T, parts []PatternPart, expected ...string) {
	if len(parts) != len(expected) {
		t.Errorf("expecting %d parts, got %d", len(expected), len(parts))
		t.FailNow()
		return
	}
	for i, part := range parts {
		if part.Value == expected[i] {
			continue
		}
		t.Errorf("got %s, expect %s", part.Value, expected[i])
		t.FailNow()
		return
	}
}
