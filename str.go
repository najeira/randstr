package randstr

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxBits = 6
	letterIdxMask = 0x3F // 63 0b111111
	letterIdxMax  = 63 / letterIdxBits
)

var (
	mtSource mathrand.Source
)

func init() {
	mtSource = mathrand.NewSource(time.Now().UnixNano())
}

// String generates a random string using math/rand
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func String(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, mtSource.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = mtSource.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
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
		if idx < len(letterBytes) {
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
