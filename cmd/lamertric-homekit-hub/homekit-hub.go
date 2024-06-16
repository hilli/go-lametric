package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	lmd "github.com/hilli/go-lametric/device"
	myDIY "github.com/hilli/lametric-my-data-diy-go"
)

func main() {
	// Create a new laMetricDevice
	hostname := os.Getenv("LAMETRIC_HOSTNAME")
	api_key := os.Getenv("LAMETRIC_API_KEY")

	if hostname == "" || api_key == "" {
		log.Println("LAMETRIC_HOSTNAME or LAMETRIC_API_KEY not set - They shoud be set in the environment variables")
		os.Exit(1)
	}

	laMetricDevice := lmd.NewDevice(hostname, api_key)

	// Send notifications to the device about brightness levels
	diyPushURL := os.Getenv("LAMETRIC_DIY_PUSH_URL")

	// Get the name of the device
	deviceName := laMetricDevice.Status.Name

	// Create the light accessory
	LMdevice := accessory.NewLightbulb(accessory.Info{
		Name:         deviceName,
		Model:        laMetricDevice.Status.Model,
		Manufacturer: "LaMetric (github.com/hilli/go-lametric)",
		SerialNumber: laMetricDevice.Status.SerialNumber,
	})

	// Set to the current state of the light
	LMdevice.Lightbulb.On.SetValue(laMetricDevice.Status.Display.On)

	// Add brightness to the light accessory
	brightness := characteristic.NewBrightness()
	brightness.SetValue(laMetricDevice.Status.Display.Brightness)
	brightness.SetMinValue(laMetricDevice.Status.Display.BrightnessLimit.Min)
	brightness.SetMaxValue(laMetricDevice.Status.Display.BrightnessLimit.Max)
	brightness.SetStepValue(1)
	brightness.OnSetRemoteValue(func(value int) error {
		err := laMetricDevice.SetBrightness(value)
		pushBrightness(diyPushURL, laMetricDevice.API_KEY, value)
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

func pushBrightness(url, api_key string, brightness int) {
	if url == "" {
		return
	}
	// Setup frame
	frame := myDIY.MyDataFrame{
		Text: fmt.Sprintf("%d%% light", brightness),
		Icon: "22581",
	}
	frames := myDIY.MyDataFrames{}
	frames.AddFrame(frame)

	// Send the data to the device
	err := frames.Push(url, api_key)
	if err != nil {
		log.Printf("Error pushing data: %s", err)
	}
}
