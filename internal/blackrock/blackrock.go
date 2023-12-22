package blackrock

import (
	"math"
)

var SBOX = []uint64{
	0x91, 0x58, 0xb3, 0x31, 0x6c, 0x33, 0xda, 0x88,
	0x57, 0xdd, 0x8c, 0xf2, 0x29, 0x5a, 0x08, 0x9f,
	0x49, 0x34, 0xce, 0x99, 0x9e, 0xbf, 0x0f, 0x81,
	0xd4, 0x2f, 0x92, 0x3f, 0x95, 0xf5, 0x23, 0x00,
	0x0d, 0x3e, 0xa8, 0x90, 0x98, 0xdd, 0x20, 0x00,
	0x03, 0x69, 0x0a, 0xca, 0xba, 0x12, 0x08, 0x41,
	0x6e, 0xb9, 0x86, 0xe4, 0x50, 0xf0, 0x84, 0xe2,
	0xb3, 0xb3, 0xc8, 0xb5, 0xb2, 0x2d, 0x18, 0x70,
	0x0a, 0xd7, 0x92, 0x90, 0x9e, 0x1e, 0x0c, 0x1f,
	0x08, 0xe8, 0x06, 0xfd, 0x85, 0x2f, 0xaa, 0x5d,
	0xcf, 0xf9, 0xe3, 0x55, 0xb9, 0xfe, 0xa6, 0x7f,
	0x44, 0x3b, 0x4a, 0x4f, 0xc9, 0x2f, 0xd2, 0xd3,
	0x8e, 0xdc, 0xae, 0xba, 0x4f, 0x02, 0xb4, 0x76,
	0xba, 0x64, 0x2d, 0x07, 0x9e, 0x08, 0xec, 0xbd,
	0x52, 0x29, 0x07, 0xbb, 0x9f, 0xb5, 0x58, 0x6f,
	0x07, 0x55, 0xb0, 0x34, 0x74, 0x9f, 0x05, 0xb2,
	0xdf, 0xa9, 0xc6, 0x2a, 0xa3, 0x5d, 0xff, 0x10,
	0x40, 0xb3, 0xb7, 0xb4, 0x63, 0x6e, 0xf4, 0x3e,
	0xee, 0xf6, 0x49, 0x52, 0xe3, 0x11, 0xb3, 0xf1,
	0xfb, 0x60, 0x48, 0xa1, 0xa4, 0x19, 0x7a, 0x2e,
	0x90, 0x28, 0x90, 0x8d, 0x5e, 0x8c, 0x8c, 0xc4,
	0xf2, 0x4a, 0xf6, 0xb2, 0x19, 0x83, 0xea, 0xed,
	0x6d, 0xba, 0xfe, 0xd8, 0xb6, 0xa3, 0x5a, 0xb4,
	0x48, 0xfa, 0xbe, 0x5c, 0x69, 0xac, 0x3c, 0x8f,
	0x63, 0xaf, 0xa4, 0x42, 0x25, 0x50, 0xab, 0x65,
	0x80, 0x65, 0xb9, 0xfb, 0xc7, 0xf2, 0x2d, 0x5c,
	0xe3, 0x4c, 0xa4, 0xa6, 0x8e, 0x07, 0x9c, 0xeb,
	0x41, 0x93, 0x65, 0x44, 0x4a, 0x86, 0xc1, 0xf6,
	0x2c, 0x97, 0xfd, 0xf4, 0x6c, 0xdc, 0xe1, 0xe0,
	0x28, 0xd9, 0x89, 0x7b, 0x09, 0xe2, 0xa0, 0x38,
	0x74, 0x4a, 0xa6, 0x5e, 0xd2, 0xe2, 0x4d, 0xf3,
	0xf4, 0xc6, 0xbc, 0xa2, 0x51, 0x58, 0xe8, 0xae,
}

const (
	DefaultRounds = 14
	MaxRounds     = 64
)

type Blackrock struct {
	rangeLen uint64
	rounds   int
	seed     int64
	a, b     uint64
}

func New(rangeLen uint64, rounds int, seed int64) *Blackrock {
	split := uint64(math.Sqrt(float64(rangeLen)))

	var a, b uint64
	if rangeLen == 0 {
		a, b = 0, 0
	} else if rangeLen == 1 {
		a, b = 1, 1
	} else if rangeLen == 2 {
		a, b = 1, 2
	} else if rangeLen == 3 {
		a, b = 2, 2
	} else if rangeLen >= 4 && rangeLen <= 6 {
		a, b = 2, 3
	} else if rangeLen >= 7 && rangeLen <= 8 {
		a, b = 3, 3
	} else {
		a = split - 2
		b = split + 3
	}

	if rangeLen > 0 {
		for a*b <= rangeLen {
			b++
		}
	}

	return &Blackrock{
		rangeLen: rangeLen,
		rounds:   min(MaxRounds, rounds),
		seed:     seed,
		a:        a,
		b:        b,
	}
}

func (b *Blackrock) getByte(rr uint64, n int, j int) uint64 {
	return ((rr >> (uint(n) * 8)) ^ uint64(b.seed) ^ uint64(j)) & 0xFF
}

func (b *Blackrock) read(j int, rr uint64) uint64 {
	rr ^= (uint64(b.seed) << uint(j)) ^ (uint64(b.seed) >> (64 - uint(j)))

	r0 := SBOX[b.getByte(rr, 0, j)]<<0 | SBOX[b.getByte(rr, 1, j)]<<8
	r1 := SBOX[b.getByte(rr, 2, j)]<<16 | SBOX[b.getByte(rr, 3, j)]<<24
	r2 := SBOX[b.getByte(rr, 4, j)]<<0 | SBOX[b.getByte(rr, 5, j)]<<8
	r3 := SBOX[b.getByte(rr, 6, j)]<<16 | SBOX[b.getByte(rr, 7, j)]<<24

	return r0 ^ r1 ^ r2<<23 ^ r3<<33
}

func (b *Blackrock) encrypt(m uint64) uint64 {
	rr, ll := divMod(m, b.a)
	for j := 1; j <= b.rounds; j++ {
		rr, ll = (ll+b.read(j, rr))%b.getMod(j), rr
	}
	if b.rounds&1 == 1 {
		return b.a*ll + rr
	}
	return b.a*rr + ll
}

func (b *Blackrock) decrypt(m uint64) uint64 {
	rr, ll := divMod(m, b.a)
	if b.rounds&1 == 1 {
		rr, ll = ll, rr
	}

	for j := b.rounds; j > 0; j-- {
		var tmp uint64
		if j&1 == 1 {
			tmp = b.read(j, ll)
			if tmp > rr {
				tmp -= rr
				tmp = b.a - (tmp % b.a)
				if tmp == b.a {
					tmp = 0
				}
			} else {
				tmp = rr - tmp
				tmp %= b.a
			}
		} else {
			tmp = b.read(j, ll)
			if tmp > rr {
				tmp -= rr
				tmp = b.b - (tmp % b.b)
				if tmp == b.b {
					tmp = 0
				}
			} else {
				tmp = rr - tmp
				tmp %= b.b
			}
		}
		rr, ll = ll, tmp
	}

	if b.rounds&1 == 1 {
		return b.a*rr + ll
	}
	return b.a*ll + rr
}

func (b *Blackrock) Shuffle(m uint64) uint64 {
	c := b.encrypt(m)
	for c >= b.rangeLen {
		c = b.encrypt(c)
	}
	return c
}

func (b *Blackrock) unshuffle(m uint64) uint64 {
	c := b.decrypt(m)
	for c >= b.rangeLen {
		c = b.decrypt(c)
	}
	return c
}

func (b *Blackrock) getMod(j int) uint64 {
	if j&1 == 1 {
		return b.a
	}
	return b.b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func divMod(n, d uint64) (uint64, uint64) {
	return n / d, n % d
}
