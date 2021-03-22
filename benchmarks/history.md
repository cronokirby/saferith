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
```
