package models

import (
	"net"
	"sync"
	"time"

	"github.com/yl2chen/cidranger"
)

// Custom RangerEntry with a RingBuffer
type CustomRangerEntry struct {
	ipNet      net.IPNet
	ringBuffer *RingBuffer
	mu         sync.Mutex // To ensure thread-safe access
}

// Implement the RangerEntry interface
func (e *CustomRangerEntry) Network() net.IPNet {
	return e.ipNet
}

// AddStats adds stats to the RingBuffer (thread-safe)
func (e *CustomRangerEntry) AddStats(bytes, packets uint64, ts time.Time) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ringBuffer.Add(bytes, packets, ts)
}

// GetStats aggregates stats from the RingBuffer (thread-safe)
func (e *CustomRangerEntry) GetStats(window time.Duration, now time.Time) (uint64, uint64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.ringBuffer.GetAggregate(window, now)
}

// Constructor for CustomRangerEntry
func NewCustomRangerEntry(ipNet net.IPNet, bufferSize int) cidranger.RangerEntry {
	return &CustomRangerEntry{
		ipNet:      ipNet,
		ringBuffer: NewRingBuffer(bufferSize),
	}
}
