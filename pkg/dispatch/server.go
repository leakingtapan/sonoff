package dispatch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type DispatchHandler struct {
	wsServerIp   string
	wsServerPort int
}

func (h *DispatchHandler) DispatchDeivce(w http.ResponseWriter, r *http.Request) {
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
		Ip:     h.wsServerIp,
		Port:   h.wsServerPort,
	}

	output, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Failed to marshal response: %s", err)
	}
	log.Println(string(output))
	w.Write(output)
}

type DispatchServer struct {
	serverPort   int
	wsServerIp   string
	wsServerPort int
}

func NewDispatchServer(serverPort int, wsServerIp string, wsServerPort int) *DispatchServer {
	return &DispatchServer{
		serverPort:   serverPort,
		wsServerIp:   wsServerIp,
		wsServerPort: wsServerPort,
	}
}

// TODO: configurable ports
func (s *DispatchServer) ServeHTTPS() error {
	svr := s.server(8443)
	return svr.ListenAndServeTLS("./certs/server.crt", "./certs/server.key")
}

// TODO: configurable ports
func (s *DispatchServer) Serve() error {
	svr := s.server(80)
	return svr.ListenAndServe()
}

func (s *DispatchServer) server(port int) http.Server {
	h := DispatchHandler{
		wsServerIp:   s.wsServerIp,
		wsServerPort: s.wsServerPort,
	}
	r := mux.NewRouter()
	r.HandleFunc("/dispatch/device", h.DispatchDeivce).
		Methods(http.MethodPost, http.MethodGet)

	addr := fmt.Sprintf(":%d", port)
	svr := http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	return svr
}
