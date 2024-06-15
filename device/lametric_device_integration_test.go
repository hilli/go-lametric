//go:build integration
// +build integration

package device

import (
	"log"
	"os"
	"testing"
)

var lmd *LaMetricDevice

func TestMain(m *testing.M) {
	hostname := os.Getenv("LAMETRIC_HOSTNAME")
	api_key := os.Getenv("LAMETRIC_API_KEY")
	if hostname == "" {
		log.Println("LAMETRIC_HOSTNAME not set for the integration test")
		os.Exit(1)
	}
	if api_key == "" {
		log.Println("LAMETRIC_API_KEY not set for the integration test")
		os.Exit(1)
	}

	lmd = NewDevice(hostname, api_key)

	// Save status before running tests
	_ = lmd.GetStatus()
	pre_status := lmd.Status

	// Run tests
	m.Run()

	// Restore status of the device after tests
	if pre_status.Display.On {
		_ = lmd.On()
	} else {
		_ = lmd.Off()
	}
	lmd.SetBrightness(pre_status.Display.Brightness)
}

func Test_Status(t *testing.T) {
	err := lmd.GetStatus()
	if err != nil {
		t.Errorf("Error getting status: %s", err)
	}
}

func Test_On(t *testing.T) {
	err := lmd.On()
	if err != nil {
		t.Errorf("Error turning on device: %s", err)
	}
	if lmd.Status.Display.On != true {
		t.Errorf("Expected display to be on, got %v", lmd.Status.Display.On)
	}
}

func Test_Off(t *testing.T) {
	err := lmd.Off()
	if err != nil {
		t.Errorf("Error turning off device: %s", err)
	}
}

func Test_On_Again(t *testing.T) {
	err := lmd.On()
	if err != nil {
		t.Errorf("Error turning on device: %s", err)
	}
}

func Test_Display_Brightness(t *testing.T) {
	expectedBrightness := 50
	err := lmd.SetBrightness(expectedBrightness)
	if err != nil {
		t.Errorf("Error setting brightness: %s", err)
	}
	if lmd.Status.Display.Brightness != expectedBrightness {
		t.Errorf("Expected brightness to be 50, got %d", lmd.Status.Display.Brightness)
	}
}
