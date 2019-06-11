package main

import (
	"github.com/google/uuid"
)

type WsEvent struct {
	Event string
	ID    uuid.UUID
	Data  interface{}
}

type DLocation struct {
	Longitude float64
	Latitude  float64
	Accuracy  float64
	Speed     float64
}

type DVisibleArea struct {
	Radius float64
}
