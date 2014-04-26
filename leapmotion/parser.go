package gobotLeap

import (
	"encoding/json"
	"regexp"
)

type LeapGesture struct {
	Direction     []float64       `json:"direction"`
	Duration      int             `json:"duration"`
	Hands         []LeapHand      `json:"hands"`
	ID            int             `json:"id"`
	Pointables    []LeapPointable `json:"pointables"`
	Position      []float64       `json:"position"`
	Speed         float64         `json:"speed"`
	StartPosition []float64       `json:"StartPosition"`
	State         string          `json:"state"`
	Type          string          `json:"type"`
}

type LeapHand struct {
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

type LeapPointable struct {
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

type LeapInteractionBox struct {
	Center []int     `json:"center"`
	Size   []float64 `json:"size"`
}

type LeapFrame struct {
	CurrentFrameRate float64            `json:"currentFrameRate"`
	Gestures         []LeapGesture      `json:"gestures"`
	Hands            []LeapHand         `json:"hands"`
	ID               int                `json:"id"`
	InteractionBox   LeapInteractionBox `json:"interactionBox"`
	Pointables       []LeapPointable    `json:"pointables"`
	R                [][]float64        `json:"r"`
	S                float64            `json:"s"`
	T                []float64          `json:"t"`
	Timestamp        int                `json:"timestamp"`
}

func (this *LeapHand) X() float64 {
	return this.PalmPosition[0]
}
func (this *LeapHand) Y() float64 {
	return this.PalmPosition[1]
}
func (this *LeapHand) Z() float64 {
	return this.PalmPosition[2]
}

func (l *LeapDriver) ParseLeapFrame(data []byte) LeapFrame {
	var frame LeapFrame
	json.Unmarshal(data, &frame)
	return frame
}

func (l *LeapDriver) isAFrame(data []byte) bool {
	match, _ := regexp.Match("currentFrameRate", data)
	return match
}
