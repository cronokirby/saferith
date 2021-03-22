## b51537ffd0a710b798e8adb5692df30fc299da80 (2021-03-31)

```
[safenum] → go test -bench=. -benchtime=5s
goos: linux
goarch: amd64
pkg: github.com/cronokirby/safenum
cpu: Intel(R) Core(TM) i5-4690K CPU @ 3.50GHz
BenchmarkAddBig-4               37461423               150.5 ns/op
BenchmarkModAddBig-4             5770950              1134 ns/op
BenchmarkMulBig-4                2161972              2797 ns/op
BenchmarkModMulBig-4             1600071              4151 ns/op
BenchmarkModBig-4                6910508              1127 ns/op
BenchmarkModInverseBig-4         4563404              1582 ns/op
BenchmarkExpBig-4                  45710            133276 ns/op
BenchmarkSetBytesBig-4          32617444               211.1 ns/op
BenchmarkAddNat-4               43392159               128.9 ns/op
BenchmarkModAddNat-4               71904             82362 ns/op
BenchmarkMulNat-4                 755776              7755 ns/op
BenchmarkModMulNat-4               66014             90602 ns/op
BenchmarkModNat-4                 143535             41030 ns/op
BenchmarkModInverseNat-4          127548             46615 ns/op
BenchmarkExpNat-4                    518          11763090 ns/op
BenchmarkSetBytesNat-4           4154806              1633 ns/op
PASS
ok      github.com/cronokirby/safenum   123.119s
```
