package gobot

type Master struct {
	Robots []Robot
}

func GobotMaster() *Master {
	m := new(Master)
	return m
}

func (m *Master) Start() {
	for s := range m.Robots {
		go m.Robots[s].Start()
	}
	select {}
}

func (m *Master) FindRobot(name string) *Robot {
	for s := range m.Robots {
		if m.Robots[s].Name == name {
			return &m.Robots[s]
		}
	}
	return nil
}

func (m *Master) FindRobotDevice(name string, device string) *device {
	for r := range m.Robots {
		if m.Robots[r].Name == name {
			for d := range m.Robots[r].devices {
				if m.Robots[r].devices[d].Name == device {
					return m.Robots[r].devices[d]
				}
			}
		}
	}
	return nil
}

func (m *Master) FindRobotConnection(name string, connection string) *connection {
	for r := range m.Robots {
		if m.Robots[r].Name == name {
			for c := range m.Robots[r].connections {
				if m.Robots[r].connections[c].Name == connection {
					return m.Robots[r].connections[c]
				}
			}
		}
	}
	return nil
}
