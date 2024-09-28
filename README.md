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
purpose implementations but they're not optimally efficient for common cases such as
base10-64bit, hex-16bit, etc. On error `strconv.ParseInt` allocates, which makes it
unnecessarily slow when dealing with invalid values and unsuitable for situations
when dynamic memory allocation is unacceptable.

## Benchmark

- On X86 processors `parseint` achieves a geomean* of around `-79.43%`.
- On Apple M1 processors `parseint` achieves a geomean* of around `-73.48%`.

\* _on average across benchmarks_

Use `./bench.sh . 8` to execute benchmark and compare results.

### Linux AMD64

```
goos: linux
goarch: amd64
pkg: github.com/romshark/parseint
cpu: AMD Ryzen 7 5700X 8-Core Processor
                                       │ .strconv_..txt │          .parseint_..txt           │
                                       │     sec/op     │   sec/op     vs base               │
Base16Uint16/min/string-16                 5.258n ±  1%   2.829n ± 1%  -46.21% (p=0.000 n=8)
Base16Uint16/min/bytes-16                  6.574n ±  1%   2.848n ± 1%  -56.68% (p=0.000 n=8)
Base16Uint16/max_low/string-16             8.322n ±  1%   3.724n ± 1%  -55.25% (p=0.000 n=8)
Base16Uint16/max_low/bytes-16             10.705n ±  1%   3.943n ± 1%  -63.16% (p=0.000 n=8)
Base16Uint16/max_upp/string-16             8.324n ±  1%   3.725n ± 1%  -55.25% (p=0.000 n=8)
Base16Uint16/max_upp/bytes-16             10.710n ±  0%   3.943n ± 1%  -63.18% (p=0.000 n=8)
Base16Uint16/syntax/string-16             84.630n ±  1%   3.506n ± 1%  -95.86% (p=0.000 n=8)
Base16Uint16/syntax/bytes-16              89.940n ± 14%   3.457n ± 1%  -96.16% (p=0.000 n=8)
Base16Uint16/overflow/string-16           87.480n ±  4%   3.504n ± 1%  -95.99% (p=0.000 n=8)
Base16Uint16/overflow/bytes-16            95.095n ±  1%   3.045n ± 1%  -96.80% (p=0.000 n=8)
Base16Uint16_uint16/min/string-16          5.261n ±  1%   3.289n ± 1%  -37.48% (p=0.000 n=8)
Base16Uint16_uint16/min/bytes-16           6.574n ±  1%   3.286n ± 1%  -50.02% (p=0.000 n=8)
Base16Uint16_uint16/max_low/string-16      8.209n ±  1%   3.726n ± 1%  -54.62% (p=0.000 n=8)
Base16Uint16_uint16/max_low/bytes-16      10.700n ±  1%   3.725n ± 1%  -65.18% (p=0.000 n=8)
Base16Uint16_uint16/max_upp/string-16      8.321n ±  1%   3.725n ± 1%  -55.23% (p=0.000 n=8)
Base16Uint16_uint16/max_upp/bytes-16      10.700n ±  1%   3.724n ± 1%  -65.19% (p=0.000 n=8)
Base16Uint16_uint16/syntax/string-16      84.525n ±  5%   3.506n ± 1%  -95.85% (p=0.000 n=8)
Base16Uint16_uint16/syntax/bytes-16      101.750n ±  1%   3.481n ± 1%  -96.58% (p=0.000 n=8)
Base16Uint16_uint16/overflow/string-16    87.615n ±  3%   3.288n ± 1%  -96.25% (p=0.000 n=8)
Base16Uint16_uint16/overflow/bytes-16    105.450n ± 25%   3.288n ± 1%  -96.88% (p=0.000 n=8)
Base10Uint32/l1/string-16                  5.256n ±  1%   2.631n ± 1%  -49.94% (p=0.000 n=8)
Base10Uint32/l1/bytes-16                   6.573n ±  1%   2.629n ± 1%  -60.00% (p=0.000 n=8)
Base10Uint32/l3/string-16                  7.013n ±  0%   3.726n ± 1%  -46.87% (p=0.000 n=8)
Base10Uint32/l3/bytes-16                   9.206n ±  1%   3.726n ± 1%  -59.53% (p=0.000 n=8)
Base10Uint32/l6/string-16                  9.637n ±  1%   5.521n ± 1%  -42.71% (p=0.000 n=8)
Base10Uint32/l6/bytes-16                  11.880n ±  1%   5.561n ± 1%  -53.19% (p=0.000 n=8)
Base10Uint32/max/string-16                13.135n ±  1%   7.874n ± 1%  -40.05% (p=0.000 n=8)
Base10Uint32/max/bytes-16                 15.985n ±  1%   7.912n ± 2%  -50.51% (p=0.000 n=8)
Base10Uint32/syntax/string-16             58.045n ± 28%   2.192n ± 1%  -96.22% (p=0.000 n=8)
Base10Uint32/syntax/bytes-16              60.810n ±  3%   2.190n ± 1%  -96.40% (p=0.000 n=8)
Base10Uint32/overflow/string-16          102.350n ±  1%   7.452n ± 1%  -92.72% (p=0.000 n=8)
Base10Uint32/overflow/bytes-16           108.150n ±  0%   7.346n ± 0%  -93.21% (p=0.000 n=8)
Base10Uint64/min/string-16                 5.260n ±  1%   4.160n ± 1%  -20.90% (p=0.000 n=8)
Base10Uint64/min/bytes-16                  6.574n ±  1%   4.165n ± 1%  -36.65% (p=0.000 n=8)
Base10Uint64/small_3/string-16             7.008n ±  1%   5.222n ± 1%  -25.49% (p=0.000 n=8)
Base10Uint64/small_3/bytes-16              9.418n ±  1%   5.255n ± 1%  -44.20% (p=0.000 n=8)
Base10Uint64/small_4/string-16             7.883n ±  1%   5.005n ± 1%  -36.51% (p=0.000 n=8)
Base10Uint64/small_4/bytes-16             10.075n ±  1%   5.042n ± 1%  -49.96% (p=0.000 n=8)
Base10Uint64/max/string-16                 21.91n ±  1%   12.11n ± 1%  -44.75% (p=0.000 n=8)
Base10Uint64/max/bytes-16                  24.77n ±  1%   11.18n ± 1%  -54.86% (p=0.000 n=8)
Base10Uint64/syntax/string-16             88.600n ± 25%   3.289n ± 1%  -96.29% (p=0.000 n=8)
Base10Uint64/syntax/bytes-16              95.495n ±  1%   3.286n ± 1%  -96.56% (p=0.000 n=8)
Base10Uint64/overflow/string-16           122.25n ±  1%   11.32n ± 2%  -90.74% (p=0.000 n=8)
Base10Uint64/overflow/bytes-16            128.25n ±  1%   10.49n ± 1%  -91.82% (p=0.000 n=8)
geomean                                    20.82n         4.282n       -79.43%
```

### macOS ARM64

```
goos: darwin
goarch: arm64
pkg: github.com/romshark/parseint
                                       │ .strconv_..txt │           .parseint_..txt           │
                                       │     sec/op     │    sec/op     vs base               │
Base16Uint16/min/string-10                  4.967n ± 0%   2.310n ±  0%  -53.49% (p=0.000 n=8)
Base16Uint16/min/bytes-10                   5.899n ± 0%   2.460n ± 10%  -58.30% (p=0.000 n=8)
Base16Uint16/max_low/string-10              8.384n ± 0%   3.247n ±  0%  -61.27% (p=0.000 n=8)
Base16Uint16/max_low/bytes-10              10.870n ± 0%   3.188n ±  0%  -70.67% (p=0.000 n=8)
Base16Uint16/max_upp/string-10              8.383n ± 0%   3.248n ±  1%  -61.26% (p=0.000 n=8)
Base16Uint16/max_upp/bytes-10              10.870n ± 0%   3.197n ±  0%  -70.59% (p=0.000 n=8)
Base16Uint16/syntax/string-10              40.275n ± 0%   3.260n ±  0%  -91.90% (p=0.000 n=8)
Base16Uint16/syntax/bytes-10               43.345n ± 0%   3.186n ±  0%  -92.65% (p=0.000 n=8)
Base16Uint16/overflow/string-10            42.805n ± 0%   2.834n ±  1%  -93.38% (p=0.000 n=8)
Base16Uint16/overflow/bytes-10             45.285n ± 0%   2.839n ±  1%  -93.73% (p=0.000 n=8)
Base16Uint16_uint16/min/string-10           4.966n ± 0%   2.306n ±  0%  -53.57% (p=0.000 n=8)
Base16Uint16_uint16/min/bytes-10            5.898n ± 0%   2.545n ±  9%  -56.85% (p=0.000 n=8)
Base16Uint16_uint16/max_low/string-10       8.383n ± 0%   3.240n ±  1%  -61.34% (p=0.000 n=8)
Base16Uint16_uint16/max_low/bytes-10       10.870n ± 0%   3.168n ±  1%  -70.86% (p=0.000 n=8)
Base16Uint16_uint16/max_upp/string-10       8.384n ± 0%   3.248n ±  1%  -61.25% (p=0.000 n=8)
Base16Uint16_uint16/max_upp/bytes-10       10.870n ± 0%   3.158n ±  1%  -70.95% (p=0.000 n=8)
Base16Uint16_uint16/syntax/string-10       40.310n ± 0%   3.226n ±  1%  -92.00% (p=0.000 n=8)
Base16Uint16_uint16/syntax/bytes-10        43.395n ± 0%   3.150n ±  1%  -92.74% (p=0.000 n=8)
Base16Uint16_uint16/overflow/string-10     42.765n ± 0%   2.796n ±  1%  -93.46% (p=0.000 n=8)
Base16Uint16_uint16/overflow/bytes-10      45.290n ± 0%   2.798n ±  0%  -93.82% (p=0.000 n=8)
Base10Uint32/l1/string-10                   4.691n ± 2%   2.484n ±  1%  -47.05% (p=0.000 n=8)
Base10Uint32/l1/bytes-10                    5.728n ± 1%   2.483n ±  0%  -56.65% (p=0.000 n=8)
Base10Uint32/l3/string-10                   7.313n ± 0%   5.404n ±  0%  -26.10% (p=0.000 n=8)
Base10Uint32/l3/bytes-10                   10.165n ± 0%   5.402n ±  0%  -46.86% (p=0.000 n=8)
Base10Uint32/l6/string-10                  10.175n ± 0%   7.451n ±  0%  -26.77% (p=0.000 n=8)
Base10Uint32/l6/bytes-10                   12.700n ± 2%   7.449n ±  0%  -41.35% (p=0.000 n=8)
Base10Uint32/max/string-10                 13.985n ± 0%   9.933n ±  0%  -28.97% (p=0.000 n=8)
Base10Uint32/max/bytes-10                  15.840n ± 1%   9.932n ±  0%  -37.29% (p=0.000 n=8)
Base10Uint32/syntax/string-10              26.570n ± 0%   2.174n ±  0%  -91.82% (p=0.000 n=8)
Base10Uint32/syntax/bytes-10               27.780n ± 0%   2.179n ±  1%  -92.16% (p=0.000 n=8)
Base10Uint32/overflow/string-10            47.020n ± 0%   8.911n ±  0%  -81.05% (p=0.000 n=8)
Base10Uint32/overflow/bytes-10             49.735n ± 0%   8.932n ±  0%  -82.04% (p=0.000 n=8)
Base10Uint64/min/string-10                  4.693n ± 1%   4.034n ±  0%  -14.06% (p=0.000 n=8)
Base10Uint64/min/bytes-10                   5.617n ± 1%   4.035n ±  0%  -28.16% (p=0.000 n=8)
Base10Uint64/small_3/string-10              7.141n ± 1%   7.039n ±  0%   -1.42% (p=0.000 n=8)
Base10Uint64/small_3/bytes-10               9.933n ± 0%   7.039n ±  0%  -29.14% (p=0.000 n=8)
Base10Uint64/small_4/string-10              8.072n ± 0%   4.345n ±  0%  -46.17% (p=0.000 n=8)
Base10Uint64/small_4/bytes-10              10.560n ± 0%   4.341n ±  0%  -58.89% (p=0.000 n=8)
Base10Uint64/max/string-10                  24.54n ± 0%   12.09n ±  0%  -50.72% (p=0.000 n=8)
Base10Uint64/max/bytes-10                   26.73n ± 0%   12.21n ±  0%  -54.32% (p=0.000 n=8)
Base10Uint64/syntax/string-10              41.060n ± 0%   2.663n ±  1%  -93.52% (p=0.000 n=8)
Base10Uint64/syntax/bytes-10               43.625n ± 0%   2.708n ±  2%  -93.79% (p=0.000 n=8)
Base10Uint64/overflow/string-10             59.27n ± 0%   11.51n ±  0%  -80.58% (p=0.000 n=8)
Base10Uint64/overflow/bytes-10              61.73n ± 0%   11.51n ±  0%  -81.35% (p=0.000 n=8)
geomean                                     15.73n        4.173n        -73.48%
```
