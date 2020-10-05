package server

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/leakingtapan/sonoff/pkg/types"
)

type Device struct {
	types.Device
	Conn *websocket.Conn `json:"-"`
}

type Devices struct {
	mu      sync.Mutex
	devices map[string]*Device
}

func NewDeviceStore() *Devices {
	return &Devices{
		devices: map[string]*Device{},
	}
}

//TODO: move to upper layer
func (ds *Devices) TurnOn(id string) (*types.Device, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	d, found := ds.devices[id]
	if !found {
		return nil, fmt.Errorf("Device %s is not found", id)
	}

	d.State = "on"

	err := pushMessage(d)
	if err != nil {
		return nil, err
	}

	return &d.Device, nil
}

func (ds *Devices) TurnOff(id string) (*types.Device, error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	d, found := ds.devices[id]
	if !found {
		return nil, fmt.Errorf("Device %s is not found", id)
	}

	d.State = "off"

	err := pushMessage(d)
	if err != nil {
		return nil, err
	}
	return &d.Device, nil
}

func (ds *Devices) Get(id string) (*Device, bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	d, found := ds.devices[id]
	return d, found
}

func (ds *Devices) ListDevices() []*Device {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	var res []*Device
	for _, d := range ds.devices {
		res = append(res, d)
	}
	return res
}

// TODO: improve
func (ds *Devices) AddOrUpdateDevice(d *Device) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.devices[d.DeviceId] = d
}

func pushMessage(d *Device) error {
	resp := struct {
		ApiKey   string `json:"apikey"`
		Action   string `json:"action"`
		DeviceId string `json:"deviceid"`
		Params   struct {
			Switch string `json:"switch"`
		} `json:"params"`
		UserAgent string `json:"userAgent"`
		Sequence  string `json:"sequence"`
		Ts        int    `json:"ts"`
		From      string `json:"from"`
	}{
		ApiKey:   "111111111-1111-1111-1111-111111111111",
		Action:   "update",
		DeviceId: d.DeviceId,
		Params: struct {
			Switch string `json:"switch"`
		}{
			Switch: d.State,
		},
		UserAgent: "app",
		Sequence:  time.Now().Format("2006-01-02T15:04:05Z"),
		Ts:        0,
		From:      "app",
	}

	payload, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	log.Printf("REQ | WS | APP | %s", string(payload))
	return d.Conn.WriteMessage(websocket.TextMessage, payload)
}
