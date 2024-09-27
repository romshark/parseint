// Package parseint provides very efficient generic implementations of integer parsers.
package parseint

import "errors"

var (
	ErrSyntax   = errors.New("syntax error")
	ErrOverflow = errors.New("overflow")
)

// Base16Uint16 parses s as a base-16 (hexadecimal) unsigned 16-bit number.
// If any error is encountered ok=false is returned.
// Base16Uint16 is comparable to strconv.ParseUint(s, 16, 16) but more efficient.
// ErrSyntax is returned in any error case. ErrOverflow will never be returned
// because it would cost extra to determine overflow errors and this computation
// would be wasted in most cases where we don't care what kind of error there was.
func Base16Uint16[S string | []byte, U ~uint64 | ~uint32 | ~uint16](s S) (U, error) {
	switch len(s) {
	case 1:
		switch c := s[0]; {
		case c >= '0' && c <= '9':
			return U(c - '0'), nil
		case c >= 'a' && c <= 'f':
			return U(c-'a') + 10, nil
		case c >= 'A' && c <= 'F':
			return U(c-'A') + 10, nil
		}
		return 0, ErrSyntax
	case 2:
		v1, v2 := uint16(lutHex[s[0]]), uint16(lutHex[s[1]])
		if v1|v2 == uint16base16InvalidByte {
			return 0, ErrSyntax
		}
		return U((v1 << 4) | v2), nil
	case 3:
		v1, v2, v3 := uint16(lutHex[s[0]]), uint16(lutHex[s[1]]), uint16(lutHex[s[2]])
		if v1|v2|v3 == uint16base16InvalidByte {
			return 0, ErrSyntax
		}
		return U((v1 << 8) | (v2 << 4) | v3), nil
	case 4:
		v1 := uint16(lutHex[s[0]])
		v2 := uint16(lutHex[s[1]])
		v3 := uint16(lutHex[s[2]])
		v4 := uint16(lutHex[s[3]])
		if v1|v2|v3|v4 == uint16base16InvalidByte {
			return 0, ErrSyntax
		}
		return U((v1 << 12) | (v2 << 8) | (v3 << 4) | v4), nil
	}
	return 0, ErrSyntax // Invalid or overflow
}

// uint16base16InvalidByte is used in lutHex to mark invalid characters.
const uint16base16InvalidByte = 0xff

// lutHex is a lookup table mapping hex characters to their respective base-10 value.
// All other bytes are mapped to uint16base16InvalidCharacter.
var lutHex = [256]uint8{}

func init() {
	for i := range lutHex {
		lutHex[i] = uint16base16InvalidByte
	}
	lutHex['0'] = 0
	lutHex['1'] = 1
	lutHex['2'] = 2
	lutHex['3'] = 3
	lutHex['4'] = 4
	lutHex['5'] = 5
	lutHex['6'] = 6
	lutHex['7'] = 7
	lutHex['8'] = 8
	lutHex['9'] = 9
	lutHex['a'] = 10
	lutHex['b'] = 11
	lutHex['c'] = 12
	lutHex['d'] = 13
	lutHex['e'] = 14
	lutHex['f'] = 15
	lutHex['A'] = 10
	lutHex['B'] = 11
	lutHex['C'] = 12
	lutHex['D'] = 13
	lutHex['E'] = 14
	lutHex['F'] = 15
}

// Base10Uint32 parses s as a base-10 unsigned 32-bit integer.
// Returns ok=false if s contains an invalid character or overflows a uint32.
// Base10Uint32 is comparable to strconv.ParseUint(s, 10, 32) but more more efficient.
func Base10Uint32[S string | []byte, U ~uint64 | ~uint32](s S) (U, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	const maxValue = uint64(1)<<uint(32) - 1

	var n uint64
	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return 0, ErrSyntax
		}
		d := c - '0'

		n *= uint64(10)

		n1 := n + uint64(d)
		if n1 < n || n1 > maxValue {
			// n+d overflows
			return 0, ErrOverflow
		}
		n = n1
	}

	return U(n), nil
}
