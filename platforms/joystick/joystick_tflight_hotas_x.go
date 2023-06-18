package joystick

var tflightHotasXConfig = joystickConfig{
	Name: "Thrustmaster T-Flight Hotas X Joystick",
	GUID: "4444",
	Axis: []pair{
		{Name: "right_x", ID: 0},
		{Name: "right_y", ID: 1},
		{Name: "left_y", ID: 2},
		{Name: "r1", ID: 3}, // RH Twist
		{Name: "left_x", ID: 4},
	},
	Buttons: []pair{
		{Name: "r1", ID: 0},
		{Name: "l1", ID: 1},
		{Name: "r3", ID: 2},
		{Name: "l3", ID: 3},
		{Name: "square", ID: 4},
		{Name: "x", ID: 5},
		{Name: "circle", ID: 6},
		{Name: "triangle", ID: 7},
		{Name: "r2", ID: 8},
		{Name: "l2", ID: 9},
		{Name: "select", ID: 10},
		{Name: "start", ID: 11},
	},
	Hats: []hat{
		{Hat: 0, Name: "down", ID: 4},
		{Hat: 0, Name: "up", ID: 1},
		{Hat: 0, Name: "left", ID: 8},
		{Hat: 0, Name: "right", ID: 2},
		{Hat: 0, Name: "released", ID: 0},
	},
}
