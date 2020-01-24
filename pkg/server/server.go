package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DeviceHandler struct {
	ip            string
	webSocketPort int
}

func (h *DeviceHandler) Root(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQ | %s | %s ", r.Method, r.URL)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *DeviceHandler) GetDevices(w http.ResponseWriter, r *http.Request) {

}

func (h *DeviceHandler) GetDeviceState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(string(vars["deviceId"]))
	w.Write([]byte("on"))
}

func (h *DeviceHandler) SetDeviceState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println(string(vars["deviceId"]))
	log.Println(string(vars["state"]))
}

func (h *DeviceHandler) GetDeivce(w http.ResponseWriter, r *http.Request) {
}

func (h *DeviceHandler) DispatchDeivce(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQ | %s | %s ", r.Method, r.URL)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read payload: %s", err)
		return
	}
	log.Printf("REQ | %s", string(body))

	resp := struct {
		Err    int    `json:"error"`
		Reason string `json:"reason"`
		Ip     string `json:"IP"`
		Port   int    `json:"port"`
	}{
		Err:    0,
		Reason: "ok",
		Ip:     h.ip,
		Port:   h.webSocketPort,
	}

	output, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response: %s", err)
	}
	log.Println(string(output))
	w.Write(output)
}

func (h *DeviceHandler) SetRoutes(r *mux.Router) {
	r.HandleFunc("/", h.Root).Methods(http.MethodGet)
	r.HandleFunc("/devices", h.GetDevices).Methods(http.MethodGet)
	r.HandleFunc("/devices/{deviceId}/status", h.GetDeviceState).Methods(http.MethodGet)
	r.HandleFunc("/devices/{deviceId}/{state}", h.SetDeviceState).Methods(http.MethodPost)
	r.HandleFunc("/dispatch/device", h.DispatchDeivce).Methods(http.MethodPost)
}

type DeviceService struct {
	svr           *http.Server
	serviceIp     string
	websocketPort int
	devices       *Devices
}

func NewDeviceService(serviceIp string, websocketPort int, devices *Devices) *DeviceService {
	return &DeviceService{
		serviceIp:     serviceIp,
		websocketPort: websocketPort,
		devices:       devices,
	}
}

func (s *DeviceService) Serve() {
	deviceHandler := &DeviceHandler{
		ip:            s.serviceIp,
		webSocketPort: s.websocketPort,
	}
	r := mux.NewRouter()
	deviceHandler.SetRoutes(r)

	svr := http.Server{
		Addr:         ":8443",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Fatal(svr.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"))
}
