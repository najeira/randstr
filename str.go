package randstr

import (
	cryptorand "crypto/rand"
	"math"
	"math/big"
	mathrand "math/rand"
	"sync"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterCount   = len(letterBytes)
	letterIdxBits = 6
	letterIdxMask = 0x3F // 63 0b111111
	letterIdxMax  = 63 / letterIdxBits
)

var (
	mutex sync.Mutex
	mtSource mathrand.Source
)

func init() {
	var seed int64
	bint := big.NewInt(math.MaxInt64)
	bseed, err := cryptorand.Int(cryptorand.Reader, bint)
	if bseed != nil && err == nil {
		seed = bseed.Int64()
	} else {
		seed = time.Now().UnixNano()
	}
	mtSource = mathrand.NewSource(seed)
}

// String generates a random string using math/rand
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func String(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, mathRandInt63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = mathRandInt63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < letterCount {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func cryptoString(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := cryptorand.Read(buf); err != nil {
		return "", err
	}
	for i := 0; i < n; {
		idx := int(buf[i] & letterIdxMask)
		if idx < letterCount {
			buf[i] = letterBytes[idx]
			i++
		} else {
			if _, err := cryptorand.Read(buf[i : i+1]); err != nil {
				return "", err
			}
		}
	}
	return string(buf), nil
}

// CryptoString generates a random string using crypto/rand
func CryptoString(n int) string {
	s, err := cryptoString(n)
	if err != nil {
		return String(n)
	}
	return s
}

func Seed(seed int64) {
	mutex.Lock()
	mtSource.Seed(seed)
	mutex.Unlock()
}

func mathRandInt63() int64 {
	mutex.Lock()
	n := mtSource.Int63()
	mutex.Unlock()
	return n
}
