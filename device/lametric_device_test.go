package device

import (
	"testing"
)

func Test_New(t *testing.T) {
	hostname := "lametric.local"
	api_key := "123"

	lmd := NewDevice(hostname, api_key)
	if lmd.hostname != hostname {
		t.Errorf("Expected hostname to be as set in env var, got %s", lmd.hostname)
	}
	if lmd.APIURL != "https://lametric.local:4343/api/v2" {
		t.Errorf("Expected APIURL to be 'https://lametric.local:4343/api/v2', got %s", lmd.APIURL)
	}
}
