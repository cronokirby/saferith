# safenum

The purpose of this package is to provide a version of arbitrary sized
arithmetic, in a safe (i.e. constant-time) way, for cryptography.

*This is experimental software, use at your own peril*.


# Benchmarks

Run with assembly routines:

```
go test -bench=.
```

Run with pure Go code:

```
go test -bench=. -tags math_big_pure_go
```
