package server

type Device struct {
}
type Devices struct {
	devices map[string]*Devices
}

func (d *Devices) TurnOn(id string) {
}

func (d *Devices) TurnOff(id string) {
}

func (d *Devices) Get(id string) *Device {
	return nil
}

func (d *Devices) ListConnected() []*Device {
	return []*Device{}
}
