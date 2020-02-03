package switchdev

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/leakingtapan/sonoff/pkg/device"
	"github.com/leakingtapan/sonoff/pkg/types"
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
	cmd.Flags().StringVar(&switchCmd.deviceSpecPath, "device-spec-path", "", "the path to the device spec json file")

	cmd.MarkFlagRequired("device-spec-path")

	return cmd
}

type switchCmd struct {
	serverIp          string
	websocketServerIp string
	websocketPort     int
	deviceSpecPath    string
}

func (c *switchCmd) Run(cmd *cobra.Command, args []string) error {
	specPath, err := filepath.Abs(filepath.Clean(c.deviceSpecPath))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(specPath)
	if err != nil {
		return err
	}

	var spec types.Device
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return err
	}

	ctx := context.Background()
	sw := device.NewSonoffSwitch(
		c.serverIp,
		c.websocketServerIp,
		c.websocketPort,
		spec,
	)
	err = sw.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
