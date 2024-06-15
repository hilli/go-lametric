package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	lmd "github.com/hilli/go-lametric/device"
)

func main() {
	// Create a new laMetricDevice
	laMetricDevice := lmd.NewDevice(os.Getenv("LAMETRIC_HOSTNAME"), os.Getenv("LAMETRIC_API_KEY"))

	// Get the name of the device
	deviceName := laMetricDevice.Status.Name

	// Create the light accessory
	LMdevice := accessory.NewLightbulb(accessory.Info{Name: deviceName, Model: laMetricDevice.Status.Model, Manufacturer: "LaMetric"})
	LMdevice.Lightbulb.On.SetValue(laMetricDevice.Status.Display.On) // Set to the current state of the light

	// Add brightness to the light accessory
	brightness := characteristic.NewBrightness()
	brightness.SetValue(laMetricDevice.Status.Display.Brightness)
	brightness.SetMinValue(laMetricDevice.Status.Display.BrightnessLimit.Min)
	brightness.SetMaxValue(laMetricDevice.Status.Display.BrightnessLimit.Max)
	brightness.SetStepValue(1)
	brightness.OnSetRemoteValue(func(value int) error {
		err := laMetricDevice.SetBrightness(value)
		log.Printf("%s brightness changed to: %d", laMetricDevice.Status.Name, value)
		return err
	})
	LMdevice.Lightbulb.Cs = append(LMdevice.Lightbulb.Cs, brightness.C)

	// Handle the on/off state of the light
	LMdevice.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on {
			log.Printf("%s switched on\n", laMetricDevice.Status.Name)
			laMetricDevice.On()
		} else {
			log.Printf("%s switched off\n", laMetricDevice.Status.Name)
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
