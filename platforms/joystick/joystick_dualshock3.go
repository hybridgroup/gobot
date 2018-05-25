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
			ID:   3,
		},
		pair{
			Name: "right_y",
			ID:   4,
		},
	},
	Buttons: []pair{
		pair{
			Name: "square",
			ID:   3,
		},
		pair{
			Name: "triangle",
			ID:   2,
		},
		pair{
			Name: "circle",
			ID:   1,
		},
		pair{
			Name: "x",
			ID:   0,
		},
		pair{
			Name: "up",
			ID:   13,
		},
		pair{
			Name: "down",
			ID:   14,
		},
		pair{
			Name: "left",
			ID:   15,
		},
		pair{
			Name: "right",
			ID:   16,
		},
		pair{
			Name: "l1",
			ID:   4,
		},
		{
			Name: "l2",
			ID:   6,
		},
		pair{
			Name: "r1",
			ID:   5,
		},
		pair{
			Name: "r2",
			ID:   7,
		},
		pair{
			Name: "start",
			ID:   9,
		},
		pair{
			Name: "select",
			ID:   8,
		},
		pair{
			Name: "home",
			ID:   10,
		},
	},
	Hats: []hat{},
}
