package etrie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
	var parts []PatternPart

	parts, _ = s.Split("/abc/bcd/:user")
	equalParts(t, parts, "/abc/bcd/", ":user")

	parts, _ = s.Split("/abc/bcd/:user")
	equalParts(t, parts, "/abc/bcd/", ":user")

	parts, _ = s.Split("/")
	equalParts(t, parts, "/")

	parts, _ = s.Split("/abc")
	equalParts(t, parts, "/abc")
	parts, _ = s.Split("/abc/bcd")
	equalParts(t, parts, "/abc/bcd")

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
