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

	return cmd
}

type switchCmd struct {
}

func (c *switchCmd) Run(cmd *cobra.Command, args []string) error {
	sw := device.NewSonoffSwitch()
	ctx := context.Background()

	err := sw.Run(ctx)
	if err != nil {
		return err
	}
	return nil
}
