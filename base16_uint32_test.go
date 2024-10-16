package parseint_test

import (
	"math"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/romshark/parseint"
	"github.com/stretchr/testify/require"
)

var validBase16Uint32 = map[string]uint32{
	"0":                                0,
	"0000":                             0,
	"1":                                0x1,
	"12":                               0x12,
	"123":                              0x123,
	"1234":                             0x1234,
	"12345":                            0x12345,
	"123456":                           0x123456,
	"1234567":                          0x1234567,
	"12345678":                         0x12345678,
	"0000000012345678":                 0x12345678,
	"9999":                             0x9999,
	"ffffffff":                         math.MaxUint32,
	"FFFFFFFF":                         math.MaxUint32,
	"fFfFfFfF":                         math.MaxUint32,
	"a":                                0xa,
	"af":                               0xaf,
	"A":                                0xa,
	"AF":                               0xaf,
	"aaaa":                             0xaaaa,
	"AAAA":                             0xaaaa,
	"1a2b":                             0x1a2b,
	"1A2B":                             0x1a2b,
	"ffff":                             0xffff,
	"FFFF":                             0xffff,
	"abcd":                             0xabcd,
	"aBcD":                             0xabcd,
	"eeFF":                             0xeeff,
	"0000000000000000000ffffffff":      math.MaxUint32,
	"0000000000000000000000000000a":    0xa,
	"0000000000000000000000000000eeff": 0xeeff,
	"0000000000000000000000000000EEFF": 0xeeff,
}

var invalidBase16Uint32 = map[string]error{
	"":    parseint.ErrSyntax,
	" ":   parseint.ErrSyntax,
	" 1":  parseint.ErrSyntax,
	" a":  parseint.ErrSyntax,
	"-":   parseint.ErrSyntax,
	"-0":  parseint.ErrSyntax,
	"0x":  parseint.ErrSyntax,
	"xff": parseint.ErrSyntax,
	"fxf": parseint.ErrSyntax,

	"x":        parseint.ErrSyntax,
	"fx":       parseint.ErrSyntax,
	"ffx":      parseint.ErrSyntax,
	"fffx":     parseint.ErrSyntax,
	"ffffx":    parseint.ErrSyntax,
	"fffffx":   parseint.ErrSyntax,
	"ffffffx":  parseint.ErrSyntax,
	"fffffffx": parseint.ErrSyntax,

	"00000000x":        parseint.ErrSyntax,
	"00000000fx":       parseint.ErrSyntax,
	"00000000ffx":      parseint.ErrSyntax,
	"00000000fffx":     parseint.ErrSyntax,
	"00000000ffffx":    parseint.ErrSyntax,
	"00000000fffffx":   parseint.ErrSyntax,
	"00000000ffffffx":  parseint.ErrSyntax,
	"00000000fffffffx": parseint.ErrSyntax,

	"xfffffff": parseint.ErrSyntax,
	"fxffffff": parseint.ErrSyntax,
	"ffxfffff": parseint.ErrSyntax,
	"fffxffff": parseint.ErrSyntax,
	"ffffxfff": parseint.ErrSyntax,
	"fffffxff": parseint.ErrSyntax,
	"ffffffxf": parseint.ErrSyntax,

	"xfffffff00000000": parseint.ErrSyntax,
	"fxffffff00000000": parseint.ErrSyntax,
	"ffxfffff00000000": parseint.ErrSyntax,
	"fffxffff00000000": parseint.ErrSyntax,
	"ffffxfff00000000": parseint.ErrSyntax,
	"fffffxff00000000": parseint.ErrSyntax,
	"ffffffxf00000000": parseint.ErrSyntax,
	"fffffffx00000000": parseint.ErrSyntax,

	"ж":  parseint.ErrSyntax,
	"🙂":  parseint.ErrSyntax,
	"🗿":  parseint.ErrSyntax,
	"♻︎": parseint.ErrSyntax,

	// Base16Uint32 returns ErrSyntax even for overflow errors, see documentation.
	"ffffffff1":  parseint.ErrSyntax,
	"FFFFFFFF1":  parseint.ErrSyntax,
	"FFFFFFFFf":  parseint.ErrSyntax,
	"FFFFFFFFff": parseint.ErrSyntax,
}

func TestBase16Uint32(t *testing.T) {
	callBase16Uint32 := func(input string, fn func(any, error)) {
		lower, upper := strings.ToLower(input), strings.ToUpper(input)

		fn(parseint.Base16Uint32[string, uint64](lower))
		fn(parseint.Base16Uint32[string, uint64](upper))
		fn(parseint.Base16Uint32[string, uint32](lower))
		fn(parseint.Base16Uint32[string, uint32](upper))
		fn(parseint.Base16Uint32[[]byte, uint64]([]byte(lower)))
		fn(parseint.Base16Uint32[[]byte, uint64]([]byte(upper)))
		fn(parseint.Base16Uint32[[]byte, uint32]([]byte(lower)))
		fn(parseint.Base16Uint32[[]byte, uint32]([]byte(upper)))
	}

	requireOK := func(t *testing.T, expect uint64, input string) {
		callBase16Uint32(input, func(actual any, err error) {
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
		for input, expect := range validBase16Uint32 {
			requireOK(t, uint64(expect), input)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for input, expectedErr := range invalidBase16Uint32 {
			callBase16Uint32(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr, "%q", input)
				require.Zero(t, a)
			})
		}
	})

	t.Run("range_0_10k", func(t *testing.T) {
		for i := uint64(0); i <= 10_000; i++ {
			hex := strconv.FormatUint(i, 16)
			requireOK(t, i, hex)
		}
	})

	t.Run("range_mid_10k", func(t *testing.T) {
		mid := uint64(math.MaxUint32 / 2)
		for i := uint64(mid); i <= mid+10_000; i++ {
			hex := strconv.FormatUint(i, 16)
			requireOK(t, i, hex)
		}
	})

	t.Run("range_last_10k", func(t *testing.T) {
		max := uint64(math.MaxUint32)
		for i := max - 10_000; i <= max; i++ {
			hex := strconv.FormatUint(i, 16)
			requireOK(t, i, hex)
		}
	})

	t.Run("err_overflow", func(t *testing.T) {
		for i := uint64(math.MaxUint32 + 1); i <= math.MaxUint32+10_000; i++ {
			hex := strconv.FormatUint(i, 16)
			callBase16Uint32(hex, func(a any, err error) {
				require.Error(t, err)
				require.Zero(t, a)
			})
		}
	})
}

func fuzzBase16Uint32[U uint64 | uint32](f *testing.F) {
	for input := range validBase16Uint32 {
		f.Add(input)
	}
	for input := range invalidBase16Uint32 {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base16Uint32[string, U](s)
		std, errStd := strconv.ParseUint(s, 16, 32)
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
			if _, err := strconv.ParseUint(s, 16, 32); err == nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			}
		}
	})
}

func FuzzBase16Uint32_uint64(f *testing.F) { fuzzBase16Uint32[uint64](f) }
func FuzzBase16Uint32_uint32(f *testing.F) { fuzzBase16Uint32[uint32](f) }

// BenchmarkBase16Uint32_uint64 compares strconv.ParseUint
// and parseint.Base16Uint32[string, uint64]
func BenchmarkBase16Uint32_uint64(b *testing.B) {
	fn := getBenchmarkFn(b, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 16, 32)
	}, parseint.Base16Uint32[string, uint64])
	fnBytes := getBenchmarkFn(b, func(s []byte) (uint64, error) {
		return strconv.ParseUint(string(s), 16, 32)
	}, parseint.Base16Uint32[[]byte, uint64])

	var a uint64
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"min", "0"},
		{"max_low", "ffffffff"},
		{"max_upp", "FFFFFFFF"},
		{"syntax", "fffx"},
		{"overflow", "FFFFF"},
		{"leadzero31", "0000000000000000000000000000000F"},
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

// BenchmarkBase16Uint32_uint32 compares strconv.ParseUint
// and parseint.Base16Uint32[string, uint16]
func BenchmarkBase16Uint32_uint32(b *testing.B) {
	fn := getBenchmarkFn(b, func(s string) (uint32, error) {
		x, err := strconv.ParseUint(s, 16, 32)
		return uint32(x), err
	}, parseint.Base16Uint32[string, uint32])
	fnBytes := getBenchmarkFn(b, func(s []byte) (uint32, error) {
		x, err := strconv.ParseUint(string(s), 16, 32)
		return uint32(x), err
	}, parseint.Base16Uint32[[]byte, uint32])

	var a uint32
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"min", "0"},
		{"max_low", "ffffffff"},
		{"max_upp", "FFFFFFFF"},
		{"syntax", "fffx"},
		{"overflow", "FFFFF"},
		{"leadzero31", "0000000000000000000000000000000F"},
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
