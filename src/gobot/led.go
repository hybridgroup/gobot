package gobot

type Led struct {
  Driver
  //Beaglebone *AdaptorInterface
  Beaglebone *Beaglebone
  High bool
}

func NewLed(b *Beaglebone) *Led{
  l := Led{High: false, Beaglebone: b}
  return &l
}

func (l *Led) IsOn() bool {
  return l.High
}

func (l *Led) IsOff() bool {
  return !l.IsOn()
}

func (l *Led) On() bool {
  l.changeState(l.Pin, "1")
  l.High = true
  return true
}

func (l *Led) Off() bool {
  l.changeState(l.Pin, "0")
  l.High = false
  return true
}

func (l *Led) Toggle() {
  if l.IsOn() {
    l.Off()
  } else {
    l.On() 
  }
}

func (l *Led) changeState(pin string, level string) {
  l.Beaglebone.DigitalWrite(pin, level)
}
