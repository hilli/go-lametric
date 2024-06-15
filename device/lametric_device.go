package device

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LaMetricDevice is a client for the LaMetric device
type LaMetricDevice struct {
	hostname   string
	APIURL     string
	api_key    string
	Status     *Status
	httpclient *http.Client
}

const (
	httpTimeout = time.Duration(1) * time.Second
)

// NewDevice creates a new LaMetricDevice client
// hostname is the hostname of the LaMetric device
// api_key is the API key for the device
//
// Example:
//
// lmd := NewDevice("lametric.local", "123")
func NewDevice(hostname, api_key string) *LaMetricDevice {
	APIURL := "https://" + hostname + ":4343/api/v2"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	lmd := &LaMetricDevice{
		hostname:   hostname,
		api_key:    api_key,
		APIURL:     APIURL,
		httpclient: &http.Client{Timeout: httpTimeout, Transport: tr},
	}

	return lmd
}

func (lmd *LaMetricDevice) On() error {
	req, err := http.NewRequest("PUT", lmd.APIURL+"/device/display", bytes.NewBufferString(`{"on": true }`))
	if err != nil {
		return err
	}
	lmd.addRequiredHeaders(req)
	res, err := lmd.httpclient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error turning on device %s: %s", lmd.hostname, res.Status)
	}
	lmd.GetStatus()
	return nil
}

func (lmd *LaMetricDevice) Off() error {
	req, err := http.NewRequest("PUT", lmd.APIURL+"/device/display", bytes.NewBufferString(`{"on": false }`))
	if err != nil {
		return err
	}
	lmd.addRequiredHeaders(req)
	res, err := lmd.httpclient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error turning off device %s: %s", lmd.hostname, res.Status)
	}
	lmd.GetStatus()
	return nil
}

func (lmd *LaMetricDevice) SetBrightness(brightness int) error {
	req, err := http.NewRequest("PUT", lmd.APIURL+"/device/display", bytes.NewBufferString(fmt.Sprintf(`{"brightness": %d }`, brightness)))
	if err != nil {
		return err
	}
	lmd.addRequiredHeaders(req)
	res, err := lmd.httpclient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error setting brightness on device %s: %s", lmd.hostname, res.Status)
	}
	lmd.GetStatus()
	return nil
}

func (lmd *LaMetricDevice) GetStatus() error {
	req, err := http.NewRequest("GET", lmd.APIURL+"/device", nil)
	if err != nil {
		return err
	}
	lmd.addRequiredHeaders(req)
	res, err := lmd.httpclient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("error getting status of device %s: %s", lmd.hostname, res.Status)
	}
	status := &Status{}
	err = json.NewDecoder(res.Body).Decode(status)
	if err != nil {
		return err
	}
	lmd.Status = status
	return nil
}

func (lmd *LaMetricDevice) addRequiredHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth("dev", lmd.api_key)
}
