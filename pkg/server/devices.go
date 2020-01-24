package server

import "github.com/gorilla/websocket"

type Device struct {
	Id      string
	Version int
	Model   string
	Conn    *websocket.Conn
}

// TODO: add mu
type Devices struct {
	devices map[string]*Device
}

func NewDeviceStore() *Devices {
	return &Devices{
		devices: map[string]*Device{},
	}
}

func (ds *Devices) TurnOn(id string) {
}

func (ds *Devices) TurnOff(id string) {
}

func (ds *Devices) Get(id string) *Device {
	return nil
}

func (ds *Devices) ListConnected() []*Device {
	return []*Device{}
}

func (ds *Devices) AddOrUpdateDevice(d *Device) {
	ds.devices[d.Id] = d
}
