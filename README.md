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
BenchmarkString-8                	 4250818	       279.7 ns/op	     112 B/op	       1 allocs/op
BenchmarkCryptoString-8          	 3017208	       383.5 ns/op	     224 B/op	       2 allocs/op
BenchmarkNumericString-8         	 3273469	       367.1 ns/op	     112 B/op	       1 allocs/op
BenchmarkCryptoNumericString-8   	 2707904	       445.6 ns/op	     224 B/op	       2 allocs/op
```

on Mac mini (M1, 2020) Apple M1 16 GB

## License

MIT
