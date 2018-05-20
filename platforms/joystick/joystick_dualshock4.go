package joystick

var dualshock4Config = joystickConfig{
	Name: "Dualshock4 Controller",
	GUID: "2222",
	Axis: []pair{
		pair{Name: "left_x", ID: 0},
		pair{Name: "left_y", ID: 1},
		pair{Name: "right_x", ID: 3},
		pair{Name: "right_y", ID: 4},
		pair{Name: "l2", ID: 2},
		pair{Name: "r2", ID: 5},
		pair{Name: "up", ID: 7},
		pair{Name: "down", ID: 7},
		pair{Name: "left", ID: 6},
		pair{Name: "right", ID: 6},
	},
	Buttons: []pair{
		pair{Name: "x", ID: 0},
		pair{Name: "circle", ID: 1},
		pair{Name: "triangle", ID: 2},
		pair{Name: "square", ID: 3},
		pair{Name: "l1", ID: 4},
		pair{Name: "r1", ID: 5},
		pair{Name: "l2", ID: 6},
		pair{Name: "r2", ID: 7},
		pair{Name: "share", ID: 8},
		pair{Name: "options", ID: 9},
		pair{Name: "ps", ID: 10},    // the 'PlayStation' button
		pair{Name: "left", ID: 11},  // push on left stick
		pair{Name: "right", ID: 12}, // push on right stick
	},
	Hats: []hat{},
}
