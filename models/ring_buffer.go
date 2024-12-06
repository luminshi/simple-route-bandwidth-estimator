package models

import "time"

type RingBuffer struct {
	Stats []PrefixStats
	Index int
	Size  int
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		Stats: make([]PrefixStats, size),
		Index: 0,
		Size:  size,
	}
}

func (rb *RingBuffer) Add(bytes, packets uint64, ts time.Time) {
	rb.Stats[rb.Index] = PrefixStats{Bytes: bytes, Packets: packets, Ts: ts}
	rb.Index = (rb.Index + 1) % rb.Size
}

func (rb *RingBuffer) GetAggregate(window time.Duration, now time.Time) (uint64, uint64) {
	var totalBytes, totalPackets uint64
	for _, stat := range rb.Stats {
		if now.Sub(stat.Ts) <= window {
			totalBytes += stat.Bytes
			totalPackets += stat.Packets
		}
	}
	return totalBytes, totalPackets
}
