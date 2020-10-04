package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DeviceHandler struct {
	ip            string
	webSocketPort int
	devices       *Devices
}

func (h *DeviceHandler) Root(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQ | %s | %s ", r.Method, r.URL)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices := h.devices.ListDevices()
	resp, err := json.Marshal(&devices)
	if err != nil {
		msg := fmt.Sprintf("Failed to marshal devices: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *DeviceHandler) GetDeviceState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]

	device, found := h.devices.Get(deviceId)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
		return
	}

	resp, err := json.Marshal(&device)
	if err != nil {
		msg := fmt.Sprintf("Failed to marshal device: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *DeviceHandler) SetDeviceState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	state := vars["state"]
	var err error

	switch state {
	case "on":
		err = h.devices.TurnOn(deviceId)
	case "off":
		err = h.devices.TurnOff(deviceId)
	}

	if err != nil {
		msg := fmt.Sprintf("Failed to set device %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	w.Write([]byte("OK"))
}

func (h *DeviceHandler) GetDeivce(w http.ResponseWriter, r *http.Request) {
}

func (h *DeviceHandler) SetRoutes(r *mux.Router) {
	r.HandleFunc("/", h.Root).Methods(http.MethodGet)
	r.HandleFunc("/devices", h.GetDevices).Methods(http.MethodGet)
	r.HandleFunc("/devices/{deviceId}/status", h.GetDeviceState).Methods(http.MethodGet)
	r.HandleFunc("/devices/{deviceId}/{state}", h.SetDeviceState).Methods(http.MethodPost)
	//r.HandleFunc("/dispatch/device", h.DispatchDeivce).Methods(http.MethodPost)
}

type DeviceService struct {
	svr           *http.Server
	serviceIp     string
	servicePort   int
	websocketPort int
	devices       *Devices
}

func NewDeviceService(serviceIp string, servicePort int, websocketPort int, devices *Devices) *DeviceService {
	return &DeviceService{
		serviceIp:     serviceIp,
		servicePort:   servicePort,
		websocketPort: websocketPort,
		devices:       devices,
	}
}

func (s *DeviceService) ServeHTTPS() {
	deviceHandler := &DeviceHandler{
		ip:            s.serviceIp,
		webSocketPort: s.websocketPort,
		devices:       s.devices,
	}
	r := mux.NewRouter()
	deviceHandler.SetRoutes(r)

	addr := fmt.Sprintf("%s:%d", "", s.servicePort)
	svr := http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Fatal(svr.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"))
}
