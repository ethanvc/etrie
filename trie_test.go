package etrie

import (
	"testing"
)

func TestEmptyPath(t *testing.T) {
	trie := NewTrie[int](nil)
	trie.MustInsert("", 3)
}
