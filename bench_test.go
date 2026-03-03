package errors

import (
	"fmt"
	"testing"
)

func BenchmarkWith(b *testing.B) {
	err := New("base error")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = With(err, KV("k1", "v1"), KV("k2", "v2"), KV("k3", "v3"), SeverityRuntime, Code("ERR_CODE"))
	}
}

func BenchmarkValue(b *testing.B) {
	err := New("base error")
	err = With(err, KV("k1", "v1"), KV("k2", "v2"), KV("k3", "v3"), SeverityRuntime, Code("ERR_CODE"))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Value(err, "k1")
	}
}

func BenchmarkError(b *testing.B) {
	err := New("base error")
	err = With(err, Op("my.Op"), KV("k1", "v1"), KV("k2", "v2"), SeverityRuntime, Code("ERR_CODE"))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = err.Error()
	}
}

func BenchmarkDefaultFormater(b *testing.B) {
	err := New("base error")
	err = With(err, Op("my.Op"), KV("k1", "v1"), KV("k2", "v2"), SeverityRuntime, Code("ERR_CODE"))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Format(err)
	}
}

func BenchmarkFullFormater(b *testing.B) {
	err := New("base error")
	err = With(err, Op("my.Op"), KV("k1", "v1"), KV("k2", "v2"), SeverityRuntime, Code("ERR_CODE"))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = FullFormater(err)
	}
}

func BenchmarkSprintf(b *testing.B) {
	err := New("base error")
	err = With(err, Op("my.Op"), KV("k1", "v1"), KV("k2", "v2"), SeverityRuntime, Code("ERR_CODE"))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", err)
	}
}

func BenchmarkSprintfPlus(b *testing.B) {
	err := New("base error")
	err = With(err, Op("my.Op"), KV("k1", "v1"), KV("k2", "v2"), SeverityRuntime, Code("ERR_CODE"))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%+v", err)
	}
}
