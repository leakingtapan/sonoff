package server

import (
	"github.com/leakingtapan/sonoff/pkg/dispatch"
	"github.com/leakingtapan/sonoff/pkg/server"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	serverCmd := serverCmd{}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "start the sonoff backend server",
		RunE:  serverCmd.Run,
	}

	cmd.Flags().StringVar(&serverCmd.serverIp, "server-ip", "", "the IP address of the server")
	cmd.Flags().IntVar(&serverCmd.serverPort, "server-port", 8443, "the port of the server (default to 8443)")
	cmd.Flags().IntVar(&serverCmd.websocketPort, "websocket-port", 1443, "the websocket port of the server (default to 1443")

	return cmd
}

type serverCmd struct {
	serverIp      string
	serverPort    int
	websocketPort int
}

func (c *serverCmd) Run(cmd *cobra.Command, args []string) error {
	ds := server.NewDeviceStore()

	wsServie := server.NewWsServer(c.websocketPort, ds)
	go wsServie.Serve()

	dispatchServer := dispatch.NewDispatchServer(80, c.serverIp, c.websocketPort)
	go dispatchServer.Serve()
	go dispatchServer.ServeHTTPS()

	deviceService := server.NewDeviceService(c.serverIp, c.serverPort, c.websocketPort, ds)
	deviceService.ServeHTTPS()

	return nil
}
