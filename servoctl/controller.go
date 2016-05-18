package servoctl

import (
	"errors"
	"fmt"
	"time"

	"github.com/kidoman/embd/motion/servo"
)

//Controller definition
type Controller struct {
	Settings CtlSettings
	Servo    *servo.Servo
}

//CtlSettings - controller settings
type CtlSettings struct {
	Value int
	Step  int
	Min   int
	Max   int
}

//NewController constructor
func NewController(servo *servo.Servo, settings CtlSettings) *Controller {
	s := &Controller{
		Servo:    servo,
		Settings: settings,
	}
	s.SetSlow(settings.Value, 20)
	return s
}

//Inc increments servo value
func (c *Controller) Inc() error {
	if c.Settings.Value < c.Settings.Max {
		c.Settings.Value += c.Settings.Step
		c.Set(c.Settings.Value)
		return nil
	}
	return errors.New("Unable to increase value")
}

//Dec decrements servo value
func (c *Controller) Dec() error {
	if c.Settings.Value > c.Settings.Min {
		c.Settings.Value -= c.Settings.Step
		c.Set(c.Settings.Value)
		return nil
	}
	return errors.New("Unable to decrease value")
}

//Set sets servo value
func (c *Controller) Set(value int) error {
	if value >= c.Settings.Min && value <= c.Settings.Max {
		c.Settings.Value = value
		fmt.Println("Value: ", value)
		c.Servo.SetAngle(c.Settings.Value)
		return nil
	}
	return errors.New("Unable to set value")
}

//SetSlow sets value with delay
func (c *Controller) SetSlow(value, delay int) error {
	if value > c.Settings.Value {
		for x := c.Settings.Value; x <= value; x++ {
			err := c.Set(x)
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
		return nil
	}
	if value < c.Settings.Value {
		for x := c.Settings.Value; x >= value; x-- {
			err := c.Set(x)
			if err != nil {
				return err
			}
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
		return nil
	}
	return errors.New("Unable to set value")
}
