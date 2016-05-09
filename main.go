package main

import (
	"net/http"
	"strconv"
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

func main() {

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	d := pca9685.New(bus, 0x40)
	d.Freq = 60
	defer d.Close()

	vertCtl = ctl.NewController(
		servo.New(d.ServoChannel(2)),
		ctl.CtlSettings{
			Value: 85,
			Step:  1,
			Max:   175,
			Min:   1,
		})
	horCtl = ctl.NewController(
		servo.New(d.ServoChannel(3)),
		ctl.CtlSettings{
			Value: 95,
			Step:  1,
			Max:   200,
			Min:   20,
		})
	gripCtl = ctl.NewController(
		servo.New(d.ServoChannel(0)),
		ctl.CtlSettings{
			Value: 110,
			Step:  1,
			Max:   184,
			Min:   105,
		})

	r := mux.NewRouter()
	r.HandleFunc("/api", apiHandler).Methods("GET")
	go http.ListenAndServe(":3000", r)

	for {
		time.Sleep(time.Second)
	}
}
