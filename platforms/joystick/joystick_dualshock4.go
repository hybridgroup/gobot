package joystick

var dualshock4Config = joystickConfig{
	Name: "Dualshock4 Controller",
	GUID: "2222",
	Axis: []pair{
		pair{
			Name: "left_x",
			ID:   0,
		},
		pair{
			Name: "left_y",
			ID:   1,
		},
		pair{
			Name: "right_x",
			ID:   2,
		},
		pair{
			Name: "right_y",
			ID:   5,
		},
		pair{
			Name: "l2",
			ID:   3,
		},
		pair{
			Name: "r2",
			ID:   4,
		},
		pair{
			Name: "up",
			ID:   7,
		},
		pair{
			Name: "down",
			ID:   6,
		},
		pair{
			Name: "left",
			ID:   7,
		},
		pair{
			Name: "right",
			ID:   8,
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
			Name: "r2",
			ID:   7,
		},
		pair{
			Name: "r1",
			ID:   5,
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
	},
	Hats: []hat{},
}
