package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	lmd "github.com/hilli/go-lametric/device"
)

func main() {
	// Create a new laMetricDevice
	laMetricDevice := lmd.NewDevice(os.Getenv("LAMETRIC_HOSTNAME"), os.Getenv("LAMETRIC_API_KEY"))

	// Get the name of the device
	deviceName := laMetricDevice.Status.Name

	// Create the light accessory
	LMdevice := accessory.NewLightbulb(accessory.Info{Name: deviceName, Model: laMetricDevice.Status.Model, Manufacturer: "LaMetric"})

	LMdevice.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			log.Println("Light is on")
			laMetricDevice.On()
		} else {
			log.Println("Light is off")
			laMetricDevice.Off()
		}
	})

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore("./db")

	// Create the hap server.
	server, err := hap.NewServer(fs, LMdevice.A)
	if err != nil {
		// stop if an error happens
		log.Panic(err)
	}

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	server.ListenAndServe(ctx)

}
