package joystick

var dualshock3Config = joystickConfig{
	Name: "Dualshock3 Controller",
	GUID: "1111",
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
	},
	Buttons: []pair{
		{
			Name: "square",
			ID:   15,
		},
		{
			Name: "triangle",
			ID:   12,
		},
		{
			Name: "circle",
			ID:   13,
		},
		{
			Name: "x",
			ID:   14,
		},
		{
			Name: "up",
			ID:   4,
		},
		{
			Name: "down",
			ID:   6,
		},
		{
			Name: "left",
			ID:   17,
		},
		{
			Name: "right",
			ID:   5,
		},
		{
			Name: "l1",
			ID:   10,
		},
		{
			Name: "l2",
			ID:   8,
		},
		{
			Name: "r1",
			ID:   11,
		},
		{
			Name: "r2",
			ID:   9,
		},
		{
			Name: "start",
			ID:   3,
		},
		{
			Name: "select",
			ID:   0,
		},
		{
			Name: "home",
			ID:   16,
		},
	},
}
