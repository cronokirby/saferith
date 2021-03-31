# safenum

The purpose of this package is to provide a version of arbitrary sized
arithmetic, in a safe (i.e. constant-time) way, for cryptography.

*This is experimental software, use at your own peril*.

This code is structured to be easily moved inside of (a fork of)
Go's standard library, hence the `arith*.go` files, which are copied from there.
If you wanted to use this in a more standalone context, you'd likely want to salvage
only the necessary internal routines.

# Benchmarks

Run with assembly routines:

```
go test -bench=.
```

Run with pure Go code:

```
go test -bench=. -tags math_big_pure_go
```

# Licensing

The files `arith*.go` come from Go's standard library, and are licensed under
a BSD license in `LICENSE_go`. The rest of the code is under an MIT license.
