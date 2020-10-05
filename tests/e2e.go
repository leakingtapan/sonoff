package e2e

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/leakingtapan/sonoff/apis/client"
	"github.com/leakingtapan/sonoff/apis/client/operations"
	"github.com/leakingtapan/sonoff/pkg/device"
	"github.com/leakingtapan/sonoff/pkg/types"
)

var (
	deviceSpec = types.Device{
		DeviceId: "10001aa123",
		ApiKey:   "test-key",
		Version:  2,
		Model:    "ITA-GZ1-GL",
		State:    "on",
	}
)
var _ = Describe("Sonoff Switch", func() {
	sw := setupSwitch()

	It("should set switch status to off", func() {
		describeDevices()
		setDeviceState(deviceSpec.DeviceId, "off")

		deviceClient := client.Default
		getDeviceStateParam := &operations.GetDeviceByIDParams{
			DeviceID:   deviceSpec.DeviceId,
			Context:    context.Background(),
			HTTPClient: insecureHttpClient(),
		}

		resp, err := deviceClient.Operations.GetDeviceByID(getDeviceStateParam)

		Expect(err).To(BeNil())
		Expect(resp.Payload.State).To(Equal("off"))
		Expect(sw.State).To(Equal("off"))
	})

	It("should set switch status to on", func() {
		describeDevices()
		setDeviceState(deviceSpec.DeviceId, "on")

		deviceClient := client.Default
		getDeviceStateParam := &operations.GetDeviceByIDParams{
			DeviceID:   deviceSpec.DeviceId,
			Context:    context.Background(),
			HTTPClient: insecureHttpClient(),
		}

		resp, err := deviceClient.Operations.GetDeviceByID(getDeviceStateParam)

		Expect(err).To(BeNil())
		Expect(resp.Payload.State).To(Equal("on"))
		Expect(sw.State).To(Equal("on"))
	})

})

func setDeviceState(deviceID string, state string) {
	deviceClient := client.Default
	setDeviceStateParam := &operations.SetDeviceStateByIDParams{
		DeviceID:   deviceID,
		State:      state,
		Context:    context.Background(),
		HTTPClient: insecureHttpClient(),
	}

	_, err := deviceClient.Operations.SetDeviceStateByID(setDeviceStateParam)
	Expect(err).To(BeNil())
}

func describeDevices() {
	deviceClient := client.Default
	params := &operations.DescribeDevicesParams{
		Context:    context.Background(),
		HTTPClient: insecureHttpClient(),
	}

	devices, err := deviceClient.Operations.DescribeDevices(params)
	Expect(err).To(BeNil())

	for id, d := range devices.Payload {
		fmt.Printf("device[%d] = %+v\n", id, d)
	}
}

func setupSwitch() *device.SonoffSwitch {
	ctx := context.Background()
	serverIp := "127.0.0.1"
	serverPort := 8080
	sw := device.NewSonoffSwitch(
		serverIp,
		serverPort,
		"",
		0,
		deviceSpec,
	)

	go func() {
		err := sw.Run(ctx)
		if err != nil {
			fmt.Printf("err = %+v\n", err)
		}
	}()

	// wait 1 sec for switch to initialize
	time.Sleep(time.Second)

	return sw
}

func insecureHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: tr}
}
