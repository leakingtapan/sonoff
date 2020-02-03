package types

type Action string

const (
	Register Action = "register"
	Date     Action = "date"
	Query    Action = "query"
	Update   Action = "update"
)

type Message struct {
	Action     Action `json:"action"`
	DeviceId   string `json:"deviceid"`
	UserAgent  string `json:"userAgent"`
	ApiKey     string `json:"apikey"`
	Version    int    `json:"version"`
	RomVersion string `json:"romVersion"`
	Model      string `json:"model"`
	Ts         int64  `json:"ts"`
}

type UpdateParams struct {
	Switch string `json:"switch"`
}

type UpdateMessage struct {
	Message
	Params   UpdateParams `json:"params"`
	Sequence string       `json:"sequence"`
}

type QueryMessage struct {
	Message
	Params []string `json:"params"`
}

type Device struct {
	DeviceId   string `json:"deviceid"`
	ApiKey     string `json:"apikey"`
	Version    int    `json:"version"`
	RomVersion string `json:"romVersion"`
	Model      string `json:"model"`

	// the state of the switch as "on" or "off"
	State string `json:"state"`
}
