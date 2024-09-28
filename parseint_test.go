package parseint_test

import (
	"flag"
	"math"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/romshark/parseint"
	"github.com/stretchr/testify/require"
)

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

	t.Run("range_valid", func(t *testing.T) {
		// Iterating over 65535 values is relatively cheap.
		for i := uint64(0); i <= math.MaxUint16; i++ {
			hex := strconv.FormatUint(i, 16)
			callBase16Uint16(hex, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				case uint32:
					require.Equal(t, uint32(i), a)
				case uint16:
					require.Equal(t, uint16(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("err_syntax", func(t *testing.T) {
		f := func(t *testing.T, input string) {
			t.Helper()
			callBase16Uint16(input, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrSyntax)
				require.Zero(t, a)
			})
		}
		f(t, "")
		f(t, "-")
		f(t, "-0")
		f(t, "0x")
		f(t, "0x")
		f(t, "xff")
		f(t, "fxf")
		f(t, "ffx")
		f(t, "xfff")
		f(t, "fxff")
		f(t, "ffxf")
		f(t, "fffx")
		f(t, "ж")
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

func TestBase10Uint32(t *testing.T) {
	callBase10Uint32 := func(input string, fn func(any, error)) {
		fn(parseint.Base10Uint32[string, uint64](input))
		fn(parseint.Base10Uint32[string, uint32](input))
		fn(parseint.Base10Uint32[[]byte, uint64]([]byte(input)))
		fn(parseint.Base10Uint32[[]byte, uint32]([]byte(input)))
	}

	t.Run("range_0_10k", func(t *testing.T) {
		for i := uint64(0); i <= 10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint32(dec, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				case uint32:
					require.Equal(t, uint32(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("range_mid10k", func(t *testing.T) {
		mid := uint64(math.MaxUint32) / 2
		for i := mid; i <= mid+10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint32(dec, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				case uint32:
					require.Equal(t, uint32(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("range_last10k", func(t *testing.T) {
		max := uint64(math.MaxUint32)
		for i := max; i <= max-10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint32(dec, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				case uint32:
					require.Equal(t, uint32(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("err_syntax", func(t *testing.T) {
		f := func(t *testing.T, input string, expectedErr error) {
			t.Helper()
			callBase10Uint32(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr)
				require.Zero(t, a)
			})
		}

		// Invalid input
		f(t, "", parseint.ErrSyntax)
		f(t, "-", parseint.ErrSyntax)
		f(t, "-0", parseint.ErrSyntax)
		f(t, "-x", parseint.ErrSyntax)
		f(t, "0x0", parseint.ErrSyntax)
		f(t, "123x", parseint.ErrSyntax)
		f(t, "ж", parseint.ErrSyntax)

		// Overflow
		f(t, "4294967296", parseint.ErrOverflow)
		f(t, "42949672960", parseint.ErrOverflow)
		f(t, "14294967296", parseint.ErrOverflow)
		f(t, "12331232123123123", parseint.ErrOverflow)
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

func TestBase10Uint64(t *testing.T) {
	callBase10Uint64 := func(input string, fn func(any, error)) {
		fn(parseint.Base10Uint64(input))
		fn(parseint.Base10Uint64([]byte(input)))
	}

	t.Run("range_0_10k", func(t *testing.T) {
		for i := uint64(0); i <= 10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint64(dec, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("range_mid10k", func(t *testing.T) {
		mid := uint64(math.MaxUint64) / 2
		for i := mid; i <= mid+10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint64(dec, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("range_last10k", func(t *testing.T) {
		max := uint64(math.MaxUint64)
		for i := max; i <= max-10_000; i++ {
			dec := strconv.FormatUint(i, 10)
			callBase10Uint64(dec, func(a any, err error) {
				require.NoError(t, err)
				switch a := a.(type) {
				case uint64:
					require.Equal(t, uint64(i), a)
				default:
					t.Fatalf("unexpected type: %T", a)
				}
			})
		}
	})

	t.Run("err_syntax", func(t *testing.T) {
		f := func(t *testing.T, input string, expectedErr error) {
			t.Helper()
			callBase10Uint64(input, func(a any, err error) {
				require.ErrorIs(t, err, expectedErr)
				require.Zero(t, a)
			})
		}

		// Invalid input
		f(t, "", parseint.ErrSyntax)
		f(t, "-", parseint.ErrSyntax)
		f(t, "-0", parseint.ErrSyntax)
		f(t, "0.0", parseint.ErrSyntax)
		f(t, "-x", parseint.ErrSyntax)
		f(t, "0x0", parseint.ErrSyntax)
		f(t, "123x", parseint.ErrSyntax)
		f(t, "ж", parseint.ErrSyntax)
		f(t, ".18446744073709551615", parseint.ErrSyntax)
		f(t, "1.8446744073709551615", parseint.ErrSyntax)
		f(t, "18.446744073709551615", parseint.ErrSyntax)
		f(t, "184.46744073709551615", parseint.ErrSyntax)
		f(t, "1844.6744073709551615", parseint.ErrSyntax)
		f(t, "18446.744073709551615", parseint.ErrSyntax)
		f(t, "184467.44073709551615", parseint.ErrSyntax)
		f(t, "1844674.4073709551615", parseint.ErrSyntax)
		f(t, "18446744.073709551615", parseint.ErrSyntax)
		f(t, "184467440.73709551615", parseint.ErrSyntax)
		f(t, "1844674407.3709551615", parseint.ErrSyntax)
		f(t, "18446744073.709551615", parseint.ErrSyntax)
		f(t, "184467440737.09551615", parseint.ErrSyntax)
		f(t, "1844674407370.9551615", parseint.ErrSyntax)
		f(t, "18446744073709.551615", parseint.ErrSyntax)
		f(t, "184467440737095.51615", parseint.ErrSyntax)
		f(t, "1844674407370955.1615", parseint.ErrSyntax)
		f(t, "18446744073709551.615", parseint.ErrSyntax)
		f(t, "184467440737095516.15", parseint.ErrSyntax)
		f(t, "1844674407370955161.5", parseint.ErrSyntax)
		f(t, "18446744073709551615.", parseint.ErrSyntax)

		f(t, ".1844674407370955161500000000", parseint.ErrSyntax)
		f(t, "1.844674407370955161500000000", parseint.ErrSyntax)
		f(t, "18.44674407370955161500000000", parseint.ErrSyntax)
		f(t, "184.4674407370955161500000000", parseint.ErrSyntax)
		f(t, "1844.674407370955161500000000", parseint.ErrSyntax)
		f(t, "18446.74407370955161500000000", parseint.ErrSyntax)
		f(t, "184467.4407370955161500000000", parseint.ErrSyntax)
		f(t, "1844674.407370955161500000000", parseint.ErrSyntax)
		f(t, "18446744.07370955161500000000", parseint.ErrSyntax)
		f(t, "184467440.7370955161500000000", parseint.ErrSyntax)
		f(t, "1844674407.370955161500000000", parseint.ErrSyntax)
		f(t, "18446744073.70955161500000000", parseint.ErrSyntax)
		f(t, "184467440737.0955161500000000", parseint.ErrSyntax)
		f(t, "1844674407370.955161500000000", parseint.ErrSyntax)
		f(t, "18446744073709.55161500000000", parseint.ErrSyntax)
		f(t, "184467440737095.5161500000000", parseint.ErrSyntax)
		f(t, "1844674407370955.161500000000", parseint.ErrSyntax)
		f(t, "18446744073709551.61500000000", parseint.ErrSyntax)
		f(t, "184467440737095516.1500000000", parseint.ErrSyntax)
		f(t, "1844674407370955161.500000000", parseint.ErrSyntax)
		f(t, "18446744073709551615.00000000", parseint.ErrSyntax)

		// Overflow
		f(t, "18446744073709551616", parseint.ErrOverflow)
		f(t, "184467440737095516150", parseint.ErrOverflow)
		f(t, "118446744073709551615", parseint.ErrOverflow)
		f(t, "999999999999999999999999", parseint.ErrOverflow)
		f(t, "123123123123123123123123123123123", parseint.ErrOverflow)
	})

	t.Run("err_overflow", func(t *testing.T) {
		maxUint64 := new(big.Int).SetUint64(math.MaxUint64)
		start := new(big.Int).Add(maxUint64, big.NewInt(1))
		end := new(big.Int).Add(maxUint64, big.NewInt(10_000))
		for i := new(big.Int).Set(start); i.Cmp(end) <= 0; i.Add(i, big.NewInt(1)) {
			str := i.String()
			callBase10Uint64(str, func(a any, err error) {
				require.ErrorIs(t, err, parseint.ErrOverflow)
				require.Zero(t, a)
			})
		}
	})
}

func fuzzBase16Uint16[U uint64 | uint32 | uint16](f *testing.F) {
	// Valid inputs.
	f.Add("0")
	f.Add("1")
	f.Add("12")
	f.Add("123")
	f.Add("1234")
	f.Add("9999")
	f.Add("a")
	f.Add("af")
	f.Add("aaaa")
	f.Add("1a2b")
	f.Add("FFFF")
	f.Add("abcd")
	f.Add("eeFF")

	// Invalid inputs.
	f.Add("11111") // Overflow
	f.Add("")
	f.Add("defg")
	f.Add("123456789") // Overflow
	f.Add("xyz")
	f.Add("GHIJ")
	f.Add("🙂")
	f.Add("🗿")
	f.Add("♻︎")

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base16Uint16[string, U](s)
		if err == nil && x > math.MaxUint16 {
			t.Errorf("%q: returned value out of 16-bit range: %x", s, x)
		} else if err != nil && x != 0 {
			t.Errorf("%q: failed but returned non-zero value: %x", s, x)
		}
	})
}

func FuzzBase16Uint16_uint64(f *testing.F) { fuzzBase16Uint16[uint64](f) }
func FuzzBase16Uint16_uint32(f *testing.F) { fuzzBase16Uint16[uint32](f) }
func FuzzBase16Uint16_uint16(f *testing.F) { fuzzBase16Uint16[uint16](f) }

func fuzzBase10Uint32[U uint64 | uint32](f *testing.F) {
	// Valid inputs.
	f.Add("0")
	f.Add("1")
	f.Add("12")
	f.Add("123")
	f.Add("1234")
	f.Add("999999999")
	f.Add("1234567890")
	f.Add("4294967294")
	f.Add("4294967295") // Max

	// Invalid inputs.
	f.Add("")
	f.Add("a")
	f.Add("af")
	f.Add("aaaa")
	f.Add("1a2b")
	f.Add("FFFF")
	f.Add("abcd")
	f.Add("eeFF")
	f.Add("defg")
	f.Add("xyz")
	f.Add("GHIJ")
	f.Add("🙂")
	f.Add("🗿")
	f.Add("♻︎")

	// Overflow
	f.Add("4294967296")
	f.Add("9999999999")
	f.Add("123456789123456789")

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base10Uint32[string, U](s)
		if err == nil {
			std, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			} else if std != uint64(x) {
				t.Errorf("expected %d; received: %d", std, uint64(x))
			}
		} else if x != 0 {
			t.Errorf("%q: failed but returned non-zero value: %x", s, x)
		}
	})
}

func FuzzBase10Uint64(f *testing.F) {
	// Valid inputs.
	f.Add("0")
	f.Add("1")
	f.Add("10")
	f.Add("100")
	f.Add("1000")
	f.Add("10000")
	f.Add("100000")
	f.Add("1000000")
	f.Add("10000000")
	f.Add("100000000")
	f.Add("1000000000")
	f.Add("10000000000")
	f.Add("100000000000")
	f.Add("1000000000000")
	f.Add("10000000000000")
	f.Add("100000000000000")
	f.Add("1000000000000000")
	f.Add("10000000000000000")
	f.Add("100000000000000000")
	f.Add("1000000000000000000")
	f.Add("10000000000000000000")
	f.Add("12")
	f.Add("123")
	f.Add("1234")
	f.Add("999999999")
	f.Add("1234567890")
	f.Add("4294967294")
	f.Add("184")
	f.Add("1844")
	f.Add("18446")
	f.Add("184467")
	f.Add("1844674")
	f.Add("18446744")
	f.Add("184467440")
	f.Add("1844674407")
	f.Add("18446744073")
	f.Add("184467440737")
	f.Add("1844674407370")
	f.Add("18446744073709")
	f.Add("184467440737095")
	f.Add("1844674407370955")
	f.Add("18446744073709551")
	f.Add("184467440737095516")
	f.Add("1844674407370955161")
	f.Add("18446744073709551615") // Max

	// Invalid input
	f.Add("")
	f.Add("-")
	f.Add("-0")
	f.Add("0.0")
	f.Add("-x")
	f.Add("0x0")
	f.Add("123x")
	f.Add("ж")
	f.Add(".18446744073709551615")
	f.Add("1.8446744073709551615")
	f.Add("18.446744073709551615")
	f.Add("184.46744073709551615")
	f.Add("1844.6744073709551615")
	f.Add("18446.744073709551615")
	f.Add("184467.44073709551615")
	f.Add("1844674.4073709551615")
	f.Add("18446744.073709551615")
	f.Add("184467440.73709551615")
	f.Add("1844674407.3709551615")
	f.Add("18446744073.709551615")
	f.Add("184467440737.09551615")
	f.Add("1844674407370.9551615")
	f.Add("18446744073709.551615")
	f.Add("184467440737095.51615")
	f.Add("1844674407370955.1615")
	f.Add("18446744073709551.615")
	f.Add("184467440737095516.15")
	f.Add("1844674407370955161.5")
	f.Add("18446744073709551615.")

	f.Add(".1844674407370955161500000000")
	f.Add("1.844674407370955161500000000")
	f.Add("18.44674407370955161500000000")
	f.Add("184.4674407370955161500000000")
	f.Add("1844.674407370955161500000000")
	f.Add("18446.74407370955161500000000")
	f.Add("184467.4407370955161500000000")
	f.Add("1844674.407370955161500000000")
	f.Add("18446744.07370955161500000000")
	f.Add("184467440.7370955161500000000")
	f.Add("1844674407.370955161500000000")
	f.Add("18446744073.70955161500000000")
	f.Add("184467440737.0955161500000000")
	f.Add("1844674407370.955161500000000")
	f.Add("18446744073709.55161500000000")
	f.Add("184467440737095.5161500000000")
	f.Add("1844674407370955.161500000000")
	f.Add("18446744073709551.61500000000")
	f.Add("184467440737095516.1500000000")
	f.Add("1844674407370955161.500000000")
	f.Add("18446744073709551615.00000000")

	// Overflow
	f.Add("18446744073709551616")
	f.Add("184467440737095516150")
	f.Add("118446744073709551615")
	f.Add("999999999999999999999999")
	f.Add("123123123123123123123123123123123")

	f.Fuzz(func(t *testing.T, s string) {
		x, err := parseint.Base10Uint64[string](s)
		if err == nil {
			std, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", s, err)
			} else if std != x {
				t.Errorf("expected %d; received: %d", std, x)
			}
		} else if x != 0 {
			t.Errorf("%q: failed but returned non-zero value: %x", s, x)
		}
	})
}

func FuzzBase10Uint32_uint64(f *testing.F) { fuzzBase10Uint32[uint64](f) }
func FuzzBase10Uint32_uint32(f *testing.F) { fuzzBase10Uint32[uint32](f) }

var fBenchmarkFn = flag.String(
	"benchfunc",
	BenchmarkFnParseint,
	`function to benchmark, use either "strconv" or "parseint"`)

const (
	BenchmarkFnStrconv  = "strconv"
	BenchmarkFnParseint = "parseint"
)

// BenchmarkBase16Uint16 compares strconv.ParseUint
// and parseint.Base16Uint16[string, uint64]
func BenchmarkBase16Uint16(b *testing.B) {
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

func BenchmarkBase10Uint32(b *testing.B) {
	var fn func(string) (uint64, error)
	var fnBytes func([]byte) (uint64, error)
	switch *fBenchmarkFn {
	case BenchmarkFnStrconv:
		fn = func(s string) (uint64, error) {
			return strconv.ParseUint(s, 10, 32)
		}
		fnBytes = func(s []byte) (uint64, error) {
			return strconv.ParseUint(string(s), 10, 32)
		}
	case BenchmarkFnParseint:
		fn = parseint.Base10Uint32[string, uint64]
		fnBytes = parseint.Base10Uint32[[]byte, uint64]
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
		{"max", "4294967295"},
		{"syntax", "-"},
		{"overflow", "99999999999"},
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
		{"max", "18446744073709551615"},
		{"syntax", "0.000000000000001"},
		{"overflow", "18446744073709551616"},
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
