package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type DeviceHandler struct {
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

func (h *DeviceHandler) SetRoutes(r *mux.Router) {
	r.HandleFunc("/devices/{deviceId}/status", h.GetDeviceState).
		Methods(http.MethodGet)
	r.HandleFunc("/devices/{deviceId}/{state}", h.SetDeviceState).
		Methods(http.MethodPost)
}
