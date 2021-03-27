## 63b23810d231b288391aee879bc456a3a17e9ef7 (2021-03-27)

Free limb injection

```
[safenum] → go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        7599657               156.0 ns/op
BenchmarkModAddBig-4                     1000000              1032 ns/op
BenchmarkLargeModAddBig-4                 691902              1747 ns/op
BenchmarkMulBig-4                         474602              2534 ns/op
BenchmarkModMulBig-4                      341878              3436 ns/op
BenchmarkLargeModMulBig-4                 276097              4125 ns/op
BenchmarkModBig-4                        1318466               915.3 ns/op
BenchmarkLargeModBig-4                    695412              1712 ns/op
BenchmarkModInverseBig-4                  760770              1375 ns/op
BenchmarkLargeModInverseBig-4             106210             11211 ns/op
BenchmarkExpBig-4                           8623            137208 ns/op
BenchmarkLargeExpBig-4                        44          25980018 ns/op
BenchmarkSetBytesBig-4                   5724906               200.9 ns/op
BenchmarkAddNat-4                        7070061               169.3 ns/op
BenchmarkModAddNat-4                       55671             21655 ns/op
BenchmarkLargeModAddNat-4                 113676             10201 ns/op
BenchmarkMulNat-4                         150108              7998 ns/op
BenchmarkModMulNat-4                       57106             21033 ns/op
BenchmarkLargeModMulNat-4                 118340             10228 ns/op
BenchmarkModNat-4                          57546             20826 ns/op
BenchmarkLargeModNat-4                    118358             10063 ns/op
BenchmarkModInverseNat-4                   44703             26640 ns/op
BenchmarkLargeModInverseNat-4                278           4313069 ns/op
BenchmarkExpNat-4                            201           5924363 ns/op
BenchmarkLargeExpNat-4                         3         339242312 ns/op
BenchmarkSetBytesNat-4                    718911              1432 ns/op
PASS
ok      github.com/cronokirby/safenum   35.747s

[safenum] → go test -bench=. -tags math_big_pure_go
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        5930488               203.0 ns/op
BenchmarkModAddBig-4                     1000000              1075 ns/op
BenchmarkLargeModAddBig-4                 239420              5165 ns/op
BenchmarkMulBig-4                         248128              4550 ns/op
BenchmarkModMulBig-4                      219744              5472 ns/op
BenchmarkLargeModMulBig-4                 124629              9412 ns/op
BenchmarkModBig-4                        1324681               892.5 ns/op
BenchmarkLargeModBig-4                    227462              4950 ns/op
BenchmarkModInverseBig-4                  811485              1329 ns/op
BenchmarkLargeModInverseBig-4              60489             19763 ns/op
BenchmarkExpBig-4                           8653            135583 ns/op
BenchmarkLargeExpBig-4                        22          50462736 ns/op
BenchmarkSetBytesBig-4                   5924011               203.6 ns/op
BenchmarkAddNat-4                        5505318               209.1 ns/op
BenchmarkModAddNat-4                       56461             21460 ns/op
BenchmarkLargeModAddNat-4                 102577             11116 ns/op
BenchmarkMulNat-4                          79017             15225 ns/op
BenchmarkModMulNat-4                       56889             21109 ns/op
BenchmarkLargeModMulNat-4                 108670             11012 ns/op
BenchmarkModNat-4                          56896             20528 ns/op
BenchmarkLargeModNat-4                    109075             10884 ns/op
BenchmarkModInverseNat-4                   45232             25925 ns/op
BenchmarkLargeModInverseNat-4                181           6408879 ns/op
BenchmarkExpNat-4                            200           5940260 ns/op
BenchmarkLargeExpNat-4                         3         401659100 ns/op
BenchmarkSetBytesNat-4                    836647              1423 ns/op
PASS
ok      github.com/cronokirby/safenum   36.844s
```

## f0bb3a76e90610a4e46c19909b6612af5825619b (2021-03-27)

Implemented limb by limb reduction, instead of bit by bit reduction

```
[safenum] → go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        5810145               189.7 ns/op
BenchmarkModAddBig-4                     1000000              1127 ns/op
BenchmarkLargeModAddBig-4                 690636              1899 ns/op
BenchmarkMulBig-4                         372877              2711 ns/op
BenchmarkModMulBig-4                      334502              3572 ns/op
BenchmarkLargeModMulBig-4                 278947              4417 ns/op
BenchmarkModBig-4                        1257882               970.3 ns/op
BenchmarkLargeModBig-4                    601494              1748 ns/op
BenchmarkModInverseBig-4                  858920              1388 ns/op
BenchmarkLargeModInverseBig-4             106639             11492 ns/op
BenchmarkExpBig-4                           8581            138093 ns/op
BenchmarkLargeExpBig-4                        43          25913686 ns/op
BenchmarkSetBytesBig-4                   5212320               241.2 ns/op
BenchmarkAddNat-4                        6347133               219.0 ns/op
BenchmarkModAddNat-4                       57331             20652 ns/op
BenchmarkLargeModAddNat-4                  27400             44292 ns/op
BenchmarkMulNat-4                         155056              8088 ns/op
BenchmarkModMulNat-4                       57728             20910 ns/op
BenchmarkLargeModMulNat-4                  27795             43191 ns/op
BenchmarkModNat-4                          59110             20022 ns/op
BenchmarkLargeModNat-4                     24699             46187 ns/op
BenchmarkModInverseNat-4                   44248             26677 ns/op
BenchmarkLargeModInverseNat-4                274           4354721 ns/op
BenchmarkExpNat-4                            206           6001053 ns/op
BenchmarkLargeExpNat-4                         2         641230924 ns/op
BenchmarkSetBytesNat-4                    673185              1510 ns/op
PASS
ok      github.com/cronokirby/safenum   37.809s

[safenum] → go test -bench=. -tags math_big_pure_go
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        6187009               195.9 ns/op
BenchmarkModAddBig-4                      973845              1099 ns/op
BenchmarkLargeModAddBig-4                 229160              5033 ns/op
BenchmarkMulBig-4                         263078              4610 ns/op
BenchmarkModMulBig-4                      216165              5469 ns/op
BenchmarkLargeModMulBig-4                 127482              9404 ns/op
BenchmarkModBig-4                        1333504               904.6 ns/op
BenchmarkLargeModBig-4                    241681              4929 ns/op
BenchmarkModInverseBig-4                  785236              1347 ns/op
BenchmarkLargeModInverseBig-4              56412             19969 ns/op
BenchmarkExpBig-4                           8462            137817 ns/op
BenchmarkLargeExpBig-4                        22          51594426 ns/op
BenchmarkSetBytesBig-4                   5712288               206.0 ns/op
BenchmarkAddNat-4                        5621718               217.7 ns/op
BenchmarkModAddNat-4                       56742             21254 ns/op
BenchmarkLargeModAddNat-4                  24981             47774 ns/op
BenchmarkMulNat-4                          79909             15119 ns/op
BenchmarkModMulNat-4                       57012             20927 ns/op
BenchmarkLargeModMulNat-4                  25014             47860 ns/op
BenchmarkModNat-4                          57055             20523 ns/op
BenchmarkLargeModNat-4                     25486             47605 ns/op
BenchmarkModInverseNat-4                   45637             25879 ns/op
BenchmarkLargeModInverseNat-4                183           6441945 ns/op
BenchmarkExpNat-4                            202           5913828 ns/op
BenchmarkLargeExpNat-4                         2         700008856 ns/op
BenchmarkSetBytesNat-4                    824199              1428 ns/op
PASS
ok      github.com/cronokirby/safenum   37.963s
```

## 553a30bd957056a42a02cd771810ea30b3c64160 (2021-03-22)

```
[safenum] → go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        7389273               159.7 ns/op
BenchmarkModAddBig-4                     1135212              1054 ns/op
BenchmarkLargeModAddBig-4                 602158              1762 ns/op
BenchmarkMulBig-4                         459238              2612 ns/op
BenchmarkModMulBig-4                      338905              3538 ns/op
BenchmarkLargeModMulBig-4                 253484              4260 ns/op
BenchmarkModBig-4                        1287739               921.3 ns/op
BenchmarkLargeModBig-4                    590497              1702 ns/op
BenchmarkModInverseBig-4                  854334              1393 ns/op
BenchmarkLargeModInverseBig-4              98396             11563 ns/op
BenchmarkExpBig-4                           8260            140352 ns/op
BenchmarkLargeExpBig-4                        40          26345668 ns/op
BenchmarkSetBytesBig-4                   4869546               207.6 ns/op
BenchmarkAddNat-4                        7881738               145.2 ns/op
BenchmarkModAddNat-4                       27153             44084 ns/op
BenchmarkLargeModAddNat-4                   2624            441829 ns/op
BenchmarkMulNat-4                         145112              8054 ns/op
BenchmarkModMulNat-4                       27182             44607 ns/op
BenchmarkLargeModMulNat-4                   2710            440824 ns/op
BenchmarkModNat-4                          26890             43682 ns/op
BenchmarkLargeModNat-4                      2320            447318 ns/op
BenchmarkModInverseNat-4                   24141             49643 ns/op
BenchmarkLargeModInverseNat-4                235           5055468 ns/op
BenchmarkExpNat-4                             92          12272048 ns/op
BenchmarkLargeExpNat-4                         1        5905129218 ns/op
BenchmarkSetBytesNat-4                    810808              1496 ns/op
PASS
ok      github.com/cronokirby/safenum   44.839s

[safenum] → go test -bench=. -tags math_big_pure_go
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                        5996150               223.6 ns/op
BenchmarkModAddBig-4                     1065055              1156 ns/op
BenchmarkLargeModAddBig-4                 235039              5100 ns/op
BenchmarkMulBig-4                         247996              4614 ns/op
BenchmarkModMulBig-4                      216322              5486 ns/op
BenchmarkLargeModMulBig-4                 123004              9443 ns/op
BenchmarkModBig-4                        1307606               929.3 ns/op
BenchmarkLargeModBig-4                    238051              4924 ns/op
BenchmarkModInverseBig-4                  868454              1362 ns/op
BenchmarkLargeModInverseBig-4              59748             20067 ns/op
BenchmarkExpBig-4                           7464            136994 ns/op
BenchmarkLargeExpBig-4                        22          50682980 ns/op
BenchmarkSetBytesBig-4                   4736103               214.3 ns/op
BenchmarkAddNat-4                        5792606               189.2 ns/op
BenchmarkModAddNat-4                       32984             36057 ns/op
BenchmarkLargeModAddNat-4                   1731            676239 ns/op
BenchmarkMulNat-4                          78020             15072 ns/op
BenchmarkModMulNat-4                       32888             35863 ns/op
BenchmarkLargeModMulNat-4                   1742            673623 ns/op
BenchmarkModNat-4                          33306             35515 ns/op
BenchmarkLargeModNat-4                      1496            670187 ns/op
BenchmarkModInverseNat-4                   29176             40850 ns/op
BenchmarkLargeModInverseNat-4                160           7335082 ns/op
BenchmarkExpNat-4                            117          10122154 ns/op
BenchmarkLargeExpNat-4                         1        8909942141 ns/op
BenchmarkSetBytesNat-4                    809414              1403 ns/op
PASS
ok      github.com/cronokirby/safenum   44.924s
```
