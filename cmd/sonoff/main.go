package main

import (
	"fmt"
	"os"

	server "github.com/leakingtapan/sonoff/cmd/sonoff-server"
	switchdev "github.com/leakingtapan/sonoff/cmd/sonoff-switch"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:        "sonoff",
	Short:      "sonoff cli for both server and client",
	SuggestFor: []string{"sonoff"},
}

func init() {
	cobra.EnablePrefixMatching = true

	rootCmd.AddCommand(
		server.NewServerCommand(),
		switchdev.NewSwitchCommand(),
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "sonoff failed %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
