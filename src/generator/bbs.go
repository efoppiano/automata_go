package generator

import (
	"math"

	"lukechampine.com/uint128"
)

type BlumBlumShub struct {
	curr uint64
	m    uint64
}

func NewBlumBlumShub(seed uint64) *BlumBlumShub {
	p := nextUsablePrime(800000000)
	q := nextUsablePrime(400000000)
	m := p * q
	return &BlumBlumShub{seed, m}
}

func (bbs *BlumBlumShub) next() float64 {
	curr_128 := uint128.From64(bbs.curr)

	bbs.curr = curr_128.Mul(curr_128).Mod64(bbs.m)
	ret := float64(bbs.curr) / float64(bbs.m)
	return ret
}

func nextUsablePrime(x uint64) uint64 {
	p := nextPrime(x)
	for p%4 != 3 {
		p = nextPrime(p)
	}
	return p
}

func nextPrime(x uint64) uint64 {
	for i := x + 1; ; i++ {
		if isPrime(i) {
			return i
		}
	}
}

func isPrime(x uint64) bool {
	if x < 2 {
		return false
	}
	var i uint64
	for i = 2; i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func (bbs *BlumBlumShub) Random() float64 {
	return bbs.next()
}

// / Returns a random integer in the range [a, b)
func (bbs *BlumBlumShub) RandInt(a int, b int) int {
	p := (float64(b-a))*bbs.Random() + float64(a)
	return int(p)
}

func (bbs *BlumBlumShub) Poi(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0
	for p > L {
		k++
		v := bbs.Random()
		p *= v
	}
	return k - 1
}

func (bbs *BlumBlumShub) Choice(a []any) int {
	i := bbs.RandInt(0, len(a))
	return i
}
