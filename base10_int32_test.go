package parseint_test

import (
	"math"
	"runtime"
	"strconv"
	"testing"

	"github.com/romshark/parseint"
	"github.com/stretchr/testify/require"
)

var validBase10Int32 = map[string]int32{
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
	"2147483647":                       math.MaxInt32,
	"02147483647":                      math.MaxInt32,
	"01":                               1,
	"00000000000000000000000000000001": 1,
	"0000":                             0,
	"00000000000000000000000000000000": 0,

	"-0":                                0,
	"-1":                                -1,
	"-12":                               -12,
	"-123":                              -123,
	"-1234":                             -1234,
	"-12345":                            -12345,
	"-123456":                           -123456,
	"-1234567":                          -1234567,
	"-12345678":                         -12345678,
	"-999999999":                        -999999999,
	"-1234567890":                       -1234567890,
	"-01":                               -1,
	"-00000000000000000000000000000001": -1,
	"-0000":                             0,
	"-00000000000000000000000000000000": 0,
	"-2147483648":                       math.MinInt32,
	"-02147483648":                      math.MinInt32,
}

var invalidBase10Int32 = map[string]error{
	"":     parseint.ErrSyntax,
	" ":    parseint.ErrSyntax,
	" 1":   parseint.ErrSyntax,
	" -1":  parseint.ErrSyntax,
	" +1":  parseint.ErrSyntax,
	"\t":   parseint.ErrSyntax,
	"\t1":  parseint.ErrSyntax,
	"-":    parseint.ErrSyntax,
	"+":    parseint.ErrSyntax,
	"a":    parseint.ErrSyntax,
	"-a":   parseint.ErrSyntax,
	"af":   parseint.ErrSyntax,
	"aaaa": parseint.ErrSyntax,
	"1a2b": parseint.ErrSyntax,
	"FFFF": parseint.ErrSyntax,
	"abcd": parseint.ErrSyntax,
	"eeFF": parseint.ErrSyntax,
	"defg": parseint.ErrSyntax,
	"xyz":  parseint.ErrSyntax,
	"GHIJ": parseint.ErrSyntax,

	"ж":  parseint.ErrSyntax,
	"🙂":  parseint.ErrSyntax,
	"🗿":  parseint.ErrSyntax,
	"♻︎": parseint.ErrSyntax,

	"2147483648":                       parseint.ErrOverflow,
	"4294967296":                       parseint.ErrOverflow,
	"9999999999":                       parseint.ErrOverflow,
	"123456789123456789":               parseint.ErrOverflow,
	"00000000000000000000002147483648": parseint.ErrOverflow,

	"+2147483648":                       parseint.ErrOverflow,
	"+4294967296":                       parseint.ErrOverflow,
	"+9999999999":                       parseint.ErrOverflow,
	"+123456789123456789":               parseint.ErrOverflow,
	"+00000000000000000000002147483648": parseint.ErrOverflow,

	"-2147483649":                       parseint.ErrOverflow,
	"-4294967296":                       parseint.ErrOverflow,
	"-9999999999":                       parseint.ErrOverflow,
	"-123456789123456789":               parseint.ErrOverflow,
	"-00000000000000000000002147483649": parseint.ErrOverflow,
}

func TestBase10Int32(t *testing.T) {
	callBase10Int32 := func(input string, fn func(any, error)) {
		fn(parseint.Base10Int32[string, int64](input))
		fn(parseint.Base10Int32[string, int32](input))
		fn(parseint.Base10Int32[[]byte, int64]([]byte(input)))
		fn(parseint.Base10Int32[[]byte, int32]([]byte(input)))
	}

	requireOK := func(t *testing.T, expect int64, input string) {
		callBase10Int32(input, func(actual any, err error) {
			require.NoError(t, err)
			switch actual := actual.(type) {
			case int64:
				require.Equal(t, int64(expect), actual)
			case int32:
				require.Equal(t, int32(expect), actual)
			default:
				t.Fatalf("unexpected type: %T", actual)
			}
		})
	}

	t.Run("valid", func(t *testing.T) {
		for input, expect := range validBase10Int32 {
			requireOK(t, int64(expect), input)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for input, expectedErr := range invalidBase10Int32 {
			callBase10Int32(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr, "%q", input)
				require.Zero(t, a)
			})
		}
	})

	t.Run("range_min_10k", func(t *testing.T) {
		for input, expect := range validBase10Int32 {
			requireOK(t, int64(expect), input)
		}
	})

	t.Run("range_min_10k", func(t *testing.T) {
		min := int64(math.MinInt32)
		for i := min; i <= min+10_000; i++ {
			requireOK(t, i, strconv.FormatInt(i, 10))
		}
	})

	t.Run("range_0_10k", func(t *testing.T) {
		for i := int64(0); i <= 10_000; i++ {
			dec := strconv.FormatInt(i, 10)
			requireOK(t, i, dec)
			requireOK(t, i, "+"+dec)
		}
	})

	t.Run("range_mid_neg_10k", func(t *testing.T) {
		mid_neg := int64(math.MinInt32) / 2
		for i := mid_neg; i <= mid_neg+10_000; i++ {
			requireOK(t, i, strconv.FormatInt(i, 10))
		}
	})

	t.Run("range_mid_pos_10k", func(t *testing.T) {
		mid_pos := int64(math.MaxInt32) / 2
		for i := mid_pos; i <= mid_pos+10_000; i++ {
			dec := strconv.FormatInt(i, 10)
			requireOK(t, i, dec)
			requireOK(t, i, "+"+dec)
		}
	})

	t.Run("range_last10k", func(t *testing.T) {
		max := int64(math.MaxInt32)
		for i := max; i <= max-10_000; i++ {
			dec := strconv.FormatInt(i, 10)
			requireOK(t, i, dec)
			requireOK(t, i, "+"+dec)
		}
	})

	t.Run("range_overflow_pos", func(t *testing.T) {
		for i := int64(math.MaxInt32 + 1); i <= math.MaxInt32+10_000; i++ {
			dec := strconv.FormatInt(i, 10)
			callBase10Int32(dec, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})

	t.Run("range_overflow_neg", func(t *testing.T) {
		for i := int64(math.MinInt32 - 10_000); i < math.MinInt32; i++ {
			dec := strconv.FormatInt(i, 10)
			callBase10Int32(dec, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})
}

func fuzzBase10Int32[I int64 | int32](f *testing.F) {
	for input := range validBase10Int32 {
		f.Add(input)
	}
	for input := range invalidBase10Int32 {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base10Int32[string, I](s)
		std, errStd := strconv.ParseInt(s, 10, 32)
		if err == nil {
			if errStd != nil {
				t.Fatalf("must have returned error %v but didn't: %q", errStd, s)
			} else if std != int64(x) {
				t.Errorf("expected %d; received: %d", std, int64(x))
			}
		} else {
			if x != 0 {
				t.Errorf("%q: failed but returned non-zero value: %x", s, x)
			}
			if _, err := strconv.ParseInt(s, 10, 32); err == nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			}
		}
	})
}

func FuzzBase10Int32_int64(f *testing.F) { fuzzBase10Int32[int64](f) }
func FuzzBase10Int32_int32(f *testing.F) { fuzzBase10Int32[int32](f) }

func BenchmarkBase10Int32(b *testing.B) {
	fn := getBenchmarkFn(b, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 32)
	}, parseint.Base10Int32[string, int64])
	fnBytes := getBenchmarkFn(b, func(s []byte) (int64, error) {
		return strconv.ParseInt(string(s), 10, 32)
	}, parseint.Base10Int32[[]byte, int64])

	var a int64
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"min", "-2147483648"},
		{"neg7", "-429495"},
		{"plus", "+429495"},
		{"pos1", "0"},
		{"pos3", "100"},
		{"pos6", "429495"},
		{"max", "2147483647"},
		{"syntax", "-"},
		{"overflow_min", "-2147483649"},
		{"overflow_max", "2147483648"},
		{"overflow_len", "999999999999999999"},
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
