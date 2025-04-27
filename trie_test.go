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
