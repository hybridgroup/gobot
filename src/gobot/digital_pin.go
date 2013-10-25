package gobot

import (
  "os"
  "strconv"
)

type DigitalPin struct {
  PinNum string
  Mode string
  PinFile *os.File
  Status string
}

const GPIO_PATH = "/sys/class/gpio"
const GPIO_DIRECTION_READ = "in"
const GPIO_DIRECTION_WRITE = "out"
const HIGH = 1
const LOW = 0

func NewDigitalPin(pinNum int, mode string) *DigitalPin {
  d := DigitalPin{PinNum: strconv.Itoa(pinNum)}

  fi, _ := os.OpenFile(GPIO_PATH + "/export", os.O_WRONLY | os.O_APPEND, 0666)
  fi.WriteString(d.PinNum)
  fi.Close()

  d.SetMode(mode)

  return &d
}

func (d *DigitalPin) SetMode(mode string) {
  d.Mode = mode

  if mode == "w" {
    fi, _ := os.OpenFile(GPIO_PATH + "/gpio" + d.PinNum + "/direction", os.O_WRONLY, 0666)
    fi.WriteString(GPIO_DIRECTION_WRITE)
    fi.Close()
    d.PinFile, _ = os.OpenFile(GPIO_PATH + "/gpio" + d.PinNum + "/value", os.O_WRONLY, 0666)
  } else if mode =="r" {
    fi, _ := os.OpenFile(GPIO_PATH + "/gpio" + d.PinNum + "/direction", os.O_WRONLY, 0666)
    fi.WriteString(GPIO_DIRECTION_READ)
    fi.Close()
    d.PinFile, _ = os.OpenFile(GPIO_PATH + "/gpio" + d.PinNum + "/value", os.O_RDONLY, 0666)
  }
}

func (d *DigitalPin) DigitalWrite(value string) {

  if d.Mode != "w" {
    d.SetMode("w")
  }

  d.PinFile.WriteString(value)
  d.PinFile.Sync()
}

func (d *DigitalPin) IsOn() bool {
  if d.Status == "1" {
    return true
  } else {
    return false
  }
}

func (d *DigitalPin) IsOff() bool {
  return !d.IsOn()
}

func (d *DigitalPin) On() {
  d.DigitalWrite("1")
}

func (d *DigitalPin) Off() {
  d.DigitalWrite("0")
}

func (d *DigitalPin) Close() {
  fi, _ := os.OpenFile(GPIO_PATH + "/unexport", os.O_WRONLY | os.O_APPEND, 0666)
  fi.WriteString(d.PinNum)
  fi.Close()
  d.PinFile.Close()
}
