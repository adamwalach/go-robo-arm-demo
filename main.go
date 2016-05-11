package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	ctl "github.com/adamwalach/go-robo-arm-demo/servoctl"
	"github.com/gorilla/mux"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/motion/servo"
)

var (
	vertCtl *ctl.Controller
	horCtl  *ctl.Controller
	gripCtl *ctl.Controller
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"handler":         "urlHandler",
		"query":           r.URL.Query(),
		"ip":              r.RemoteAddr,
		"x-forwarded-for": r.Header.Get("X-Forwarded-For"),
	}).Info("Scripts request")

	op := r.FormValue("operation")
	servo := r.FormValue("servo")

	output := ""

	if op != "" {
		switch servo {
		case "v":
			if op == "i" {
				vertCtl.Inc()
			} else {
				vertCtl.Dec()
			}
		case "h":
			if op == "i" {
				horCtl.Inc()
			} else {
				horCtl.Dec()
			}
		case "g":
			if op == "i" {
				gripCtl.Inc()
			} else {
				gripCtl.Dec()
			}
		default:
			output = "Error"
		}
	} else {
		v, err := strconv.Atoi(r.FormValue("value"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch servo {
		case "v":
			vertCtl.Set(v)
		case "h":
			horCtl.Set(v)
		case "g":
			gripCtl.Set(v)
		default:
			output = "Error"
		}
	}

	w.Write([]byte(output))
}

// func keyboardHandler() {
// 	for {
// 		ascii, keyCode, _ := keys.GetChar()
// 		fmt.Println("A: ", ascii, "C: ", keyCode)
// 		switch ascii {
// 		case keys.AsciiEsc:
// 			return
// 		case keys.AsciiW:
// 			vertCtl.Inc()
// 		case keys.AsciiS:
// 			vertCtl.Dec()
// 		case keys.AsciiA:
// 			horCtl.Inc()
// 		case keys.AsciiD:
// 			horCtl.Dec()
// 		}
//
// 		switch keyCode {
// 		case keys.CodeUpArrow:
// 			gripCtl.Inc()
// 		case keys.CodeDownArrow:
// 			gripCtl.Dec()
// 		}
// 	}
// }

func demo() {
	vertCtl.SetSlow(vertCtl.Settings.Max, 20)
	vertCtl.SetSlow(vertCtl.Settings.Min, 20)
	vertCtl.SetSlow(vertCtl.Settings.Min+30, 20)

	horCtl.SetSlow(horCtl.Settings.Min, 10)
	horCtl.SetSlow(horCtl.Settings.Max, 10)
	horCtl.SetSlow(90, 10)

	gripCtl.SetSlow(gripCtl.Settings.Min, 5)
	gripCtl.SetSlow(gripCtl.Settings.Max, 5)
	gripCtl.SetSlow(gripCtl.Settings.Min, 5)
	gripCtl.SetSlow(gripCtl.Settings.Max, 5)

	gripCtl.SetSlow(110, 5)
}

func captureCtrlC(d *pca9685.PCA9685) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			fmt.Println(fmt.Sprintf("Captured '%v', exiting..", sig))
			d.Close()
			embd.CloseI2C()
			os.Exit(1)
		}
	}()
}

func main() {

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	d := pca9685.New(bus, 0x40)
	d.Freq = 60
	defer d.Close()
	captureCtrlC(d)

	vertCtl = ctl.NewController(
		servo.New(d.ServoChannel(2)),
		ctl.CtlSettings{
			Value: 1,
			Step:  4,
			Max:   115,
			Min:   1,
		})
	horCtl = ctl.NewController(
		servo.New(d.ServoChannel(3)),
		ctl.CtlSettings{
			Value: 90,
			Step:  4,
			Max:   170,
			Min:   20,
		})
	gripCtl = ctl.NewController(
		servo.New(d.ServoChannel(0)),
		ctl.CtlSettings{
			Value: 110,
			Step:  4,
			Max:   184,
			Min:   105,
		})
	demo()
	r := mux.NewRouter()
	r.HandleFunc("/api", apiHandler).Methods("GET")
	go http.ListenAndServe(":3000", r)

	for {
		time.Sleep(time.Second)
	}
}
