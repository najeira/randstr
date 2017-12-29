# randstr

Generate random string using crypto/rand and math/rand for Go.

## Usage

String generates a random string using math/rand:

```go
s := String(10)
```

CryptoString generates a random string using crypto/rand:

```go
s := CryptoString(10)
```

## Benchmark

```
BenchmarkString-4                5000000               306 ns/op             224 B/op          2 allocs/op
BenchmarkCryptoString-4           200000              8272 ns/op             224 B/op          2 allocs/op
```

on Macbook Pro 3.3 GHz Intel Core i7

## License

MIT
