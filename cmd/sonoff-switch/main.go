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

	cmd.Flags().StringVar(&switchCmd.serverIp, "server-endpoint", "disp.coolkit.cc", "the endpoint of the dispatch server")
	cmd.Flags().StringVar(&switchCmd.websocketServerIp, "websocket-server-ip", "", "the optional IP address of the websocket server")
	cmd.Flags().IntVar(&switchCmd.websocketPort, "websocket-server-port", 0, "the optional port of the websocket server")

	return cmd
}

type switchCmd struct {
	serverIp          string
	websocketServerIp string
	websocketPort     int
}

func (c *switchCmd) Run(cmd *cobra.Command, args []string) error {
	sw := device.NewSonoffSwitch(c.serverIp, c.websocketServerIp, c.websocketPort)
	ctx := context.Background()

	if c.websocketServerIp == "" || c.websocketPort == 0 {
		err := sw.Dispatch()
		if err != nil {
			return err
		}
	}

	err := sw.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
