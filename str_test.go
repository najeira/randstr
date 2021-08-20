package randstr_test

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/najeira/randstr"
)

const (
	genLen    = 100
	loopCount = 1000000
)

// tests

func TestString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s := randstr.String(genLen)
		check(t, store, nil, s, false)
	}
}

func TestStringVariation(t *testing.T) {
	store := make(map[rune]int)
	for i := 1; i < loopCount; i++ {
		s := randstr.String(genLen)
		for _, r := range s {
			store[r]++
		}
	}
	checkVariation(t, store)
}

func TestCryptoString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s := randstr.CryptoString(genLen)
		check(t, store, nil, s, false)
	}
}

func TestPrivateCryptoString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s, err := randstr.ExportCryptoString(genLen)
		check(t, store, err, s, false)
	}
}

func TestCryptoStringVariation(t *testing.T) {
	store := make(map[rune]int)
	for i := 1; i < loopCount; i++ {
		s, err := randstr.ExportCryptoString(genLen)
		if err != nil {
			t.Error(err)
		}
		for _, r := range s {
			store[r]++
		}
	}
	checkVariation(t, store)
}

func TestNumericString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s := randstr.NumericString(genLen)
		check(t, store, nil, s, true)
	}
}

func TestNumericStringVariation(t *testing.T) {
	store := make(map[rune]int)
	for i := 1; i < loopCount; i++ {
		s := randstr.NumericString(genLen)
		for _, r := range s {
			store[r]++
		}
	}
	checkVariation(t, store)
}

func TestCryptoNumericString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s, err := randstr.ExportCryptoNumericString(genLen)
		check(t, store, err, s, true)
	}
}

func TestCryptoNumericStringVariation(t *testing.T) {
	store := make(map[rune]int)
	for i := 1; i < loopCount; i++ {
		s, err := randstr.ExportCryptoNumericString(genLen)
		if err != nil {
			t.Error(err)
		}
		for _, r := range s {
			store[r]++
		}
	}
	checkVariation(t, store)
}

func check(t *testing.T, store map[string]bool, err error, s string, numeric bool) {
	if err != nil {
		t.Error(err)
	} else if _, exists := store[s]; exists {
		t.Errorf("already generated %s", s)
	} else if len(s) != genLen {
		t.Errorf("invalid length %s", s)
	} else {
		var invalid bool
		for i := range s {
			if r := s[i]; '0' <= r && r <= '9' {
				// ok
			} else if 'A' <= r && r <= 'z' && !numeric {
				// ok
			} else {
				invalid = true
			}
		}
		if invalid {
				t.Errorf("invalid character %s", s)
		}
	}
	store[s] = true
}

func checkVariation(t *testing.T, store map[rune]int) {
	var max, min int
	for _, n := range store {
		if n > max {
			max = n
		}
		if n < min || min == 0 {
			min = n
		}
	}
	diff := max - min

	total := loopCount * genLen
	avg := total / len(store)
	threshold := avg * 3 / 100
	if diff > threshold {
		t.Errorf("%d %d", min, max)
	}
	//println(max, "-", min, "=", diff, "<", threshold)
}

// other implementasions

func cryptoPhase1(n int) string {
	buf := make([]byte, n)
	max := new(big.Int)
	max.SetInt64(int64(len(randstr.LetterBytes)))
	for i := range buf {
		r, err := cryptorand.Int(cryptorand.Reader, max)
		if err != nil {
			panic(err)
		}
		buf[i] = randstr.LetterBytes[r.Int64()]
	}
	return string(buf)
}

func cryptoPhase2(n int) string {
	src := make([]byte, 1)
	buf := make([]byte, n)
	for i := 0; i < n; {
		if _, err := cryptorand.Read(src); err != nil {
			panic(err)
		}
		idx := int(src[0] & randstr.LetterIdxMask)
		if idx < len(randstr.LetterBytes) {
			buf[i] = randstr.LetterBytes[idx]
			i++
		}
	}
	return string(buf)
}

func cryptoPhase3(n int) string {
	src := make([]byte, n)
	buf := make([]byte, n)
	for i, j := 0, 0; i < n; j++ {
		pos := j % n
		if pos == 0 {
			if _, err := cryptorand.Read(src); err != nil {
				panic(err)
			}
		}
		idx := int(src[pos] & randstr.LetterIdxMask)
		if idx < len(randstr.LetterBytes) {
			buf[i] = randstr.LetterBytes[idx]
			i++
		}
	}
	return string(buf)
}

func cryptoPhase4(n int) string {
	var cache uint64
	b := make([]byte, n)
	for i, remain := n-1, 0; i >= 0; {
		if remain == 0 {
			err := binary.Read(cryptorand.Reader, binary.LittleEndian, &cache)
			if err != nil {
				panic(err)
			}
			remain = randstr.LetterIdxTimes
		}
		if idx := int(cache & randstr.LetterIdxMask); idx < len(randstr.LetterBytes) {
			b[i] = randstr.LetterBytes[idx]
			i--
		}
		cache >>= 6
		remain--
	}
	return string(b)
}

// benchmarks

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randstr.String(genLen)
	}
}

func BenchmarkCryptoString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = randstr.ExportCryptoString(genLen)
	}
}

func BenchmarkCryptoPhase1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cryptoPhase1(genLen)
	}
}

func BenchmarkCryptoPhase2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cryptoPhase2(genLen)
	}
}

func BenchmarkCryptoPhase3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cryptoPhase3(genLen)
	}
}

func BenchmarkCryptoPhase4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cryptoPhase4(genLen)
	}
}

func BenchmarkNumericString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randstr.NumericString(genLen)
	}
}

func BenchmarkCryptoNumericString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = randstr.ExportCryptoNumericString(genLen)
	}
}
