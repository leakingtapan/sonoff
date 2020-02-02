package switchdev

import (
	"context"

	"github.com/leakingtapan/sonoff/pkg/device"
	"github.com/spf13/cobra"
)

func NewSwitchCommand() *cobra.Command {
	switchCmd := switchCmd{}
	cmd := &cobra.Command{
		Use:   "switch",
		Short: "start the sonoff switch",
		RunE:  switchCmd.Run,
	}

	cmd.Flags().StringVar(&switchCmd.serverIp, "server-ip", "50.18.84.251", "the IP address of the server")
	cmd.Flags().IntVar(&switchCmd.websocketPort, "websocket-port", 443, "the websocket port of the server")

	return cmd
}

type switchCmd struct {
	serverIp      string
	websocketPort int
}

func (c *switchCmd) Run(cmd *cobra.Command, args []string) error {
	sw := device.NewSonoffSwitch(c.serverIp, c.websocketPort)
	ctx := context.Background()

	err := sw.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
