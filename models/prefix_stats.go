package models

import "time"

type PrefixStats struct {
	Bytes   uint64
	Packets uint64
	Ts      time.Time
}
