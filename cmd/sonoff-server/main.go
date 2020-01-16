package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/leakingtapan/sonoff/pkg/server"
)

func main() {
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "/, %q", html.EscapeString(r.URL.Path))
	//	body, _ := ioutil.ReadAll(r.Body)
	//	log.Println(string(body))
	//})
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "ws, %q", html.EscapeString(r.URL.Path))
	//	body, _ := ioutil.ReadAll(r.Body)
	//	log.Println(string(body))
	//})

	//log.Fatal(http.ListenAndServe(":8080", nil))
	deviceSvc := &server.DeviceHandler{}
	r := mux.NewRouter()

	deviceSvc.SetRoutes(r)
	svr := http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Fatal(svr.ListenAndServe())
}
