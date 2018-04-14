package joystick

var xbox360Config = joystickConfig{
	Name: "XBox 360 Controller",
	GUID: "3333",
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
			ID:   3,
		},
		pair{
			Name: "right_y",
			ID:   4,
		},
		pair{
			Name: "rt",
			ID:   5,
		},
		pair{
			Name: "lt",
			ID:   2,
		},
	},
	Buttons: []pair{
		pair{
			Name: "x",
			ID:   2,
		},
		pair{
			Name: "a",
			ID:   0,
		},
		pair{
			Name: "b",
			ID:   1,
		},
		pair{
			Name: "y",
			ID:   3,
		},
		pair{
			Name: "lb",
			ID:   4,
		},
		pair{
			Name: "rb",
			ID:   5,
		},
		pair{
			Name: "back",
			ID:   6,
		},
		pair{
			Name: "start",
			ID:   7,
		},
		pair{
			Name: "home",
			ID:   8,
		},
		pair{
			Name: "right_stick",
			ID:   10,
		},
		pair{
			Name: "left_stick",
			ID:   9,
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
