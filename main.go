package main

import (
	"fmt"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/motion/servo"
)

var (
	servoMin = 400
	servoMax = 600
)

func main() {

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	d := pca9685.New(bus, 0x40)
	d.Freq = 50
	defer d.Close()

	vertCtl := NewController(
		servo.New(d.ServoChannel(2)),
		CtlSettings{
			Value: 90,
			Step:  1,
			Max:   200,
			Min:   20,
		})
	horCtl := NewController(
		servo.New(d.ServoChannel(3)),
		CtlSettings{
			Value: 90,
			Step:  1,
			Max:   200,
			Min:   20,
		})
	gripCtl := NewController(
		servo.New(d.ServoChannel(0)),
		CtlSettings{
			Value: 90,
			Step:  1,
			Max:   200,
			Min:   20,
		})

	for {
		ascii, keyCode, _ := getChar()
		fmt.Println("A: ", ascii, "C: ", keyCode)
		switch ascii {
		case AsciiEsc:
			return
		case AsciiW:
			vertCtl.Inc()
		case AsciiS:
			vertCtl.Dec()
		case AsciiA:
			horCtl.Dec()
		case AsciiD:
			horCtl.Inc()
		}

		switch ascii {
		case CodeUpArrow:
			gripCtl.Inc()
		case CodeDownArrow:
			gripCtl.Dec()
		}

	}
}
