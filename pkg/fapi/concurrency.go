package fapi

import "sync"

type positionInfo struct {
	m map[string]float64
	sync.Mutex
}
