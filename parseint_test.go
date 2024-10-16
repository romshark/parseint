package parseint_test

import (
	"flag"
	"fmt"
	"testing"
)

var fBenchmarkFn = flag.String(
	"benchfunc",
	BenchmarkFnParseint,
	fmt.Sprintf(
		`function to benchmark, use either %q or %q`,
		BenchmarkFnStrconv, BenchmarkFnParseint,
	),
)

const (
	BenchmarkFnStrconv  = "strconv"
	BenchmarkFnParseint = "parseint"
)

func getBenchmarkFn[I any, S []byte | string](b *testing.B,
	strconvImpl func(s S) (I, error), parseintImpl func(s S) (I, error),
) func(S) (I, error) {
	switch *fBenchmarkFn {
	case BenchmarkFnStrconv:
		return strconvImpl
	case BenchmarkFnParseint:
		return parseintImpl
	default:
		b.Fatalf("unknown benchmark function: %q", *fBenchmarkFn)
	}
	return nil
}
