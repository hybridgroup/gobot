package joystick

var dualshock3Config = joystickConfig{
	Name: "Dualshock3 Controller",
	GUID: "1111",
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
	},
	Buttons: []pair{
		pair{
			Name: "square",
			ID:   15,
		},
		pair{
			Name: "triangle",
			ID:   12,
		},
		pair{
			Name: "circle",
			ID:   13,
		},
		pair{
			Name: "x",
			ID:   14,
		},
		pair{
			Name: "up",
			ID:   4,
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
			ID:   5,
		},
		pair{
			Name: "left_stick",
			ID:   1,
		},
		pair{
			Name: "right_stick",
			ID:   2,
		},
		pair{
			Name: "l1",
			ID:   10,
		},
		{
			Name: "l2",
			ID:   8,
		},
		pair{
			Name: "r1",
			ID:   11,
		},
		pair{
			Name: "r2",
			ID:   9,
		},
		pair{
			Name: "start",
			ID:   3,
		},
		pair{
			Name: "select",
			ID:   0,
		},
		pair{
			Name: "home",
			ID:   16,
		},
	},
	Hats: []hat{},
}
