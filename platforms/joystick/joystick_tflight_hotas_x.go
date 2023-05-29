package joystick

var tflightHotasXConfig = joystickConfig{
	Name: "Thrustmaster T-Flight Hotas X Joystick",
	GUID: "4444",
	Axis: []pair{
		pair{Name: "right_x", ID: 0},
		pair{Name: "right_y", ID: 1},
		pair{Name: "left_y", ID: 2},
		pair{Name: "r1", ID: 3}, // RH Twist
		pair{Name: "left_x", ID: 4},
	},
	Buttons: []pair{
		pair{Name: "r1", ID: 0},
		pair{Name: "l1", ID: 1},
		pair{Name: "r3", ID: 2},
		pair{Name: "l3", ID: 3},
		pair{Name: "square", ID: 4},
		pair{Name: "x", ID: 5},
		pair{Name: "circle", ID: 6},
		pair{Name: "triangle", ID: 7},
		pair{Name: "r2", ID: 8},
		pair{Name: "l2", ID: 9},
		pair{Name: "select", ID: 10},
		pair{Name: "start", ID: 11},
	},
	Hats: []hat{
		hat{Hat: 0, Name: "down", ID: 4},
		hat{Hat: 0, Name: "up", ID: 1},
		hat{Hat: 0, Name: "left", ID: 8},
		hat{Hat: 0, Name: "right", ID: 2},
		hat{Hat: 0, Name: "released", ID: 0},
	},
}
