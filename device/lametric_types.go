package device

// Status is the status of the LaMetric device from the device API
// Don't change values trough this struct
type Status struct {
	Audio        Audio     `json:"audio"`
	Bluetooth    Bluetooth `json:"bluetooth"`
	Display      Display   `json:"display"`
	Id           string    `json:"id"`
	Mode         string    `json:"mode"`
	Model        string    `json:"model"`
	Name         string    `json:"name"`
	OSVersion    string    `json:"os_version"`
	SerialNumber string    `json:"serial_number"`
	WiFi         WiFi      `json:"wifi"`
}

type Audio struct {
	Available bool `json:"available"`
}

type Bluetooth struct {
	Available bool `json:"available"`
}

type Display struct {
	Brightness      int             `json:"brightness"`
	BrightnessLimit BrightnessLimit `json:"brightness_limit"`
	BrightnessMode  string          `json:"brightness_mode"`
	BrightnessRange BrightnessRange `json:"brightness_tange"`
	Height          int             `json:"height"`
	On              bool            `json:"on"`
	Type            string          `json:"type"`
	Width           int             `json:"width"`
}

type BrightnessLimit struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type BrightnessRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type WiFi struct {
	Active    bool   `json:"active"`
	Address   string `json:"address"`
	Available bool   `json:"available"`
	ESSID     string `json:"essid"`
	IP        string `json:"ip"`
	Mode      string `json:"mode"`
	Netmask   string `json:"netmask"`
	Strength  int    `json:"strength"`
}
