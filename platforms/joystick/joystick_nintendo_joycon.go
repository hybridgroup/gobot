package joystick

var joyconPairConfig = joystickConfig{
	Name: "Nintendo Switch Joycon Controller Pair",
	GUID: "5555",
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
			ID:   3,
		},
		pair{
			Name: "lt",
			ID:   4,
		},
		pair{
			Name: "rt",
			ID:   5,
		},
	},
	Buttons: []pair{
		pair{
			Name: "a",
			ID:   0,
		},
		pair{
			Name: "b",
			ID:   1,
		},
		pair{
			Name: "x",
			ID:   2,
		},
		pair{
			Name: "y",
			ID:   3,
		},
		pair{
			Name: "up",
			ID:   11,
		},
		pair{
			Name: "down",
			ID:   12,
		},
		pair{
			Name: "left",
			ID:   13,
		},
		pair{
			Name: "right",
			ID:   14,
		},
		pair{
			Name: "lb",
			ID:   9,
		},
		pair{
			Name: "rb",
			ID:   10,
		},
		pair{
			Name: "right_stick",
			ID:   8,
		},
		pair{
			Name: "left_stick",
			ID:   7,
		},
		pair{
			Name: "options",
			ID:   15,
		},
		pair{
			Name: "home",
			ID:   5,
		},
		pair{
			Name: "sr_left",
			ID:   17,
		},
		pair{
			Name: "sl_left",
			ID:   19,
		},
		pair{
			Name: "sr_right",
			ID:   16,
		},
		pair{
			Name: "sl_right",
			ID:   18,
		},
		pair{
			Name: "minus",
			ID:   4,
		},
		pair{
			Name: "plus",
			ID:   6,
		},
	},
	Hats: []hat{},
}
