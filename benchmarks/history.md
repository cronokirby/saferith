## b51537ffd0a710b798e8adb5692df30fc299da80 (2021-03-31)

```
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4                5606130               185.9 ns/op
BenchmarkModAddBig-4             1082536              1206 ns/op
BenchmarkMulBig-4                 434248              2664 ns/op
BenchmarkModMulBig-4              317546              3641 ns/op
BenchmarkModBig-4                1230501              1002 ns/op
BenchmarkModInverseBig-4          742441              1596 ns/op
BenchmarkExpBig-4                   8119            150092 ns/op
BenchmarkSetBytesBig-4           4231164               255.7 ns/op
BenchmarkAddNat-4                5599339               180.1 ns/op
BenchmarkModAddNat-4               13545             89190 ns/op
BenchmarkMulNat-4                 139971              8003 ns/op
BenchmarkModMulNat-4               12583             95224 ns/op
BenchmarkModNat-4                  27632             43629 ns/op
BenchmarkModInverseNat-4           24558             49429 ns/op
BenchmarkExpNat-4                     92          12278909 ns/op
BenchmarkSetBytesNat-4            817094              1466 ns/op
PASS
ok      github.com/cronokirby/safenum   28.026s
```
