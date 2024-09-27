<a href="https://pkg.go.dev/github.com/romshark/parseint">
    <img src="https://godoc.org/github.com/romshark/parseint?status.svg" alt="GoDoc">
</a>
<a href="https://goreportcard.com/report/github.com/romshark/parseint">
    <img src="https://goreportcard.com/badge/github.com/romshark/parseint" alt="GoReportCard">
</a>
<a href='https://coveralls.io/github/romshark/parseint?branch=main'>
    <img src='https://coveralls.io/repos/github/romshark/parseint/badge.svg?branch=main' alt='Coverage Status' />
</a>

# parseint

Package `parseint` is a collection of efficient generic integer parsing functions for Go.
Go's standard library provides [strconv.ParseInt](https://pkg.go.dev/strconv#ParseInt) and
[strconv.ParseUint](https://pkg.go.dev/strconv#ParseUint) which are very useful general
purpose implementations but they're not optimally efficient for common cases.

## Benchmark

Use `./bench.sh <FILTER> <COUNT>` to execute benchmark and compare results.

### base-16 (hexadecimal) 16-bit unsigned integer `"ffff"`:

```sh
./bench.sh Base16Uint16/max_low 10
```

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/parseint
                                     │ .strconv_Base16Uint16_max_low.txt │ .parseint_Base16Uint16_max_low.txt  │
                                     │              sec/op               │   sec/op     vs base                │
Base16Uint16/max_low-10                                      8.429n ± 0%   3.111n ± 0%  -63.10% (p=0.000 n=10)
Base16Uint16/max_low/bytes-10                               10.880n ± 0%   3.112n ± 1%  -71.40% (p=0.000 n=10)
Base16Uint16_uint16/max_low-10                               8.426n ± 0%   2.963n ± 0%  -64.84% (p=0.000 n=10)
Base16Uint16_uint16/max_low/bytes-10                        10.875n ± 0%   3.123n ± 0%  -71.29% (p=0.000 n=10)
geomean                                                      9.575n        3.076n       -67.87%
```

### base-16 (hexadecimal) 16-bit unsigned integer `"fffx"` (err: syntax):

```sh
./bench.sh Base16Uint16/syntax 10
```

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/parseint
                                    │ .strconv_Base16Uint16_syntax.txt │  .parseint_Base16Uint16_syntax.txt  │
                                    │              sec/op              │   sec/op     vs base                │
Base16Uint16/syntax-10                                    40.270n ± 0%   3.210n ± 0%  -92.03% (p=0.000 n=10)
Base16Uint16/syntax/bytes-10                              43.445n ± 0%   3.127n ± 0%  -92.80% (p=0.000 n=10)
Base16Uint16_uint16/syntax-10                             40.335n ± 1%   3.192n ± 1%  -92.09% (p=0.000 n=10)
Base16Uint16_uint16/syntax/bytes-10                       43.410n ± 0%   3.147n ± 1%  -92.75% (p=0.000 n=10)
geomean                                                    41.84n        3.169n       -92.43%

                                    │ .strconv_Base16Uint16_syntax.txt │   .parseint_Base16Uint16_syntax.txt    │
                                    │               B/op               │   B/op     vs base                     │
Base16Uint16/syntax-10                                      52.00 ± 0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
Base16Uint16/syntax/bytes-10                                52.00 ± 0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
Base16Uint16_uint16/syntax-10                               52.00 ± 0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
Base16Uint16_uint16/syntax/bytes-10                         52.00 ± 0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
geomean                                                     52.00                   ?                       ¹ ²
¹ summaries must be >0 to compute geomean
² ratios must be >0 to compute geomean

                                    │ .strconv_Base16Uint16_syntax.txt │    .parseint_Base16Uint16_syntax.txt    │
                                    │            allocs/op             │ allocs/op   vs base                     │
Base16Uint16/syntax-10                                      2.000 ± 0%   0.000 ± 0%  -100.00% (p=0.000 n=10)
Base16Uint16/syntax/bytes-10                                2.000 ± 0%   0.000 ± 0%  -100.00% (p=0.000 n=10)
Base16Uint16_uint16/syntax-10                               2.000 ± 0%   0.000 ± 0%  -100.00% (p=0.000 n=10)
Base16Uint16_uint16/syntax/bytes-10                         2.000 ± 0%   0.000 ± 0%  -100.00% (p=0.000 n=10)
geomean                                                     2.000                    ?                       ¹ ²
¹ summaries must be >0 to compute geomean
² ratios must be >0 to compute geomean
```

### base-10 32-bit unsigned integer `"4294967295"`:

```sh
./bench.sh Base10Uint32/max 10
```

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/parseint
                          │ .strconv_Base10Uint32_max.txt │   .parseint_Base10Uint32_max.txt    │
                          │            sec/op             │   sec/op     vs base                │
Base10Uint32/max-10                          13.680n ± 0%   9.636n ± 0%  -29.56% (p=0.000 n=10)
Base10Uint32/max/bytes-10                    15.860n ± 1%   9.645n ± 0%  -39.19% (p=0.000 n=10)
geomean                                       14.73n        9.640n       -34.55%
```

### base-10 32-bit unsigned integer `"99999999999"` (error: overflow):

```sh
./bench.sh Base10Uint32/overflow 10
```

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/parseint
                               │ .strconv_Base10Uint32_overflow.txt │ .parseint_Base10Uint32_overflow.txt │
                               │               sec/op               │   sec/op     vs base                │
Base10Uint32/overflow-10                               47.105n ± 0%   8.941n ± 0%  -81.02% (p=0.000 n=10)
Base10Uint32/overflow/bytes-10                         49.785n ± 0%   8.868n ± 0%  -82.19% (p=0.000 n=10)
geomean                                                 48.43n        8.904n       -81.61%

                               │ .strconv_Base10Uint32_overflow.txt │  .parseint_Base10Uint32_overflow.txt   │
                               │                B/op                │   B/op     vs base                     │
Base10Uint32/overflow-10                                 64.00 ± 0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
Base10Uint32/overflow/bytes-10                           64.00 ± 0%   0.00 ± 0%  -100.00% (p=0.000 n=10)
geomean                                                  64.00                   ?                       ¹ ²
¹ summaries must be >0 to compute geomean
² ratios must be >0 to compute geomean

                               │ .strconv_Base10Uint32_overflow.txt │   .parseint_Base10Uint32_overflow.txt   │
                               │             allocs/op              │ allocs/op   vs base                     │
Base10Uint32/overflow-10                                 2.000 ± 0%   0.000 ± 0%  -100.00% (p=0.000 n=10)
Base10Uint32/overflow/bytes-10                           2.000 ± 0%   0.000 ± 0%  -100.00% (p=0.000 n=10)
geomean                                                  2.000                    ?                       ¹ ²
¹ summaries must be >0 to compute geomean
² ratios must be >0 to compute geomean
```
