# 线性搜索和哈希表搜索的性能分析
结论：虽然在少量元素(20个以内)的情况下，线性搜索性能更好，但是为了避免极端case，还是采用hash表。

BenchmarkSearchEfficiency
BenchmarkSearchEfficiency/ArraySearch_10_uint8
BenchmarkSearchEfficiency/ArraySearch_10_uint8-12         	242017050	         4.884 ns/op
BenchmarkSearchEfficiency/MapSearch_10_uint8
BenchmarkSearchEfficiency/MapSearch_10_uint8-12           	137616816	         8.693 ns/op
BenchmarkSearchEfficiency/ArraySearch_20_uint8
BenchmarkSearchEfficiency/ArraySearch_20_uint8-12         	137390774	         8.711 ns/op
BenchmarkSearchEfficiency/MapSearch_20_uint8
BenchmarkSearchEfficiency/MapSearch_20_uint8-12           	138828927	         8.652 ns/op
BenchmarkSearchEfficiency/ArraySearch_30_uint8
BenchmarkSearchEfficiency/ArraySearch_30_uint8-12         	94604268	        12.48 ns/op
BenchmarkSearchEfficiency/MapSearch_30_uint8
BenchmarkSearchEfficiency/MapSearch_30_uint8-12           	138707463	         8.665 ns/op
PASS