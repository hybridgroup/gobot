package joystick

var xboxOneConfig = joystickConfig{
	Name: "Xbox One Controller",
	GUID: "30001983600083000100",
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
			ID:   3,
		},
		{
			Name: "right_y",
			ID:   4,
		},
		{
			Name: "rt",
			ID:   5,
		},
		{
			Name: "lt",
			ID:   2,
		},
	},
	Buttons: []pair{
		{
			Name: "x",
			ID:   2,
		},
		{
			Name: "a",
			ID:   0,
		},
		{
			Name: "b",
			ID:   1,
		},
		{
			Name: "y",
			ID:   3,
		},
		{
			Name: "lb",
			ID:   4,
		},
		{
			Name: "rb",
			ID:   5,
		},
		{
			Name: "back",
			ID:   6,
		},
		{
			Name: "start",
			ID:   7,
		},
		{
			Name: "home",
			ID:   9,
		},
		{
			Name: "right_stick",
			ID:   10,
		},
		{
			Name: "left_stick",
			ID:   8,
		},
	},
	Hats: []hat{
		{
			Hat:  0,
			Name: "down",
			ID:   4,
		},
		{
			Hat:  0,
			Name: "up",
			ID:   1,
		},
		{
			Hat:  0,
			Name: "left",
			ID:   8,
		},
		{
			Hat:  0,
			Name: "right",
			ID:   2,
		},
		{
			Hat:  0,
			Name: "released",
			ID:   0,
		},
	},
}
