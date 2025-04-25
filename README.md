# 线性搜索和哈希表搜索的性能分析
结论：虽然在少量元素的情况下，线性搜索性能更好，但是为了避免极端case，还是采用hash表。

BenchmarkSearchEfficiency
BenchmarkSearchEfficiency/ArraySearch_10_uint8
BenchmarkSearchEfficiency/ArraySearch_10_uint8-10         	286758964	         4.062 ns/op
BenchmarkSearchEfficiency/MapSearch_10_uint8
BenchmarkSearchEfficiency/MapSearch_10_uint8-10           	165222480	         7.259 ns/op
BenchmarkSearchEfficiency/ArraySearch_50_uint8
BenchmarkSearchEfficiency/ArraySearch_50_uint8-10         	72184431	        16.56 ns/op
BenchmarkSearchEfficiency/MapSearch_50_uint8
BenchmarkSearchEfficiency/MapSearch_50_uint8-10           	100000000	        11.25 ns/op
BenchmarkSearchEfficiency/ArraySearch_100_uint8
BenchmarkSearchEfficiency/ArraySearch_100_uint8-10        	37307149	        32.23 ns/op
BenchmarkSearchEfficiency/MapSearch_100_uint8
BenchmarkSearchEfficiency/MapSearch_100_uint8-10          	142322832	         8.425 ns/op
PASS