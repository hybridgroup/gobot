package joystick

var joyconPairConfig = joystickConfig{
	Name: "Nintendo Switch Joycon Controller Pair",
	GUID: "5555",
	Axis: []pair{
		{
			Name: "left_x",
			ID:   0,
		},
		{
			Name: "left_y",
			ID:   1,
		},
		{
			Name: "right_x",
			ID:   2,
		},
		{
			Name: "right_y",
			ID:   3,
		},
		{
			Name: "lt",
			ID:   4,
		},
		{
			Name: "rt",
			ID:   5,
		},
	},
	Buttons: []pair{
		{
			Name: "a",
			ID:   0,
		},
		{
			Name: "b",
			ID:   1,
		},
		{
			Name: "x",
			ID:   2,
		},
		{
			Name: "y",
			ID:   3,
		},
		{
			Name: "up",
			ID:   11,
		},
		{
			Name: "down",
			ID:   12,
		},
		{
			Name: "left",
			ID:   13,
		},
		{
			Name: "right",
			ID:   14,
		},
		{
			Name: "lb",
			ID:   9,
		},
		{
			Name: "rb",
			ID:   10,
		},
		{
			Name: "right_stick",
			ID:   8,
		},
		{
			Name: "left_stick",
			ID:   7,
		},
		{
			Name: "options",
			ID:   15,
		},
		{
			Name: "home",
			ID:   5,
		},
		{
			Name: "sr_left",
			ID:   17,
		},
		{
			Name: "sl_left",
			ID:   19,
		},
		{
			Name: "sr_right",
			ID:   16,
		},
		{
			Name: "sl_right",
			ID:   18,
		},
		{
			Name: "minus",
			ID:   4,
		},
		{
			Name: "plus",
			ID:   6,
		},
	},
}
