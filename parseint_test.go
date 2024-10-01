package parseint_test

import "flag"

var fBenchmarkFn = flag.String(
	"benchfunc",
	BenchmarkFnParseint,
	`function to benchmark, use either "strconv" or "parseint"`)

const (
	BenchmarkFnStrconv  = "strconv"
	BenchmarkFnParseint = "parseint"
)
