package leap

import (
	"encoding/json"
)

// Gesture is a Leap Motion gesture tht has been detected
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

// Hand is a Leap Motion hand tht has been detected
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

// Pointable is a Leap Motion pointing motion tht has been detected
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

// InteractionBox is the area within which the gestural interaction has been detected
type InteractionBox struct {
	Center []int     `json:"center"`
	Size   []float64 `json:"size"`
}

// Frame is the base representation returned that holds every other objects
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

// X returns hand x value
func (h *Hand) X() float64 {
	return h.PalmPosition[0]
}

// Y returns hand y value
func (h *Hand) Y() float64 {
	return h.PalmPosition[1]
}

// Z returns hand z value
func (h *Hand) Z() float64 {
	return h.PalmPosition[2]
}

// ParseFrame converts json data to a Frame
func (l *Driver) ParseFrame(data []byte) Frame {
	var frame Frame
	json.Unmarshal(data, &frame)
	return frame
}
