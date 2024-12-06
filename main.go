package main

import (
	models "bwgetter/models"
	"fmt"
	"time"
)

func main() {
	store := models.NewPrefixStore()

	// Add prefixes with a ring buffer size of 60
	store.AddPrefix("192.168.0.0/23", 60)
	store.AddPrefix("192.168.1.0/24", 60)
	store.AddPrefix("10.0.0.0/8", 60)

	// Update stats for IPs
	now := time.Now()
	store.UpdateStats("192.168.1.1", 1024, 10, now)
	store.UpdateStats("192.168.0.1", 1024, 10, now)
	store.UpdateStats("10.0.0.5", 2048, 20, now)

	// Query stats for a prefix
	bytes, packets, err := store.QueryStats("192.168.0.0/23", 1*time.Minute, time.Now())
	if err != nil {
		fmt.Println("Error querying stats:", err)
	} else {
		fmt.Printf("Prefix Stats - Bytes: %d, Packets: %d\n", bytes, packets)
	}
}
