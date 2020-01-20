package main

import (
	"github.com/leakingtapan/sonoff/pkg/server"
)

func main() {
	websocketPort := 1443
	wsServie := server.NewWsServer(websocketPort)
	go wsServie.Serve()

	servcieIp := "192.168.31.110"
	deviceService := server.NewDeviceService(serviceIp, websocketPort)
	deviceService.Serve()
}
