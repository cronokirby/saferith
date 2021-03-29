# a24e618bccbc040c7121394c889e8bdd0dca2d01 (2021-03-29)

Implement free limb injection

```
[safenum] → go test -bench=.                                                                                                                                                     
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        6953644               212.4 ns/op
BenchmarkModAddBig-4                     1100115              1040 ns/op
BenchmarkLargeModAddBig-4                 585422              1841 ns/op
BenchmarkMulBig-4                         429370              3033 ns/op
BenchmarkModMulBig-4                      339429              3708 ns/op
BenchmarkLargeModMulBig-4                 284680              4709 ns/op
BenchmarkModBig-4                        1257049               936.7 ns/op
BenchmarkLargeModBig-4                    650443              1853 ns/op
BenchmarkModInverseBig-4                  823682              1468 ns/op
BenchmarkLargeModInverseBig-4              98150             12050 ns/op
BenchmarkExpBig-4                           8720            135014 ns/op
BenchmarkLargeExpBig-4                        43          25975861 ns/op
BenchmarkSetBytesBig-4                   5997328               225.9 ns/op
BenchmarkAddNat-4                        5293556               211.6 ns/op
BenchmarkModAddNat-4                       55104             21435 ns/op
BenchmarkLargeModAddNat-4                 108148             10385 ns/op
BenchmarkMulNat-4                         140707              8281 ns/op
BenchmarkModMulNat-4                       24109             51049 ns/op
BenchmarkLargeModMulNat-4                  17731             63715 ns/op
BenchmarkModNat-4                          56545             20935 ns/op
BenchmarkLargeModNat-4                    114956             10484 ns/op
BenchmarkModInverseNat-4                   44448             26838 ns/op
BenchmarkLargeModInverseNat-4                271           4304645 ns/op
BenchmarkExpNat-4                            196           5948716 ns/op
BenchmarkLargeExpNat-4                         3         359713739 ns/op
BenchmarkSetBytesNat-4                    836378              1506 ns/op
PASS
ok      github.com/cronokirby/safenum   44.859s
[safenum] → go test -bench=. -tags math_big_pure_go                                                                                                                              
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        5669407               252.4 ns/op
BenchmarkModAddBig-4                      946286              1171 ns/op
BenchmarkLargeModAddBig-4                 215196              5212 ns/op
BenchmarkMulBig-4                         222129              4557 ns/op
BenchmarkModMulBig-4                      215752              5561 ns/op
BenchmarkLargeModMulBig-4                 125840              9571 ns/op
BenchmarkModBig-4                        1275968              1004 ns/op
BenchmarkLargeModBig-4                    234193              5070 ns/op
BenchmarkModInverseBig-4                  877096              1414 ns/op
BenchmarkLargeModInverseBig-4              53406             23734 ns/op
BenchmarkExpBig-4                           8257            141085 ns/op
BenchmarkLargeExpBig-4                        21          51494745 ns/op
BenchmarkSetBytesBig-4                   4290268               292.7 ns/op
BenchmarkAddNat-4                        5286901               272.6 ns/op
BenchmarkModAddNat-4                       56404             21761 ns/op
BenchmarkLargeModAddNat-4                 109947             11317 ns/op
BenchmarkMulNat-4                          77644             15774 ns/op
BenchmarkModMulNat-4                       21037             56950 ns/op
BenchmarkLargeModMulNat-4                  16245             73966 ns/op
BenchmarkModNat-4                          57388             20625 ns/op
BenchmarkLargeModNat-4                    111604             10840 ns/op
BenchmarkModInverseNat-4                   45824             25909 ns/op
BenchmarkLargeModInverseNat-4                181           6707925 ns/op
BenchmarkExpNat-4                            192           6192261 ns/op
BenchmarkLargeExpNat-4                         3         415365259 ns/op
BenchmarkSetBytesNat-4                    613538              1669 ns/op
PASS
ok      github.com/cronokirby/safenum   39.011s
```

# 020e34e3436d885500318e1777e7364a1c3c393d (2021-03-27)

Limb by limb reduction

``` 
[safenum] → go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        7193640               174.0 ns/op
BenchmarkModAddBig-4                     1129430              1049 ns/op
BenchmarkLargeModAddBig-4                 637372              1735 ns/op
BenchmarkMulBig-4                         470226              2635 ns/op
BenchmarkModMulBig-4                      313644              3456 ns/op
BenchmarkLargeModMulBig-4                 257732              4091 ns/op
BenchmarkModBig-4                        1342916               891.9 ns/op
BenchmarkLargeModBig-4                    669978              1662 ns/op
BenchmarkModInverseBig-4                  772887              1358 ns/op
BenchmarkLargeModInverseBig-4             102363             11236 ns/op
BenchmarkExpBig-4                           8823            136720 ns/op
BenchmarkLargeExpBig-4                        40          25947153 ns/op
BenchmarkSetBytesBig-4                   5259864               209.5 ns/op
BenchmarkAddNat-4                        6865567               177.1 ns/op
BenchmarkModAddNat-4                       56781             21101 ns/op
BenchmarkLargeModAddNat-4                  26634             44499 ns/op
BenchmarkMulNat-4                         147716              7994 ns/op
BenchmarkModMulNat-4                       23832             49599 ns/op
BenchmarkLargeModMulNat-4                  12249             96907 ns/op
BenchmarkModNat-4                          56899             20651 ns/op
BenchmarkLargeModNat-4                     26866             45090 ns/op
BenchmarkModInverseNat-4                   44586             26605 ns/op
BenchmarkLargeModInverseNat-4                265           4357764 ns/op
BenchmarkExpNat-4                            200           5926589 ns/op
BenchmarkLargeExpNat-4                         2         622324074 ns/op
BenchmarkSetBytesNat-4                    786177              1445 ns/op
PASS
ok      github.com/cronokirby/safenum   40.069s

[safenum] → go test -bench=. -tags math_big_pure_go
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        5884750               222.7 ns/op
BenchmarkModAddBig-4                      888044              1178 ns/op
BenchmarkLargeModAddBig-4                 218853              5074 ns/op
BenchmarkMulBig-4                         256408              4895 ns/op
BenchmarkModMulBig-4                      202692              5761 ns/op
BenchmarkLargeModMulBig-4                 111362              9603 ns/op
BenchmarkModBig-4                        1289430               915.5 ns/op
BenchmarkLargeModBig-4                    238866              4949 ns/op
BenchmarkModInverseBig-4                  892677              1415 ns/op
BenchmarkLargeModInverseBig-4              60750             20164 ns/op
BenchmarkExpBig-4                           8632            140500 ns/op
BenchmarkLargeExpBig-4                        22          51087158 ns/op
BenchmarkSetBytesBig-4                   4668199               243.2 ns/op
BenchmarkAddNat-4                        4858521               255.4 ns/op
BenchmarkModAddNat-4                       56320             21361 ns/op
BenchmarkLargeModAddNat-4                  24938             48347 ns/op
BenchmarkMulNat-4                          78342             15215 ns/op
BenchmarkModMulNat-4                       20944             55907 ns/op
BenchmarkLargeModMulNat-4                  10000            111708 ns/op
BenchmarkModNat-4                          56568             20852 ns/op
BenchmarkLargeModNat-4                     23888             48265 ns/op
BenchmarkModInverseNat-4                   44446             26053 ns/op
BenchmarkLargeModInverseNat-4                177           6559094 ns/op
BenchmarkExpNat-4                            202           5873289 ns/op
BenchmarkLargeExpNat-4                         2         717422832 ns/op
BenchmarkSetBytesNat-4                    686104              1476 ns/op
PASS
ok      github.com/cronokirby/safenum   39.768s
```

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