package main

import (
	"github.com/leakingtapan/sonoff/pkg/server"
)

func main() {
	ds := server.NewDeviceStore()

	websocketPort := 1443
	wsServie := server.NewWsServer(websocketPort, ds)
	go wsServie.Serve()

	serviceIp := "192.168.31.110"
	deviceService := server.NewDeviceService(serviceIp, websocketPort, ds)
	deviceService.Serve()
}
