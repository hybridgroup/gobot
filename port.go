package gobot

type Port struct {
	Name string
}

func (Port) NewPort(p string) *Port {
	return new(Port)
}

func (p *Port) ToString() string {
	return p.Name
}
