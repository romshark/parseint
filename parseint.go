// Package parseint provides very efficient generic implementations of integer parsers.
package parseint

import "errors"

var (
	ErrSyntax   = errors.New("syntax error")
	ErrOverflow = errors.New("overflow")
)

// Base16Uint16 parses s as a base-16 (hexadecimal) unsigned 16-bit integer.
// ErrSyntax is returned in any error case. ErrOverflow will never be returned
// because it would cost extra to determine overflow errors and this computation
// would be wasted in most cases where we don't care what kind of error there was.
// Base16Uint16 is comparable to strconv.ParseUint(s, 16, 16) but is more efficient.
func Base16Uint16[S string | []byte, U ~uint64 | ~uint32 | ~uint16](s S) (U, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	if s[0] == '0' { // Skip all leading zeroes if any
		var i int
		var c byte
		for i, c = range []byte(s) {
			if c != '0' {
				s = s[i:]
				goto INT
			}
		}
		if i == len(s)-1 {
			return 0, nil // Input consists exclusively of zeroes.
		}
	}
INT:
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
		if v1|v2 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 4) | v2), nil
	case 3:
		v1, v2, v3 := uint16(lutHex[s[0]]), uint16(lutHex[s[1]]), uint16(lutHex[s[2]])
		if v1|v2|v3 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 8) | (v2 << 4) | v3), nil
	case 4:
		v1 := uint16(lutHex[s[0]])
		v2 := uint16(lutHex[s[1]])
		v3 := uint16(lutHex[s[2]])
		v4 := uint16(lutHex[s[3]])
		if v1|v2|v3|v4 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 12) | (v2 << 8) | (v3 << 4) | v4), nil
	}
	return 0, ErrSyntax // Invalid or overflow
}

// Base16Uint32 parses s as a base-16 (hexadecimal) unsigned 32-bit integer.
// ErrSyntax is returned in any error case. ErrOverflow will never be returned
// because it would cost extra to determine overflow errors and this computation
// would be wasted in most cases where we don't care what kind of error there was.
// Base16Uint32 is comparable to strconv.ParseUint(s, 16, 32) but is more efficient.
func Base16Uint32[S string | []byte, U ~uint64 | ~uint32](s S) (U, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	if s[0] == '0' { // Skip all leading zeroes if any
		var i int
		var c byte
		for i, c = range []byte(s) {
			if c != '0' {
				s = s[i:]
				goto INT
			}
		}
		if i == len(s)-1 {
			return 0, nil // Input consists exclusively of zeroes.
		}
	}
INT:
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
		if v1|v2 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 4) | v2), nil
	case 3:
		v1, v2, v3 := uint16(lutHex[s[0]]), uint16(lutHex[s[1]]), uint16(lutHex[s[2]])
		if v1|v2|v3 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 8) | (v2 << 4) | v3), nil
	case 4:
		v1 := uint32(lutHex[s[0]])
		v2 := uint32(lutHex[s[1]])
		v3 := uint32(lutHex[s[2]])
		v4 := uint32(lutHex[s[3]])
		if v1|v2|v3|v4 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 12) | (v2 << 8) | (v3 << 4) | v4), nil
	case 5:
		v1 := uint32(lutHex[s[0]])
		v2 := uint32(lutHex[s[1]])
		v3 := uint32(lutHex[s[2]])
		v4 := uint32(lutHex[s[3]])
		v5 := uint32(lutHex[s[4]])
		if v1|v2|v3|v4|v5 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 16) | (v2 << 12) | (v3 << 8) | (v4 << 4) | v5), nil
	case 6:
		v1 := uint32(lutHex[s[0]])
		v2 := uint32(lutHex[s[1]])
		v3 := uint32(lutHex[s[2]])
		v4 := uint32(lutHex[s[3]])
		v5 := uint32(lutHex[s[4]])
		v6 := uint32(lutHex[s[5]])
		if v1|v2|v3|v4|v5|v6 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 20) | (v2 << 16) | (v3 << 12) |
			(v4 << 8) | (v5 << 4) | v6), nil
	case 7:
		v1 := uint32(lutHex[s[0]])
		v2 := uint32(lutHex[s[1]])
		v3 := uint32(lutHex[s[2]])
		v4 := uint32(lutHex[s[3]])
		v5 := uint32(lutHex[s[4]])
		v6 := uint32(lutHex[s[5]])
		v7 := uint32(lutHex[s[6]])
		if v1|v2|v3|v4|v5|v6|v7 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 24) | (v2 << 20) | (v3 << 16) | (v4 << 12) |
			(v5 << 8) | (v6 << 4) | v7), nil
	case 8:
		v1 := uint32(lutHex[s[0]])
		v2 := uint32(lutHex[s[1]])
		v3 := uint32(lutHex[s[2]])
		v4 := uint32(lutHex[s[3]])
		v5 := uint32(lutHex[s[4]])
		v6 := uint32(lutHex[s[5]])
		v7 := uint32(lutHex[s[6]])
		v8 := uint32(lutHex[s[7]])
		if v1|v2|v3|v4|v5|v6|v7|v8 == invalidHexByte {
			return 0, ErrSyntax
		}
		return U((v1 << 28) | (v2 << 24) | (v3 << 20) | (v4 << 16) | (v5 << 12) |
			(v6 << 8) | (v7 << 4) | v8), nil
	}
	return 0, ErrSyntax // Invalid or overflow
}

// invalidHexByte is used in lutHex to mark invalid characters.
const invalidHexByte = 0xff

// lutHex is a lookup table mapping hex characters to their respective base-10 value.
// All other bytes are mapped to uint16base16InvalidByte.
var lutHex = [256]uint8{}

func init() {
	for i := range lutHex {
		lutHex[i] = invalidHexByte
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
// Returns ErrSyntax if s contains an invalid character.
// Returns ErrOverflow if the stringified value overflows a uint32.
// Base10Uint32 is comparable to strconv.ParseUint(s, 10, 32) but is more efficient.
func Base10Uint32[S string | []byte, U ~uint64 | ~uint32](s S) (U, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	const max = 1<<32 - 1

	var n uint64
	// A batch of 8 digits here might actually give a geomean of -1.7%
	// but I'm not entirely sure this is worth it.
	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return 0, ErrSyntax
		}
		n *= uint64(10)
		n1 := n + uint64(c-'0')
		if n1 < n || n1 > max {
			return 0, ErrOverflow
		}
		n = n1
	}

	return U(n), nil
}

// Base10Int32 parses s as a base-10 signed 32-bit integer.
// Returns ErrSyntax if s contains an invalid character.
// Returns ErrOverflow if the stringified value overflows an int32.
// Base10Int32 is comparable to strconv.ParseInt(s, 10, 32) but is more efficient.
func Base10Int32[S string | []byte, I ~int64 | ~int32](s S) (I, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	var n int64
	switch s[0] {
	case '-': // Negative integer.
		const max = int64(1 << 31)
		if len(s) == 1 { // Sign without any following digits.
			return 0, ErrSyntax
		}
		for _, c := range []byte(s[1:]) {
			if c < '0' || c > '9' {
				return 0, ErrSyntax
			}
			n *= 10
			n += int64(c - '0')
			if n > max {
				return 0, ErrOverflow
			}
		}
		return I(-n), nil
	case '+':
		if len(s) == 1 { // Sign without any following digits.
			return 0, ErrSyntax
		}
		s = s[1:] // Remove sign.
	}

	// Positive integer.
	const max = int64(1<<31 - 1)
	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return 0, ErrSyntax
		}
		n *= 10
		n += int64(c - '0')
		if n > max {
			return 0, ErrOverflow
		}
	}
	return I(n), nil
}

// Base10Uint64 parses s as a base-10 unsigned 64-bit integer.
// Returns ErrSyntax if s contains an invalid character.
// Returns ErrOverflow if the stringified value overflows a uint64.
// Base10Uint64 is comparable to strconv.ParseUint(s, 10, 64) but is more efficient.
func Base10Uint64[S string | []byte](s S) (uint64, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	const max = 1<<64 - 1

	var n uint64
	for len(s) > 7 { // Process 8 digits at a time as long as possible.
		c0, c1, c2, c3, c4, c5, c6, c7 := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7]
		if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' ||
			c2 < '0' || c2 > '9' || c3 < '0' || c3 > '9' ||
			c4 < '0' || c4 > '9' || c5 < '0' || c5 > '9' ||
			c6 < '0' || c6 > '9' || c7 < '0' || c7 > '9' {
			return 0, ErrSyntax
		}
		d := uint64(c0-'0')*10_000_000 +
			uint64(c1-'0')*1_000_000 +
			uint64(c2-'0')*100_000 +
			uint64(c3-'0')*10_000 +
			uint64(c4-'0')*1_000 +
			uint64(c5-'0')*100 +
			uint64(c6-'0')*10 +
			uint64(c7-'0')
		if n > (max-d)/100_000_000 {
			return 0, ErrOverflow
		}
		n = n*100_000_000 + d
		s = s[8:]
	}
	for len(s) > 3 { // Process 4 digits at a time as long as possible.
		c0, c1, c2, c3 := s[0], s[1], s[2], s[3]
		if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' ||
			c2 < '0' || c2 > '9' || c3 < '0' || c3 > '9' {
			return 0, ErrSyntax
		}
		d := uint64(c0-'0')*1_000 +
			uint64(c1-'0')*100 +
			uint64(c2-'0')*10 +
			uint64(c3-'0')
		if n > (max-d)/10_000 {
			return 0, ErrOverflow
		}
		n = n*10_000 + d
		s = s[4:]
	}
	for _, c := range []byte(s) { // Process remaining digits one at a time.
		if c < '0' || c > '9' {
			return 0, ErrSyntax
		}
		d := uint64(c - '0')
		if n > (max-d)/10 {
			return 0, ErrOverflow
		}
		n = n*10 + d
	}
	return n, nil
}

// Base10Int64 parses s as a base-10 signed 64-bit integer.
// Returns ErrSyntax if s contains an invalid character.
// Returns ErrOverflow if the stringified value overflows an int64.
// Base10Int64 is comparable to strconv.ParseInt(s, 10, 64) but is more efficient.
func Base10Int64[S string | []byte](s S) (int64, error) {
	if len(s) == 0 {
		return 0, ErrSyntax
	}
	var n uint64
	switch s[0] {
	case '-': // Negative integer.
		const max = uint64(1 << 63)
		if len(s) == 1 { // Sign without any following digits.
			return 0, ErrSyntax
		}
		s = s[1:]        // Remove the sign
		for len(s) > 7 { // Process 8 digits at a time as long as possible.
			c0, c1, c2, c3, c4, c5, c6, c7 := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7]
			if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' ||
				c2 < '0' || c2 > '9' || c3 < '0' || c3 > '9' ||
				c4 < '0' || c4 > '9' || c5 < '0' || c5 > '9' ||
				c6 < '0' || c6 > '9' || c7 < '0' || c7 > '9' {
				return 0, ErrSyntax
			}
			d := uint64(c0-'0')*10_000_000 +
				uint64(c1-'0')*1_000_000 +
				uint64(c2-'0')*100_000 +
				uint64(c3-'0')*10_000 +
				uint64(c4-'0')*1_000 +
				uint64(c5-'0')*100 +
				uint64(c6-'0')*10 +
				uint64(c7-'0')
			if n > (max-d)/100_000_000 {
				return 0, ErrOverflow
			}
			n = n*100_000_000 + d
			s = s[8:]
		}
		for len(s) > 3 { // Process 4 digits at a time as long as possible.
			c0, c1, c2, c3 := s[0], s[1], s[2], s[3]
			if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' ||
				c2 < '0' || c2 > '9' || c3 < '0' || c3 > '9' {
				return 0, ErrSyntax
			}
			d := uint64(c0-'0')*1_000 +
				uint64(c1-'0')*100 +
				uint64(c2-'0')*10 +
				uint64(c3-'0')
			if n > (max-d)/10_000 {
				return 0, ErrOverflow
			}
			n = n*10_000 + d
			s = s[4:]
		}
		for _, c := range []byte(s) { // Process remaining digits one at a time.
			if c < '0' || c > '9' {
				return 0, ErrSyntax
			}
			n *= 10
			n += uint64(c - '0')
			if n > max {
				return 0, ErrOverflow
			}
		}
		return int64(-n), nil
	case '+':
		if len(s) == 1 { // Sign without any following digits.
			return 0, ErrSyntax
		}
		s = s[1:] // Remove sign.
	}

	// Positive integer.
	const max = uint64(1<<63 - 1)
	const cutoff = max/10 + 1

	for len(s) > 7 { // Process 8 digits at a time as long as possible.
		c0, c1, c2, c3, c4, c5, c6, c7 := s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7]
		if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' ||
			c2 < '0' || c2 > '9' || c3 < '0' || c3 > '9' ||
			c4 < '0' || c4 > '9' || c5 < '0' || c5 > '9' ||
			c6 < '0' || c6 > '9' || c7 < '0' || c7 > '9' {
			return 0, ErrSyntax
		}
		d := uint64(c0-'0')*10_000_000 +
			uint64(c1-'0')*1_000_000 +
			uint64(c2-'0')*100_000 +
			uint64(c3-'0')*10_000 +
			uint64(c4-'0')*1_000 +
			uint64(c5-'0')*100 +
			uint64(c6-'0')*10 +
			uint64(c7-'0')
		if n > (max-d)/100_000_000 {
			return 0, ErrOverflow
		}
		n = n*100_000_000 + d
		s = s[8:]
	}
	for len(s) > 3 { // Process 4 digits at a time as long as possible.
		c0, c1, c2, c3 := s[0], s[1], s[2], s[3]
		if c0 < '0' || c0 > '9' || c1 < '0' || c1 > '9' ||
			c2 < '0' || c2 > '9' || c3 < '0' || c3 > '9' {
			return 0, ErrSyntax
		}
		d := uint64(c0-'0')*1_000 +
			uint64(c1-'0')*100 +
			uint64(c2-'0')*10 +
			uint64(c3-'0')
		if n > (max-d)/10_000 {
			return 0, ErrOverflow
		}
		n = n*10_000 + d
		s = s[4:]
	}
	for _, c := range []byte(s) { // Process remaining digits one at a time.
		if c < '0' || c > '9' {
			return 0, ErrSyntax
		}
		if n > cutoff {
			return 0, ErrOverflow
		}
		n *= 10
		n += uint64(c - '0')
		if n > max {
			return 0, ErrOverflow
		}
	}
	return int64(n), nil
}
