package joystick

var dualshock4Config = joystickConfig{
	Name: "Dualshock4 Controller",
	GUID: "2222",
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
			Name: "l2",
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
			ID:   0,
		},
		{
			Name: "triangle",
			ID:   3,
		},
		{
			Name: "circle",
			ID:   2,
		},
		{
			Name: "x",
			ID:   1,
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
			Name: "share",
			ID:   13,
		},
		{
			Name: "options",
			ID:   9,
		},
		{
			Name: "home",
			ID:   12,
		},
		{
			Name: "left_joystick",
			ID:   10,
		},
		{
			Name: "right_joystick",
			ID:   11,
		},
		{
			Name: "panel",
			ID:   13,
		},
	},
}
