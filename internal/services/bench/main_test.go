package main

import (
	"encoding/json"
	"strings"
	"testing"
)

var numbers = []string{"hello", "sello", "mello", "allo", "ooppo"}

func BenchmarkStringsJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for index := range numbers {
			numbers[index] = numbers[index] + " 1"
		}
		strings.Join(numbers, ",")
	}
}

func BenchmarkJsonMarshal(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		_, err = json.Marshal(numbers)
		if err != nil {
			b.Fatal(err)
		}
	}
}
