package etrie

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func generateMap(count int) map[byte]string {
	result := make(map[byte]string)
	for i := 0; i < count; i++ {
		result[byte(i)] = fmt.Sprintf("%d", i)
	}
	return result
}

func BenchmarkSearchEfficiency(b *testing.B) {
	benchmarkSearchEfficiency(b, generateMap(10))
	benchmarkSearchEfficiency(b, generateMap(20))
	benchmarkSearchEfficiency(b, generateMap(30))
}

func benchmarkSearchEfficiency[K comparable, V any](b *testing.B, contentMap map[K]V) {
	var array []K
	for k := range contentMap {
		array = append(array, k)
	}
	lastItem := array[len(array)-1]
	var searchResult V
	b.Run(fmt.Sprintf("ArraySearch_%d_%s", len(contentMap), reflect.TypeOf(array[0]).String()),
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				slices.Contains(array, lastItem)
			}
		})
	b.Run(fmt.Sprintf("MapSearch_%d_%s", len(contentMap), reflect.TypeOf(array[0]).String()),
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				searchResult = contentMap[lastItem]
			}
		})
	_ = searchResult
}
