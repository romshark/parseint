package parseint_test

import (
	"math"
	"math/big"
	"runtime"
	"strconv"
	"testing"

	"github.com/romshark/parseint"
	"github.com/stretchr/testify/require"
)

var validBase10Int64 = map[string]int64{
	"0":                   0,
	"1":                   1,
	"10":                  10,
	"100":                 100,
	"1000":                1000,
	"10000":               10000,
	"100000":              100000,
	"1000000":             1000000,
	"10000000":            10000000,
	"100000000":           100000000,
	"1000000000":          1000000000,
	"10000000000":         10000000000,
	"100000000000":        100000000000,
	"1000000000000":       1000000000000,
	"10000000000000":      10000000000000,
	"100000000000000":     100000000000000,
	"1000000000000000":    1000000000000000,
	"10000000000000000":   10000000000000000,
	"100000000000000000":  100000000000000000,
	"1000000000000000000": 1000000000000000000,

	"01":                               1,
	"001":                              1,
	"0001":                             1,
	"00001":                            1,
	"000001":                           1,
	"0000001":                          1,
	"00000001":                         1,
	"000000001":                        1,
	"0000000001":                       1,
	"00000000001":                      1,
	"000000000001":                     1,
	"0000000000001":                    1,
	"00000000000001":                   1,
	"000000000000001":                  1,
	"0000000000000001":                 1,
	"00000000000000001":                1,
	"000000000000000001":               1,
	"0000000000000000001":              1,
	"00000000000000000001":             1,
	"000000000000000000001":            1,
	"0000000000000000000001":           1,
	"00000000000000000000001":          1,
	"000000000000000000000001":         1,
	"0000000000000000000000001":        1,
	"00000000000000000000000001":       1,
	"000000000000000000000000001":      1,
	"0000000000000000000000000001":     1,
	"00000000000000000000000000001":    1,
	"000000000000000000000000000001":   1,
	"0000000000000000000000000000001":  1,
	"00000000000000000000000000000001": 1,

	"12":                  12,
	"123":                 123,
	"1234":                1234,
	"999999999":           999999999,
	"1234567890":          1234567890,
	"4294967294":          4294967294,
	"9":                   9,
	"92":                  92,
	"922":                 922,
	"9223":                9223,
	"92233":               92233,
	"922337":              922337,
	"9223372":             9223372,
	"92233720":            92233720,
	"922337203":           922337203,
	"9223372036":          9223372036,
	"92233720368":         92233720368,
	"922337203685":        922337203685,
	"9223372036854":       9223372036854,
	"92233720368547":      92233720368547,
	"922337203685477":     922337203685477,
	"9223372036854775":    9223372036854775,
	"92233720368547758":   92233720368547758,
	"922337203685477580":  922337203685477580,
	"9223372036854775807": math.MaxInt64,

	"+12":                   12,
	"+123":                  123,
	"+1234":                 1234,
	"+999999999":            999999999,
	"+1234567890":           1234567890,
	"+4294967294":           4294967294,
	"+9":                    9,
	"+92":                   92,
	"+922":                  922,
	"+9223":                 9223,
	"+92233":                92233,
	"+922337":               922337,
	"+9223372":              9223372,
	"+92233720":             92233720,
	"+922337203":            922337203,
	"+9223372036":           9223372036,
	"+92233720368":          92233720368,
	"+922337203685":         922337203685,
	"+9223372036854":        9223372036854,
	"+92233720368547":       92233720368547,
	"+922337203685477":      922337203685477,
	"+9223372036854775":     9223372036854775,
	"+92233720368547758":    92233720368547758,
	"+922337203685477580":   922337203685477580,
	"+9223372036854775807":  math.MaxInt64,
	"+09223372036854775807": math.MaxInt64,

	"-12":                   -12,
	"-123":                  -123,
	"-1234":                 -1234,
	"-999999999":            -999999999,
	"-1234567890":           -1234567890,
	"-4294967294":           -4294967294,
	"-9":                    -9,
	"-92":                   -92,
	"-922":                  -922,
	"-9223":                 -9223,
	"-92233":                -92233,
	"-922337":               -922337,
	"-9223372":              -9223372,
	"-92233720":             -92233720,
	"-922337203":            -922337203,
	"-9223372036":           -9223372036,
	"-92233720368":          -92233720368,
	"-922337203685":         -922337203685,
	"-9223372036854":        -9223372036854,
	"-92233720368547":       -92233720368547,
	"-922337203685477":      -922337203685477,
	"-9223372036854775":     -9223372036854775,
	"-92233720368547758":    -92233720368547758,
	"-922337203685477580":   -922337203685477580,
	"-9223372036854775808":  math.MinInt,
	"-09223372036854775808": math.MinInt,
}

var invalidBase10Int64 = map[string]error{
	"":     parseint.ErrSyntax,
	" ":    parseint.ErrSyntax,
	" 1":   parseint.ErrSyntax,
	" +1":  parseint.ErrSyntax,
	" -1":  parseint.ErrSyntax,
	"-":    parseint.ErrSyntax,
	"+":    parseint.ErrSyntax,
	"-x":   parseint.ErrSyntax,
	"0x0":  parseint.ErrSyntax,
	"123x": parseint.ErrSyntax,

	"ж":  parseint.ErrSyntax,
	"🙂":  parseint.ErrSyntax,
	"🗿":  parseint.ErrSyntax,
	"♻︎": parseint.ErrSyntax,

	".9223372036854775807": parseint.ErrSyntax,
	"9.223372036854775807": parseint.ErrSyntax,
	"92.23372036854775807": parseint.ErrSyntax,
	"922.3372036854775807": parseint.ErrSyntax,
	"9223.372036854775807": parseint.ErrSyntax,
	"92233.72036854775807": parseint.ErrSyntax,
	"922337.2036854775807": parseint.ErrSyntax,
	"9223372.036854775807": parseint.ErrSyntax,
	"92233720.36854775807": parseint.ErrSyntax,
	"922337203.6854775807": parseint.ErrSyntax,
	"9223372036.854775807": parseint.ErrSyntax,
	"92233720368.54775807": parseint.ErrSyntax,
	"922337203685.4775807": parseint.ErrSyntax,
	"9223372036854.775807": parseint.ErrSyntax,
	"92233720368547.75807": parseint.ErrSyntax,
	"922337203685477.5807": parseint.ErrSyntax,
	"9223372036854775.807": parseint.ErrSyntax,
	"92233720368547758.07": parseint.ErrSyntax,
	"922337203685477580.7": parseint.ErrSyntax,
	"9223372036854775807.": parseint.ErrSyntax,

	".922337203685477580700000000": parseint.ErrSyntax,
	"9.22337203685477580700000000": parseint.ErrSyntax,
	"92.2337203685477580700000000": parseint.ErrSyntax,
	"922.337203685477580700000000": parseint.ErrSyntax,
	"9223.37203685477580700000000": parseint.ErrSyntax,
	"92233.7203685477580700000000": parseint.ErrSyntax,
	"922337.203685477580700000000": parseint.ErrSyntax,
	"9223372.03685477580700000000": parseint.ErrSyntax,
	"92233720.3685477580700000000": parseint.ErrSyntax,
	"922337203.685477580700000000": parseint.ErrSyntax,
	"9223372036.85477580700000000": parseint.ErrSyntax,
	"92233720368.5477580700000000": parseint.ErrSyntax,
	"922337203685.477580700000000": parseint.ErrSyntax,
	"9223372036854.77580700000000": parseint.ErrSyntax,
	"92233720368547.7580700000000": parseint.ErrSyntax,
	"922337203685477.580700000000": parseint.ErrSyntax,
	"9223372036854775.80700000000": parseint.ErrSyntax,
	"92233720368547758.0700000000": parseint.ErrSyntax,
	"922337203685477580.700000000": parseint.ErrSyntax,
	"9223372036854775807.00000000": parseint.ErrSyntax,

	"-9223372036854775808.": parseint.ErrSyntax,
	"-922337203685477580.8": parseint.ErrSyntax,
	"-92233720368547758.08": parseint.ErrSyntax,
	"-9223372036854775.808": parseint.ErrSyntax,
	"-922337203685477.5808": parseint.ErrSyntax,
	"-92233720368547.75808": parseint.ErrSyntax,
	"-9223372036854.775808": parseint.ErrSyntax,
	"-922337203685.4775808": parseint.ErrSyntax,
	"-92233720368.54775808": parseint.ErrSyntax,
	"-9223372036.854775808": parseint.ErrSyntax,
	"-922337203.6854775808": parseint.ErrSyntax,
	"-92233720.36854775808": parseint.ErrSyntax,
	"-9223372.036854775808": parseint.ErrSyntax,
	"-922337.2036854775808": parseint.ErrSyntax,
	"-92233.72036854775808": parseint.ErrSyntax,
	"-9223.372036854775808": parseint.ErrSyntax,
	"-922.3372036854775808": parseint.ErrSyntax,
	"-92.23372036854775808": parseint.ErrSyntax,
	"-9.223372036854775808": parseint.ErrSyntax,
	"-.9223372036854775808": parseint.ErrSyntax,

	"-9223372036854775808.00000000": parseint.ErrSyntax,
	"-922337203685477580.800000000": parseint.ErrSyntax,
	"-92233720368547758.0800000000": parseint.ErrSyntax,
	"-9223372036854775.80800000000": parseint.ErrSyntax,
	"-922337203685477.580800000000": parseint.ErrSyntax,
	"-92233720368547.7580800000000": parseint.ErrSyntax,
	"-9223372036854.77580800000000": parseint.ErrSyntax,
	"-922337203685.477580800000000": parseint.ErrSyntax,
	"-92233720368.5477580800000000": parseint.ErrSyntax,
	"-9223372036.85477580800000000": parseint.ErrSyntax,
	"-922337203.685477580800000000": parseint.ErrSyntax,
	"-92233720.3685477580800000000": parseint.ErrSyntax,
	"-9223372.03685477580800000000": parseint.ErrSyntax,
	"-922337.203685477580800000000": parseint.ErrSyntax,
	"-92233.7203685477580800000000": parseint.ErrSyntax,
	"-9223.37203685477580800000000": parseint.ErrSyntax,
	"-922.337203685477580800000000": parseint.ErrSyntax,
	"-92.2337203685477580800000000": parseint.ErrSyntax,
	"-9.22337203685477580800000000": parseint.ErrSyntax,
	"-.922337203685477580800000000": parseint.ErrSyntax,

	// Overflow
	"9223372036854775808":              parseint.ErrOverflow,
	"+9223372036854775808":             parseint.ErrOverflow,
	"9323372036854775807":              parseint.ErrOverflow,
	"19223372036854775808":             parseint.ErrOverflow,
	"11111111111111111111":             parseint.ErrOverflow,
	"12331232123123123123123123123123": parseint.ErrOverflow,

	"92233720368547758070":          parseint.ErrOverflow,
	"922337203685477580700":         parseint.ErrOverflow,
	"9223372036854775807000":        parseint.ErrOverflow,
	"92233720368547758070000":       parseint.ErrOverflow,
	"922337203685477580700000":      parseint.ErrOverflow,
	"9223372036854775807000000":     parseint.ErrOverflow,
	"92233720368547758070000000":    parseint.ErrOverflow,
	"922337203685477580700000000":   parseint.ErrOverflow,
	"9223372036854775807000000000":  parseint.ErrOverflow,
	"92233720368547758070000000000": parseint.ErrOverflow,

	"-92233720368547758080":          parseint.ErrOverflow,
	"-922337203685477580800":         parseint.ErrOverflow,
	"-9223372036854775808000":        parseint.ErrOverflow,
	"-92233720368547758080000":       parseint.ErrOverflow,
	"-922337203685477580800000":      parseint.ErrOverflow,
	"-9223372036854775808000000":     parseint.ErrOverflow,
	"-92233720368547758080000000":    parseint.ErrOverflow,
	"-922337203685477580800000000":   parseint.ErrOverflow,
	"-9223372036854775808000000000":  parseint.ErrOverflow,
	"-92233720368547758080000000000": parseint.ErrOverflow,
}

func TestBase10Int64(t *testing.T) {
	callBase10Int64 := func(input string, fn func(any, error)) {
		fn(parseint.Base10Int64(input))
		fn(parseint.Base10Int64([]byte(input)))
	}

	requireOK := func(t *testing.T, expect int64, input string) {
		callBase10Int64(input, func(actual any, err error) {
			require.NoError(t, err)
			switch actual := actual.(type) {
			case int64:
				require.Equal(t, int64(expect), actual)
			default:
				t.Fatalf("unexpected type: %T", actual)
			}
		})
	}

	t.Run("valid", func(t *testing.T) {
		for input, expecxt := range validBase10Int64 {
			requireOK(t, expecxt, input)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for input, expectedErr := range invalidBase10Int64 {
			callBase10Int64(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr, "%q", input)
				require.Zero(t, a)
			})
		}
	})

	t.Run("range_min_10k", func(t *testing.T) {
		min := int64(math.MinInt64)
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
		mid_neg := int64(math.MinInt64) / 2
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

	t.Run("err_overflow_neg", func(t *testing.T) {
		min := new(big.Int).SetInt64(math.MinInt64)
		start := new(big.Int).Sub(min, big.NewInt(10_000))
		end := new(big.Int).Sub(min, big.NewInt(1))
		delta := big.NewInt(1)
		for i := new(big.Int).Set(start); i.Cmp(end) <= 0; i.Add(i, delta) {
			str := i.String()
			callBase10Int64(str, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})

	t.Run("err_overflow_pos", func(t *testing.T) {
		max := new(big.Int).SetUint64(math.MaxInt64)
		start := new(big.Int).Add(max, big.NewInt(1))
		end := new(big.Int).Add(max, big.NewInt(10_000))
		delta := big.NewInt(1)
		for i := new(big.Int).Set(start); i.Cmp(end) <= 0; i.Add(i, delta) {
			str := i.String()
			callBase10Int64(str, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})
}

func FuzzBase10Int64(f *testing.F) {
	for input := range validBase10Int64 {
		f.Add(input)
	}
	for input := range invalidBase10Int64 {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base10Uint64[string](s)
		std, errStd := strconv.ParseUint(s, 10, 64)
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
			if _, err := strconv.ParseUint(s, 10, 64); err == nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			}
		}
	})
}

func BenchmarkBase10Int64(b *testing.B) {
	var fn func(string) (int64, error)
	var fnBytes func([]byte) (int64, error)
	switch *fBenchmarkFn {
	case BenchmarkFnStrconv:
		fn = func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		}
		fnBytes = func(s []byte) (int64, error) {
			return strconv.ParseInt(string(s), 10, 64)
		}
	case BenchmarkFnParseint:
		fn = parseint.Base10Int64[string]
		fnBytes = parseint.Base10Int64[[]byte]
	default:
		b.Fatalf("unknown benchmark function: %q", *fBenchmarkFn)
	}

	var a int64
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"min", "-9223372036854775808"},
		{"zero", "0"},
		{"small_3", "987"},
		{"small_4", "9871"},
		{"n13", "-9876543210123"},
		{"p13", "9876543210123"},
		{"max", "9223372036854775807"},
		{"syntax", "0.000000000000001"},
		{"overflow", "9223372036854775808"},
		{"overflow_len", "18446744073709551616"},
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
