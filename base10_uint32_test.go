package parseint_test

import (
	"math"
	"runtime"
	"strconv"
	"testing"

	"github.com/romshark/parseint"
	"github.com/stretchr/testify/require"
)

var validBase10Uint32 = map[string]uint32{
	"0":                                0,
	"1":                                1,
	"12":                               12,
	"123":                              123,
	"1234":                             1234,
	"12345":                            12345,
	"123456":                           123456,
	"1234567":                          1234567,
	"12345678":                         12345678,
	"999999999":                        999999999,
	"1234567890":                       1234567890,
	"4294967295":                       math.MaxUint32,
	"04294967295":                      math.MaxUint32,
	"01":                               1,
	"00000000000000000000000000000001": 1,
	"0000":                             0,
	"00000000000000000000000000000000": 0,
}

var invalidBase10Uint32 = map[string]error{
	"":   parseint.ErrSyntax,
	" ":  parseint.ErrSyntax,
	" 1": parseint.ErrSyntax,

	"-":  parseint.ErrSyntax,
	"-1": parseint.ErrSyntax,
	"+1": parseint.ErrSyntax,

	"a":    parseint.ErrSyntax,
	"af":   parseint.ErrSyntax,
	"aaaa": parseint.ErrSyntax,
	"1a2b": parseint.ErrSyntax,
	"FFFF": parseint.ErrSyntax,
	"abcd": parseint.ErrSyntax,
	"eeFF": parseint.ErrSyntax,
	"defg": parseint.ErrSyntax,
	"xyz":  parseint.ErrSyntax,
	"GHIJ": parseint.ErrSyntax,
	"000a": parseint.ErrSyntax,

	"ж":  parseint.ErrSyntax,
	"🙂":  parseint.ErrSyntax,
	"🗿":  parseint.ErrSyntax,
	"♻︎": parseint.ErrSyntax,

	"4294967296":                       parseint.ErrOverflow,
	"9999999999":                       parseint.ErrOverflow,
	"123456789123456789":               parseint.ErrOverflow,
	"00000000000000000000004294967296": parseint.ErrOverflow,
}

func TestBase10Uint32(t *testing.T) {
	callBase10Uint32 := func(input string, fn func(any, error)) {
		fn(parseint.Base10Uint32[string, uint64](input))
		fn(parseint.Base10Uint32[string, uint32](input))
		fn(parseint.Base10Uint32[[]byte, uint64]([]byte(input)))
		fn(parseint.Base10Uint32[[]byte, uint32]([]byte(input)))
	}

	requireOK := func(t *testing.T, expect uint64, input string) {
		callBase10Uint32(input, func(actual any, err error) {
			require.NoError(t, err)
			switch actual := actual.(type) {
			case uint64:
				require.Equal(t, uint64(expect), actual)
			case uint32:
				require.Equal(t, uint32(expect), actual)
			default:
				t.Fatalf("unexpected type: %T", actual)
			}
		})
	}

	t.Run("valid", func(t *testing.T) {
		for input, expect := range validBase10Uint32 {
			requireOK(t, uint64(expect), input)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for input, expectedErr := range invalidBase10Uint32 {
			callBase10Uint32(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr, "%q", input)
				require.Zero(t, a)
			})
		}
	})

	t.Run("range_0_10k", func(t *testing.T) {
		for i := uint64(0); i <= 10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			requireOK(t, i, dec)
		}
	})

	t.Run("range_mid10k", func(t *testing.T) {
		mid := uint64(math.MaxUint32) / 2
		for i := mid; i <= mid+10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			requireOK(t, i, dec)
		}
	})

	t.Run("range_last10k", func(t *testing.T) {
		max := uint64(math.MaxUint32)
		for i := max; i <= max-10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			requireOK(t, i, dec)
		}
	})

	t.Run("err_overflow", func(t *testing.T) {
		for i := uint64(math.MaxUint32 + 1); i <= math.MaxUint32+10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint32(dec, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})
}

func fuzzBase10Uint32[U uint64 | uint32](f *testing.F) {
	for input := range validBase10Uint32 {
		f.Add(input)
	}
	for input := range invalidBase10Uint32 {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base10Uint32[string, U](s)
		std, errStd := strconv.ParseUint(s, 10, 32)
		if err == nil {
			if errStd != nil {
				t.Fatalf("must have returned error %v but didn't: %q", errStd, s)
			} else if std != uint64(x) {
				t.Errorf("expected %d; received: %d", std, uint64(x))
			}
		} else {
			if x != 0 {
				t.Errorf("%q: failed but returned non-zero value: %x", s, x)
			}
			if _, err := strconv.ParseUint(s, 10, 32); err == nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			}
		}
	})
}

func FuzzBase10Uint32_uint64(f *testing.F) { fuzzBase10Uint32[uint64](f) }
func FuzzBase10Uint32_uint32(f *testing.F) { fuzzBase10Uint32[uint32](f) }

func BenchmarkBase10Uint32(b *testing.B) {
	fn := getBenchmarkFn(b, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 32)
	}, parseint.Base10Uint32[string, uint64])
	fnBytes := getBenchmarkFn(b, func(s []byte) (uint64, error) {
		return strconv.ParseUint(string(s), 10, 32)
	}, parseint.Base10Uint32[[]byte, uint64])

	var a uint64
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"l1", "0"},
		{"l3", "100"},
		{"l6", "429495"},
		{"max", "4294967295"},
		{"syntax", "-"},
		{"overflow", "99999999999"},
		{"leadzero31", "00000000000000000000000000000001"},
	} {
		b.Run(td.name+"/string", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				a, err = fn(td.input)
			}
		})
		inputBytes := []byte(td.input)
		b.Run(td.name+"/bytes", func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				a, err = fnBytes(inputBytes)
			}
		})
	}
	runtime.KeepAlive(a)
	runtime.KeepAlive(err)
}
