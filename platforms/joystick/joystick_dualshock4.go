package joystick

var dualshock4Config = joystickConfig{
	Name: "Dualshock4 Controller",
	GUID: "2222",
	Axis: []pair{
		pair{
			Name: "left_x",
			ID:   1,
		},
		pair{
			Name: "left_y",
			ID:   0,
		},
		pair{
			Name: "right_x",
			ID:   5,
		},
		pair{
			Name: "right_y",
			ID:   2,
		},
	},
	Buttons: []pair{
		pair{
			Name: "square",
			ID:   0,
		},
		pair{
			Name: "triangle",
			ID:   3,
		},
		pair{
			Name: "circle",
			ID:   2,
		},
		pair{
			Name: "x",
			ID:   1,
		},
		pair{
			Name: "l1",
			ID:   4,
		},
		pair{
			Name: "l2",
			ID:   6,
		},
		pair{
			Name: "l3",
			ID:   10,
		},
		pair{
			Name: "r1",
			ID:   5,
		},
		pair{
			Name: "r2",
			ID:   7,
		},
		pair{
			Name: "r3",
			ID:   11,
		},
		pair{
			Name: "share",
			ID:   8,
		},
		pair{
			Name: "options",
			ID:   9,
		},
		pair{
			Name: "home",
			ID:   12,
		},
		pair{
			Name: "touchpad",
			ID:   13,
		},
	},
	Hats: []hat{
		hat{
			Hat:  0,
			Name: "down",
			ID:   4,
		},
		hat{
			Hat:  0,
			Name: "up",
			ID:   1,
		},
		hat{
			Hat:  0,
			Name: "left",
			ID:   8,
		},
		hat{
			Hat:  0,
			Name: "right",
			ID:   2,
		},
		hat{
			Hat:  0,
			Name: "released",
			ID:   0,
		},
	},
}
