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

var validBase10Uint64 = map[string]uint64{
	"0":                                0,
	"1":                                1,
	"10":                               10,
	"100":                              100,
	"1000":                             1000,
	"10000":                            10000,
	"100000":                           100000,
	"1000000":                          1000000,
	"10000000":                         10000000,
	"100000000":                        100000000,
	"1000000000":                       1000000000,
	"10000000000":                      10000000000,
	"100000000000":                     100000000000,
	"1000000000000":                    1000000000000,
	"10000000000000":                   10000000000000,
	"100000000000000":                  100000000000000,
	"1000000000000000":                 1000000000000000,
	"10000000000000000":                10000000000000000,
	"100000000000000000":               100000000000000000,
	"1000000000000000000":              1000000000000000000,
	"10000000000000000000":             10000000000000000000,
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
	"12":                               12,
	"123":                              123,
	"1234":                             1234,
	"999999999":                        999999999,
	"1234567890":                       1234567890,
	"4294967294":                       4294967294,
	"184":                              184,
	"1844":                             1844,
	"18446":                            18446,
	"184467":                           184467,
	"1844674":                          1844674,
	"18446744":                         18446744,
	"184467440":                        184467440,
	"1844674407":                       1844674407,
	"18446744073":                      18446744073,
	"184467440737":                     184467440737,
	"1844674407370":                    1844674407370,
	"18446744073709":                   18446744073709,
	"184467440737095":                  184467440737095,
	"1844674407370955":                 1844674407370955,
	"18446744073709551":                18446744073709551,
	"184467440737095516":               184467440737095516,
	"1844674407370955161":              1844674407370955161,
	"18446744073709551615":             math.MaxUint64,
	"018446744073709551615":            math.MaxUint64,
}

var invalidBase10Uint64 = map[string]error{
	"":     parseint.ErrSyntax,
	" ":    parseint.ErrSyntax,
	" 1":   parseint.ErrSyntax,
	"+":    parseint.ErrSyntax,
	"+1":   parseint.ErrSyntax,
	"-":    parseint.ErrSyntax,
	"-0":   parseint.ErrSyntax,
	"0.0":  parseint.ErrSyntax,
	"-x":   parseint.ErrSyntax,
	"0x0":  parseint.ErrSyntax,
	"123x": parseint.ErrSyntax,

	"ж":  parseint.ErrSyntax,
	"🙂":  parseint.ErrSyntax,
	"🗿":  parseint.ErrSyntax,
	"♻︎": parseint.ErrSyntax,

	".18446744073709551615": parseint.ErrSyntax,
	"1.8446744073709551615": parseint.ErrSyntax,
	"18.446744073709551615": parseint.ErrSyntax,
	"184.46744073709551615": parseint.ErrSyntax,
	"1844.6744073709551615": parseint.ErrSyntax,
	"18446.744073709551615": parseint.ErrSyntax,
	"184467.44073709551615": parseint.ErrSyntax,
	"1844674.4073709551615": parseint.ErrSyntax,
	"18446744.073709551615": parseint.ErrSyntax,
	"184467440.73709551615": parseint.ErrSyntax,
	"1844674407.3709551615": parseint.ErrSyntax,
	"18446744073.709551615": parseint.ErrSyntax,
	"184467440737.09551615": parseint.ErrSyntax,
	"1844674407370.9551615": parseint.ErrSyntax,
	"18446744073709.551615": parseint.ErrSyntax,
	"184467440737095.51615": parseint.ErrSyntax,
	"1844674407370955.1615": parseint.ErrSyntax,
	"18446744073709551.615": parseint.ErrSyntax,
	"184467440737095516.15": parseint.ErrSyntax,
	"1844674407370955161.5": parseint.ErrSyntax,
	"18446744073709551615.": parseint.ErrSyntax,

	".1844674407370955161500000000": parseint.ErrSyntax,
	"1.844674407370955161500000000": parseint.ErrSyntax,
	"18.44674407370955161500000000": parseint.ErrSyntax,
	"184.4674407370955161500000000": parseint.ErrSyntax,
	"1844.674407370955161500000000": parseint.ErrSyntax,
	"18446.74407370955161500000000": parseint.ErrSyntax,
	"184467.4407370955161500000000": parseint.ErrSyntax,
	"1844674.407370955161500000000": parseint.ErrSyntax,
	"18446744.07370955161500000000": parseint.ErrSyntax,
	"184467440.7370955161500000000": parseint.ErrSyntax,
	"1844674407.370955161500000000": parseint.ErrSyntax,
	"18446744073.70955161500000000": parseint.ErrSyntax,
	"184467440737.0955161500000000": parseint.ErrSyntax,
	"1844674407370.955161500000000": parseint.ErrSyntax,
	"18446744073709.55161500000000": parseint.ErrSyntax,
	"184467440737095.5161500000000": parseint.ErrSyntax,
	"1844674407370955.161500000000": parseint.ErrSyntax,
	"18446744073709551.61500000000": parseint.ErrSyntax,
	"184467440737095516.1500000000": parseint.ErrSyntax,
	"1844674407370955161.500000000": parseint.ErrSyntax,
	"18446744073709551615.00000000": parseint.ErrSyntax,

	// Overflow
	"18446744073709551616":              parseint.ErrOverflow,
	"184467440737095516150":             parseint.ErrOverflow,
	"118446744073709551615":             parseint.ErrOverflow,
	"999999999999999999999999":          parseint.ErrOverflow,
	"123123123123123123123123123123123": parseint.ErrOverflow,
}

func TestBase10Uint64(t *testing.T) {
	callBase10Uint64 := func(input string, fn func(any, error)) {
		fn(parseint.Base10Uint64(input))
		fn(parseint.Base10Uint64([]byte(input)))
	}

	requireOK := func(t *testing.T, expect uint64, input string) {
		callBase10Uint64(input, func(actual any, err error) {
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
		for input, expect := range validBase10Uint64 {
			requireOK(t, expect, input)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for input, expectedErr := range invalidBase10Uint64 {
			callBase10Uint64(input, func(a any, err error) {
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
		mid := uint64(math.MaxUint64) / 2
		for i := mid; i <= mid+10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			requireOK(t, i, dec)
		}
	})

	t.Run("range_last10k", func(t *testing.T) {
		max := uint64(math.MaxUint64)
		for i := max; i <= max-10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			requireOK(t, i, dec)
		}
	})

	t.Run("err_overflow", func(t *testing.T) {
		maxUint64 := new(big.Int).SetUint64(math.MaxUint64)
		start := new(big.Int).Add(maxUint64, big.NewInt(1))
		end := new(big.Int).Add(maxUint64, big.NewInt(10_000))
		delta := big.NewInt(1)
		for i := new(big.Int).Set(start); i.Cmp(end) <= 0; i.Add(i, delta) {
			str := i.String()
			callBase10Uint64(str, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})
}

func FuzzBase10Uint64(f *testing.F) {
	for input := range validBase10Uint64 {
		f.Add(input)
	}
	for input := range invalidBase10Uint64 {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base10Uint64(s)
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

func BenchmarkBase10Uint64(b *testing.B) {
	var fn func(string) (uint64, error)
	var fnBytes func([]byte) (uint64, error)
	switch *fBenchmarkFn {
	case BenchmarkFnStrconv:
		fn = func(s string) (uint64, error) {
			return strconv.ParseUint(s, 10, 64)
		}
		fnBytes = func(s []byte) (uint64, error) {
			return strconv.ParseUint(string(s), 10, 64)
		}
	case BenchmarkFnParseint:
		fn = parseint.Base10Uint64[string]
		fnBytes = parseint.Base10Uint64[[]byte]
	default:
		b.Fatalf("unknown benchmark function: %q", *fBenchmarkFn)
	}

	var a uint64
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"min", "0"},
		{"small_3", "987"},
		{"small_4", "9871"},
		{"l13", "9876543210123"},
		{"max", "18446744073709551615"},
		{"syntax", "0.000000000000001"},
		{"overflow", "18446744073709551616"},
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
