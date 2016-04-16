package main

import (
	"errors"
	"fmt"

	"github.com/kidoman/embd/motion/servo"
)

type Controller struct {
	Settings CtlSettings
	Servo    *servo.Servo
}

type CtlSettings struct {
	Value int
	Step  int
	Min   int
	Max   int
}

func NewController(servo *servo.Servo, settings CtlSettings) *Controller {
	s := &Controller{
		Servo:    servo,
		Settings: settings,
	}
	return s
}

func (c *Controller) Inc() error {
	if c.Settings.Value < c.Settings.Max {
		c.Settings.Value += c.Settings.Step
		c.Set(c.Settings.Value)
		return nil
	}
	return errors.New("Unable to increase value")
}

func (c *Controller) Dec() error {
	if c.Settings.Value > c.Settings.Min {
		c.Settings.Value -= c.Settings.Step
		c.Set(c.Settings.Value)
		return nil
	}
	return errors.New("Unable to decrease value")
}

func (c *Controller) Set(value int) error {
	if c.Settings.Value >= c.Settings.Min && c.Settings.Value <= c.Settings.Max {
		c.Settings.Value = value
		fmt.Println("Value: ", value)
		c.Servo.SetAngle(c.Settings.Value)
		return nil
	}
	return errors.New("Unable to set value")
}
