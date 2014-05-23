package leap

import (
	"encoding/json"
	"regexp"
)

type Gesture struct {
	Direction     []float64   `json:"direction"`
	Duration      int         `json:"duration"`
	Hands         []Hand      `json:"hands"`
	ID            int         `json:"id"`
	Pointables    []Pointable `json:"pointables"`
	Position      []float64   `json:"position"`
	Speed         float64     `json:"speed"`
	StartPosition []float64   `json:"StartPosition"`
	State         string      `json:"state"`
	Type          string      `json:"type"`
}

type Hand struct {
	Direction              []float64   `json:"direction"`
	ID                     int         `json:"id"`
	PalmNormal             []float64   `json:"palmNormal"`
	PalmPosition           []float64   `json:"PalmPosition"`
	PalmVelocity           []float64   `json:"PalmVelocity"`
	R                      [][]float64 `json:"r"`
	S                      float64     `json:"s"`
	SphereCenter           []float64   `json:"sphereCenter"`
	SphereRadius           float64     `json:"sphereRadius"`
	StabilizedPalmPosition []float64   `json:"stabilizedPalmPosition"`
	T                      []float64   `json:"t"`
	TimeVisible            float64     `json:"TimeVisible"`
}

type Pointable struct {
	Direction             []float64 `json:"direction"`
	HandID                int       `json:"handId"`
	ID                    int       `json:"id"`
	Length                float64   `json:"length"`
	StabilizedTipPosition []float64 `json:"stabilizedTipPosition"`
	TimeVisible           float64   `json:"timeVisible"`
	TipPosition           []float64 `json:"tipPosition"`
	TipVelocity           []float64 `json:"tipVelocity"`
	Tool                  bool      `json:"tool"`
	TouchDistance         float64   `json:"touchDistance"`
	TouchZone             string    `json:"touchZone"`
}

type InteractionBox struct {
	Center []int     `json:"center"`
	Size   []float64 `json:"size"`
}

type Frame struct {
	CurrentFrameRate float64        `json:"currentFrameRate"`
	Gestures         []Gesture      `json:"gestures"`
	Hands            []Hand         `json:"hands"`
	ID               int            `json:"id"`
	InteractionBox   InteractionBox `json:"interactionBox"`
	Pointables       []Pointable    `json:"pointables"`
	R                [][]float64    `json:"r"`
	S                float64        `json:"s"`
	T                []float64      `json:"t"`
	Timestamp        int            `json:"timestamp"`
}

func (h *Hand) X() float64 {
	return h.PalmPosition[0]
}
func (h *Hand) Y() float64 {
	return h.PalmPosition[1]
}
func (h *Hand) Z() float64 {
	return h.PalmPosition[2]
}

func (l *LeapMotionDriver) ParseFrame(data []byte) Frame {
	var frame Frame
	json.Unmarshal(data, &frame)
	return frame
}

func (l *LeapMotionDriver) isAFrame(data []byte) bool {
	match, _ := regexp.Match("currentFrameRate", data)
	return match
}
