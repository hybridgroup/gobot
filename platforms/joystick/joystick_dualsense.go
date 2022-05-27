package joystick

var dualsenseConfig = joystickConfig{
	Name: "Dualsense Controller",
	GUID: "E7D56FCA-A01F-4A14-B0D0-4FDAFD847E5E",
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
			Name: "triangle",
			ID:   3,
		},
		{
			Name: "square",
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
			Name: "create",
			ID:   4,
		},
		{
			Name: "options",
			ID:   6,
		},
		{
			Name: "ps",
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
			Name: "l1",
			ID:   9,
		},
		{
			Name: "r1",
			ID:   10,
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
			Name: "trackpad",
			ID:   15,
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
