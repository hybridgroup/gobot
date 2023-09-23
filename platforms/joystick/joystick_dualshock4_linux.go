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
			ID:   3,
		},
		{
			Name: "right_y",
			ID:   4,
		},
		{
			Name: "l2",
			ID:   2,
		},
		{
			Name: "r2",
			ID:   5,
		},
		{
			Name: "right_left",
			ID:   6,
		},
		{
			Name: "up_down",
			ID:   7,
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
			ID:   8,
		},
		{
			Name: "options",
			ID:   9,
		},
		{
			Name: "home",
			ID:   10,
		},
		{
			Name: "left_joystick",
			ID:   11,
		},
		{
			Name: "right_joystick",
			ID:   12,
		},
	},
}
