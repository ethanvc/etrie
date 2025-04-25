package etrie

import (
	"fmt"
	"slices"
	"testing"
)

var contentMap1 = map[byte]string{
	'a': "c",
	'c': "a",
	'd': "d",
}

func BenchmarkSearchEfficiency(b *testing.B) {
	benchmarkSearchEfficiency(b, contentMap1)
}

func benchmarkSearchEfficiency[K comparable, V any](b *testing.B, contentMap map[K]V) {
	var array []K
	for k := range contentMap {
		array = append(array, k)
	}
	lastItem := array[len(array)-1]
	var searchResult V
	b.Run(fmt.Sprintf("ArraySearch_%d", len(contentMap)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slices.Contains(array, lastItem)
		}
	})
	b.Run(fmt.Sprintf("MapSearch_%d", len(contentMap)), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			searchResult = contentMap[lastItem]
		}
	})
	_ = searchResult
}
