package models

import (
	"fmt"
	"net"
	"time"

	"github.com/yl2chen/cidranger"
)

type PrefixStore struct {
	ranger cidranger.Ranger
}

// Initialize a PrefixStore
func NewPrefixStore() *PrefixStore {
	return &PrefixStore{ranger: cidranger.NewPCTrieRanger()}
}

// Add a prefix with a ring buffer of specified size
func (ps *PrefixStore) AddPrefix(prefix string, bufferSize int) error {
	_, network, err := net.ParseCIDR(prefix)
	if err != nil {
		return err
	}
	entry := NewCustomRangerEntry(*network, bufferSize)
	ps.ranger.Insert(entry)
	return nil
}

// Update stats for an IP address
func (ps *PrefixStore) UpdateStats(ip string, bytes, packets uint64, ts time.Time) error {
	parsedIP := net.ParseIP(ip)
	entries, err := ps.ranger.ContainingNetworks(parsedIP)
	if err != nil || len(entries) == 0 {
		return fmt.Errorf("no matching prefix for IP: %s", ip)
	}

	for _, e := range entries {
		entry, ok := e.(*CustomRangerEntry)
		if ok {
			entry.AddStats(bytes, packets, ts)
		}
	}
	return nil
}

// Query aggregated stats for a prefix
func (ps *PrefixStore) QueryStats(prefix string, window time.Duration, now time.Time) (uint64, uint64, error) {
	_, network, err := net.ParseCIDR(prefix)
	if err != nil {
		return 0, 0, err
	}

	// Get the first IP address in the CIDR range
	firstIP := network.IP

	// Use ContainingNetworks to find entries for the first IP
	entries, err := ps.ranger.ContainingNetworks(firstIP)
	if err != nil || len(entries) == 0 {
		return 0, 0, fmt.Errorf("no entries found for CIDR: %s", prefix)
	}

	// Iterate over the entries and match the exact CIDR
	for _, e := range entries {
		customEntry, ok := e.(*CustomRangerEntry)
		if !ok {
			continue
		}
		n := customEntry.Network()
		// Check if the entry's network matches the requested CIDR
		if n.String() == network.String() {
			bytes, pkts := customEntry.GetStats(window, now)
			return bytes, pkts, nil
		}
	}

	return 0, 0, fmt.Errorf("no exact match found for CIDR: %s", prefix)
}
