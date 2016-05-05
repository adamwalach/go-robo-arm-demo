package main

import (
	"fmt"

	ctl "github.com/adamwalach/go-robo-arm-demo/servoctl"
	keys "github.com/adamwalach/go-robo-arm-demo/servoctl"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/motion/servo"
)

func main() {

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	d := pca9685.New(bus, 0x40)
	d.Freq = 60
	defer d.Close()

	vertCtl := ctl.NewController(
		servo.New(d.ServoChannel(2)),
		ctl.CtlSettings{
			Value: 85,
			Step:  1,
			Max:   175,
			Min:   1,
		})
	horCtl := ctl.NewController(
		servo.New(d.ServoChannel(3)),
		ctl.CtlSettings{
			Value: 95,
			Step:  1,
			Max:   200,
			Min:   20,
		})
	gripCtl := ctl.NewController(
		servo.New(d.ServoChannel(0)),
		ctl.CtlSettings{
			Value: 110,
			Step:  1,
			Max:   184,
			Min:   105,
		})

	for {
		ascii, keyCode, _ := keys.GetChar()
		fmt.Println("A: ", ascii, "C: ", keyCode)
		switch ascii {
		case keys.AsciiEsc:
			return
		case keys.AsciiW:
			vertCtl.Inc()
		case keys.AsciiS:
			vertCtl.Dec()
		case keys.AsciiA:
			horCtl.Inc()
		case keys.AsciiD:
			horCtl.Dec()
		}

		switch keyCode {
		case keys.CodeUpArrow:
			gripCtl.Inc()
		case keys.CodeDownArrow:
			gripCtl.Dec()
		}

	}
}
