package fapi

import "sync"

type positionInfo struct {
	m        map[string]float64
	position float64
	sync.RWMutex
}
