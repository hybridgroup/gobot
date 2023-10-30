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
			Name: "l2",
			ID:   2,
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
			Name: "r2",
			ID:   5,
		},
	},
	Buttons: []pair{
		{
			Name: "square",
			ID:   3,
		},
		{
			Name: "triangle",
			ID:   2,
		},
		{
			Name: "circle",
			ID:   1,
		},
		{
			Name: "x",
			ID:   0,
		},
		{
			Name: "up",
			ID:   13,
		},
		{
			Name: "down",
			ID:   14,
		},
		{
			Name: "left",
			ID:   15,
		},
		{
			Name: "right",
			ID:   16,
		},
		{
			Name: "l1",
			ID:   4,
		},
		{
			Name: "l2",
			ID:   6,
		},
		{
			Name: "r1",
			ID:   5,
		},
		{
			Name: "r2",
			ID:   7,
		},
		{
			Name: "start",
			ID:   9,
		},
		{
			Name: "select",
			ID:   8,
		},
		{
			Name: "home",
			ID:   10,
		},
	},
}
