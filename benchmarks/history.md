# b66ac649d9ca1b1b394a7706cea6495b107dcb9c (2021-03-25)

```
[safenum] → go test -bench=.                                                                                                                                                     
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        7233490               168.4 ns/op
BenchmarkModAddBig-4                     1000000              1051 ns/op
BenchmarkLargeModAddBig-4                 599292              1772 ns/op
BenchmarkMulBig-4                         458648              2567 ns/op
BenchmarkModMulBig-4                      324234              3573 ns/op
BenchmarkLargeModMulBig-4                 285141              4197 ns/op
BenchmarkModBig-4                        1239144              1043 ns/op
BenchmarkLargeModBig-4                    711975              1732 ns/op
BenchmarkModInverseBig-4                  778612              1438 ns/op
BenchmarkLargeModInverseBig-4              91058             13763 ns/op
BenchmarkExpBig-4                           7298            139626 ns/op
BenchmarkLargeExpBig-4                        42          25930457 ns/op
BenchmarkSetBytesBig-4                   4770934               240.3 ns/op
BenchmarkAddNat-4                        6586465               166.4 ns/op
BenchmarkModAddNat-4                       27013             44123 ns/op
BenchmarkLargeModAddNat-4                   2725            438992 ns/op
BenchmarkMulNat-4                         140364              7858 ns/op
BenchmarkModMulNat-4                       12786             94944 ns/op
BenchmarkLargeModMulNat-4                   1135            893462 ns/op
BenchmarkModNat-4                          27309             41986 ns/op
BenchmarkLargeModNat-4                      2733            438626 ns/op
BenchmarkModInverseNat-4                   25102             48908 ns/op
BenchmarkLargeModInverseNat-4                237           4906628 ns/op
BenchmarkExpNat-4                             98          11851448 ns/op
BenchmarkLargeExpNat-4                         1        5541042938 ns/op
BenchmarkSetBytesNat-4                    788446              1696 ns/op
PASS
ok      github.com/cronokirby/safenum   41.362s
[safenum] → go test -bench=. -tags math_big_pure_go                                                                                                                              
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        6008720               206.7 ns/op
BenchmarkModAddBig-4                     1000000              1141 ns/op
BenchmarkLargeModAddBig-4                 226074              4899 ns/op
BenchmarkMulBig-4                         261199              5061 ns/op
BenchmarkModMulBig-4                      159643              6593 ns/op
BenchmarkLargeModMulBig-4                 119554              9764 ns/op
BenchmarkModBig-4                        1298679              1030 ns/op
BenchmarkLargeModBig-4                    226292              5268 ns/op
BenchmarkModInverseBig-4                  810531              1395 ns/op
BenchmarkLargeModInverseBig-4              58944             21170 ns/op
BenchmarkExpBig-4                           8494            136181 ns/op
BenchmarkLargeExpBig-4                        22          50598450 ns/op
BenchmarkSetBytesBig-4                   5431112               199.1 ns/op
BenchmarkAddNat-4                        6781236               177.2 ns/op
BenchmarkModAddNat-4                       33860             36012 ns/op
BenchmarkLargeModAddNat-4                   1794            664556 ns/op
BenchmarkMulNat-4                          79315             14716 ns/op
BenchmarkModMulNat-4                       14188             83209 ns/op
BenchmarkLargeModMulNat-4                    872           1320999 ns/op
BenchmarkModNat-4                          34676             35121 ns/op
BenchmarkLargeModNat-4                      1726            674276 ns/op
BenchmarkModInverseNat-4                   29527             40186 ns/op
BenchmarkLargeModInverseNat-4                164           7174169 ns/op
BenchmarkExpNat-4                            120           9723271 ns/op
BenchmarkLargeExpNat-4                         1        8705644301 ns/op
BenchmarkSetBytesNat-4                    794876              1570 ns/op
PASS
ok      github.com/cronokirby/safenum   46.227s
```