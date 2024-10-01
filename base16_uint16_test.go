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

var validBase16Uint16 = map[string]uint16{
	"0":                                0,
	"0000":                             0,
	"1":                                0x1,
	"12":                               0x12,
	"123":                              0x123,
	"1234":                             0x1234,
	"9999":                             0x9999,
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
	"0000000000000000000000000000a":    0xa,
	"0000000000000000000000000000eeff": 0xeeff,
	"0000000000000000000000000000EEFF": 0xeeff,
}

var invalidBase16Uint16 = map[string]error{
	"":    parseint.ErrSyntax,
	" ":   parseint.ErrSyntax,
	" 1":  parseint.ErrSyntax,
	" a":  parseint.ErrSyntax,
	"-":   parseint.ErrSyntax,
	"-0":  parseint.ErrSyntax,
	"0x":  parseint.ErrSyntax,
	"xff": parseint.ErrSyntax,
	"fxf": parseint.ErrSyntax,
	"ffx": parseint.ErrSyntax,

	"xfff": parseint.ErrSyntax,
	"fxff": parseint.ErrSyntax,
	"ffxf": parseint.ErrSyntax,
	"fffx": parseint.ErrSyntax,

	"ж":  parseint.ErrSyntax,
	"🙂":  parseint.ErrSyntax,
	"🗿":  parseint.ErrSyntax,
	"♻︎": parseint.ErrSyntax,

	// Base16Uint16 returns ErrSyntax even for overflow errors, see documentation.
	"ffff1":    parseint.ErrSyntax,
	"FFFF1":    parseint.ErrSyntax,
	"FFFFF":    parseint.ErrSyntax,
	"FFFFFFFF": parseint.ErrSyntax,
}

func TestBase16Uint16(t *testing.T) {
	callBase16Uint16 := func(input string, fn func(any, error)) {
		lower, upper := strings.ToLower(input), strings.ToUpper(input)

		fn(parseint.Base16Uint16[string, uint64](lower))
		fn(parseint.Base16Uint16[string, uint64](upper))
		fn(parseint.Base16Uint16[string, uint32](lower))
		fn(parseint.Base16Uint16[string, uint32](upper))
		fn(parseint.Base16Uint16[string, uint16](lower))
		fn(parseint.Base16Uint16[string, uint16](upper))
		fn(parseint.Base16Uint16[[]byte, uint64]([]byte(lower)))
		fn(parseint.Base16Uint16[[]byte, uint64]([]byte(upper)))
		fn(parseint.Base16Uint16[[]byte, uint32]([]byte(lower)))
		fn(parseint.Base16Uint16[[]byte, uint32]([]byte(upper)))
		fn(parseint.Base16Uint16[[]byte, uint16]([]byte(lower)))
		fn(parseint.Base16Uint16[[]byte, uint16]([]byte(upper)))
	}

	requireOK := func(t *testing.T, expect uint64, input string) {
		callBase16Uint16(input, func(actual any, err error) {
			require.NoError(t, err)
			switch actual := actual.(type) {
			case uint64:
				require.Equal(t, uint64(expect), actual)
			case uint32:
				require.Equal(t, uint32(expect), actual)
			case uint16:
				require.Equal(t, uint16(expect), actual)
			default:
				t.Fatalf("unexpected type: %T", actual)
			}
		})
	}

	t.Run("valid", func(t *testing.T) {
		for input, expect := range validBase16Uint16 {
			requireOK(t, uint64(expect), input)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		for input, expectedErr := range invalidBase16Uint16 {
			callBase16Uint16(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr, "%q", input)
				require.Zero(t, a)
			})
		}
	})

	t.Run("range_valid", func(t *testing.T) {
		// Iterating over 65535 values is relatively cheap.
		for i := uint64(0); i <= math.MaxUint16; i++ {
			hex := strconv.FormatUint(i, 16)
			requireOK(t, i, hex)
		}
	})

	t.Run("err_overflow", func(t *testing.T) {
		for i := uint64(math.MaxUint16 + 1); i <= math.MaxUint16+10_000; i++ {
			hex := strconv.FormatUint(i, 16)
			callBase16Uint16(hex, func(a any, err error) {
				require.Error(t, err)
				require.Zero(t, a)
			})
		}
	})
}

func fuzzBase16Uint16[U uint64 | uint32 | uint16](f *testing.F) {
	for input := range validBase16Uint16 {
		f.Add(input)
	}
	for input := range invalidBase16Uint16 {
		f.Add(input)
	}

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base16Uint16[string, U](s)
		std, errStd := strconv.ParseUint(s, 16, 16)
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
			if _, err := strconv.ParseUint(s, 16, 16); err == nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			}
		}
	})
}

func FuzzBase16Uint16_uint64(f *testing.F) { fuzzBase16Uint16[uint64](f) }
func FuzzBase16Uint16_uint32(f *testing.F) { fuzzBase16Uint16[uint32](f) }
func FuzzBase16Uint16_uint16(f *testing.F) { fuzzBase16Uint16[uint16](f) }

// BenchmarkBase16Uint16_uint64 compares strconv.ParseUint
// and parseint.Base16Uint16[string, uint64]
func BenchmarkBase16Uint16_uint64(b *testing.B) {
	var fn func(string) (uint64, error)
	var fnBytes func([]byte) (uint64, error)
	switch *fBenchmarkFn {
	case BenchmarkFnStrconv:
		fn = func(s string) (uint64, error) {
			return strconv.ParseUint(s, 16, 16)
		}
		fnBytes = func(s []byte) (uint64, error) {
			return strconv.ParseUint(string(s), 16, 16)
		}
	case BenchmarkFnParseint:
		fn = parseint.Base16Uint16[string, uint64]
		fnBytes = parseint.Base16Uint16[[]byte, uint64]
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
		{"max_low", "ffff"},
		{"max_upp", "FFFF"},
		{"syntax", "fffx"},
		{"overflow", "FFFFF"},
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

// BenchmarkBase16Uint16_uint16 compares strconv.ParseUint
// and parseint.Base16Uint16[string, uint16]
func BenchmarkBase16Uint16_uint16(b *testing.B) {
	var fn func(string) (uint16, error)
	var fnBytes func([]byte) (uint16, error)
	switch *fBenchmarkFn {
	case BenchmarkFnStrconv:
		fn = func(s string) (uint16, error) {
			x, err := strconv.ParseUint(s, 16, 16)
			return uint16(x), err
		}
		fnBytes = func(s []byte) (uint16, error) {
			x, err := strconv.ParseUint(string(s), 16, 16)
			return uint16(x), err
		}
	case BenchmarkFnParseint:
		fn = parseint.Base16Uint16[string, uint16]
		fnBytes = parseint.Base16Uint16[[]byte, uint16]
	default:
		b.Fatalf("unknown benchmark function: %q", *fBenchmarkFn)
	}

	var a uint16
	var err error
	for _, td := range []struct {
		name  string
		input string
	}{
		{"min", "0"},
		{"max_low", "ffff"},
		{"max_upp", "FFFF"},
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
