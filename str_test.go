package randstr

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"math/big"
	"testing"
)

const (
	genLen    = 100
	loopCount = 1000000
)

// tests

func TestString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s := String(genLen)
		check(t, store, nil, s)
	}
}

func TestCryptoString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s := CryptoString(genLen)
		check(t, store, nil, s)
	}
}

func TestPrivateCryptoString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s, err := cryptoString(genLen)
		check(t, store, err, s)
	}
}

func TestNumericString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s := NumericString(genLen)
		check(t, store, nil, s)
		for _, c := range s {
			if c < '0' || '9' < c {
				t.Error(s)
				break
			}
		}
	}
}

func TestCryptoNumericString(t *testing.T) {
	store := make(map[string]bool)
	for i := 1; i < loopCount; i++ {
		s, err := cryptoNumericString(genLen)
		check(t, store, err, s)
		for _, c := range s {
			if c < '0' || '9' < c {
				t.Error(s)
				break
			}
		}
	}
}

func check(t *testing.T, store map[string]bool, err error, s string) {
	if err != nil {
		t.Error(err)
	} else if _, exists := store[s]; exists {
		t.Errorf("already generated %s", s)
	} else if len(s) != genLen {
		t.Errorf("invalid length %s", s)
	} else {
		for i := range s {
			if r := s[i]; '0' <= r && r <= '9' {
				// ok
			} else if 'A' <= r && r <= 'z' {
				// ok
			} else {
				t.Errorf("invalid character %s", s)
			}
		}
	}
	store[s] = true
}

// other implementasions

func cryptoPhase1(n int) string {
	buf := make([]byte, n)
	max := new(big.Int)
	max.SetInt64(int64(len(letterBytes)))
	for i := range buf {
		r, err := cryptorand.Int(cryptorand.Reader, max)
		if err != nil {
			panic(err)
		}
		buf[i] = letterBytes[r.Int64()]
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
		idx := int(src[0] & letterIdxMask)
		if idx < len(letterBytes) {
			buf[i] = letterBytes[idx]
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
		idx := int(src[pos] & letterIdxMask)
		if idx < len(letterBytes) {
			buf[i] = letterBytes[idx]
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
			remain = letterIdxTimes
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
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
		String(genLen)
	}
}

func BenchmarkCryptoString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = cryptoString(genLen)
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
		NumericString(genLen)
	}
}

func BenchmarkCryptoNumericString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CryptoNumericString(genLen)
	}
}
