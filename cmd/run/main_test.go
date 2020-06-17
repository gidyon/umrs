package main

import (
	"github.com/google/uuid"
	"testing"
)

func Benchmark_genID(b *testing.B) {
	var id string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id = genID()
		_ = id
	}

	b.StopTimer()
}

func Benchmark_uuid(b *testing.B) {
	var id string

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		id = uuid.New().String()
		_ = id
	}

	b.StopTimer()
}
