package leap

import (
	"encoding/json"
)

// Gesture is a Leap Motion gesture that has been detected
type Gesture struct {
	Center        []float64   `json:"center"`
	Direction     []float64   `json:"direction"`
	Duration      int         `json:"duration"`
	HandIDs       []int       `json:"handIds"`
	ID            int         `json:"id"`
	Normal        []float64   `json:"normal"`
	PointableIDs  []int       `json:"pointableIds"`
	Position      []float64   `json:"position"`
	Progress      float64     `json:"progress"`
	Radius        float64     `json:"radius"`
	Speed         float64     `json:"speed"`
	StartPosition []float64   `json:"StartPosition"`
	State         string      `json:"state"`
	Type          string      `json:"type"`
}

// Hand is a Leap Motion hand that has been detected
type Hand struct {
	ArmBasis               [][]float64 `json:"armBasis"`
	ArmWidth               float64     `json:"armWidth"`
	Confidence             float64     `json:"confidence"`
	Direction              []float64   `json:"direction"`
	Elbow                  []float64   `json:"elbow"`
	GrabStrength           float64     `json:"grabStrength"`
	ID                     int         `json:"id"`
	PalmNormal             []float64   `json:"palmNormal"`
	PalmPosition           []float64   `json:"PalmPosition"`
	PalmVelocity           []float64   `json:"PalmVelocity"`
	PinchStrength          float64     `json:"PinchStrength"`
	R                      [][]float64 `json:"r"`
	S                      float64     `json:"s"`
	SphereCenter           []float64   `json:"sphereCenter"`
	SphereRadius           float64     `json:"sphereRadius"`
	StabilizedPalmPosition []float64   `json:"stabilizedPalmPosition"`
	T                      []float64   `json:"t"`
	TimeVisible            float64     `json:"TimeVisible"`
	Type                   string      `json:"type"`
	Wrist                  []float64   `json:"wrist"`
}

// Pointable is a Leap Motion pointing motion that has been detected
type Pointable struct {
	Bases                 [][][]float64  `json:"bases"`
	BTipPosition          []float64    `json:"btipPosition"`
	CarpPosition          []float64    `json:"carpPosition"`
	DipPosition           []float64    `json:"dipPosition"`
	Direction             []float64    `json:"direction"`
	Extended              bool         `json:"extended"`
	HandID                int          `json:"handId"`
	ID                    int          `json:"id"`
	Length                float64      `json:"length"`
	MCPPosition           []float64    `json:"mcpPosition"`
	PIPPosition           []float64    `json:"pipPosition"`
	StabilizedTipPosition []float64    `json:"stabilizedTipPosition"`
	TimeVisible           float64      `json:"timeVisible"`
	TipPosition           []float64    `json:"tipPosition"`
	TipVelocity           []float64    `json:"tipVelocity"`
	Tool                  bool         `json:"tool"`
	TouchDistance         float64      `json:"touchDistance"`
	TouchZone             string       `json:"touchZone"`
	Type                  int          `json:"type"`
	Width                 float64      `json:"width"`
}

// InteractionBox is the area within which the gestural interaction has been detected
type InteractionBox struct {
	Center []float64 `json:"center"`
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
	Timestamp        uint64         `json:"timestamp"`
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
