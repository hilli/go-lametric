package main

import (
	"fmt"
	"os"
	"time"

	lmd "github.com/hilli/go-lametric/device"
)

func main() {
	// Create a new laMetricDevice
	laMetricDevice := lmd.NewDevice(os.Getenv("LAMETRIC_HOSTNAME"), os.Getenv("LAMETRIC_API_KEY"))

	// Turn off the device
	err := laMetricDevice.Off()
	if err != nil {
		fmt.Println(err)
	}

	// Wait 2 seconds
	time.Sleep(2 * time.Second)

	// Turn on the device
	err = laMetricDevice.On()
	if err != nil {
		fmt.Println(err)
	}

	// Set the brightness to 50
	err = laMetricDevice.SetBrightness(50)
	if err != nil {
		fmt.Println(err)
	}
}
