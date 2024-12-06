package models

import (
	"fmt"
	"testing"
)

func BenchmarkAddPrefixDifferentSizes(b *testing.B) {
	store := NewPrefixStore()
	for i := 0; i < b.N; i++ {
		prefix := fmt.Sprintf("%d.%d.%d.0/%d", i%256, (i/256)%256, (i/65536)%256, 16+(i%8))
		err := store.AddPrefix(prefix, 60)
		if err != nil {
			b.Fatalf("Error adding prefix: %v", err)
		}
	}
}

func BenchmarkAddPrefixMemory(b *testing.B) {
	store := NewPrefixStore()
	prefix := "192.168.1.0/24"
	b.ResetTimer()

	allocs := testing.AllocsPerRun(1000, func() {
		err := store.AddPrefix(prefix, 60)
		if err != nil {
			b.Fatalf("Error adding prefix: %v", err)
		}
	})
	b.Logf("Memory allocations per call: %.2f", allocs)
}
