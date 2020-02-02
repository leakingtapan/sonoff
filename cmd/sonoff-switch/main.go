package main

import (
	"context"

	"github.com/leakingtapan/sonoff/pkg/device"
)

func main() {
	sw := device.NewSonoffSwitch()
	ctx := context.Background()

	err := sw.Run(ctx)
	if err != nil {
		panic(err)
	}
}
