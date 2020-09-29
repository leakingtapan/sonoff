package device

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/leakingtapan/sonoff/pkg/types"
)

type SonoffSwitch struct {
	types.Device
	mu sync.Mutex
	// The websocket connection used by the device
	// Each device has a dedicated connection
	ws *websocket.Conn
	// the IP address of the dispatch server
	serverIp string
	// the IP address of the websocket server
	websocketServerIp string
	// the port address of the websocket server
	websocketPort int
	// the API key used to authenticate into the backend server
	serverApiKey string
	// the channel for switch state change events
	watchCh chan string
}

func NewSonoffSwitch(
	serverIp string,
	websocketServerIp string,
	websocketPort int,
	device types.Device,
) *SonoffSwitch {
	return &SonoffSwitch{
		Device:            device,
		serverIp:          serverIp,
		websocketServerIp: websocketServerIp,
		websocketPort:     websocketPort,
		watchCh:           make(chan string, 128),
	}
}

func (s *SonoffSwitch) Run(ctx context.Context) error {
	if s.websocketServerIp == "" || s.websocketPort == 0 {
		err := s.dispatch()
		if err != nil {
			return err
		}
	}

	ws, err := s.connect()
	if err != nil {
		return err
	}

	s.ws = ws

	initFuncs := []func() error{
		s.register,
		s.date,
		s.update,
		s.query,
	}

	for _, initFunc := range initFuncs {
		err = initFunc()
		if err != nil {
			log.Println(err)
			return err
		}
	}

	go s.loop(ctx)

	select {
	case <-ctx.Done():
		break
	}

	return nil
}

func (s *SonoffSwitch) GetState() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.Device.State
}

func (s *SonoffSwitch) SetState(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Device.State = state
	//TODO: handle channel full
	s.watchCh <- state
}

func (s *SonoffSwitch) Watch() <-chan string {
	return s.watchCh
}

func (s *SonoffSwitch) connect() (*websocket.Conn, error) {
	////TODO: check ws origin
	//dest, err := url.Parse(s.websocketUrl)
	//if err != nil {
	//	return err
	//}
	//originURL := *dest
	//if dest.Scheme == "wss" {
	//	originURL.Scheme = "https"
	//} else {
	//	originURL.Scheme = "http"
	//}
	//origin := originURL.String()
	headers := make(http.Header)
	//headers.Add("Origin", origin)

	dialer := websocket.Dialer{
		Proxy: http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// TODO: handle close connection
	ws, _, err := dialer.Dial(s.websocketUrl(), headers)
	if err != nil {
		return nil, err
	}
	return ws, err
}

func (s *SonoffSwitch) loop(ctx context.Context) {
	for {
		// TODO: cancellation
		//select {
		//case <-ctx.Done():
		//	break
		//default:

		//}
		msgType, buf, err := s.ws.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %s", err)
			//TODO: handle connection closed by server
			break
		}

		log.Printf("RX %s", string(buf))
		switch msgType {
		case websocket.TextMessage:
			resp, err := s.handleMessage(buf)
			if err != nil {
				log.Printf("Failed to handle message: %s", err)
				continue
			}

			if resp == nil {
				continue
			}

			log.Printf("TX %s", string(resp))
			err = s.ws.WriteMessage(websocket.TextMessage, resp)
			if err != nil {
				log.Printf("Failed to write response message: %s", err)
				continue
			}
		default:
			log.Println("Unknown message: %s", msgType)
		}
	}
}

func (s *SonoffSwitch) dispatch() error {
	dispatchUrl := fmt.Sprintf("http://%s/dispatch/device", s.serverIp)
	device := struct {
		types.Device
		ApiKey     string `json:"apikey"`
		Accept     string `json:"accept"`
		RomVersion string `json:"romVersion"`
		Ts         int    `json:"ts"`
	}{
		Device: types.Device{
			DeviceId: s.DeviceId,
			Version:  s.Version,
			Model:    s.Model,
		},
		Accept:     "ws;2",
		ApiKey:     s.ApiKey,
		RomVersion: s.RomVersion,
		Ts:         119,
	}
	payload, err := json.Marshal(&device)
	if err != nil {
		return err
	}

	log.Printf("POST %s", payload)
	resp, err := http.Post(dispatchUrl, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	log.Printf("POST DONE")

	var respData struct {
		Port   int    `json:"port"`
		Reason string `json:"reason"`
		Ip     string `json:"IP"`
		Error  int    `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return err
	}

	log.Printf("RES %+v", respData)
	s.websocketServerIp = respData.Ip
	s.websocketPort = respData.Port

	return nil
}

func (s *SonoffSwitch) register() error {
	req := types.Message{
		Action:     types.Register,
		UserAgent:  "device",
		ApiKey:     s.ApiKey,
		DeviceId:   s.DeviceId,
		Version:    s.Version,
		RomVersion: s.RomVersion,
		Model:      s.Model,
		//Ts:         time.Now().Unix(), //TODO: precision
		Ts: 193,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	buf, err := s.roundTrip(data)
	if err != nil {
		return err
	}

	var resp types.Message
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		return err
	}

	s.serverApiKey = resp.ApiKey
	return nil
}

func (s *SonoffSwitch) date() error {
	req := types.Message{
		Action:    types.Date,
		UserAgent: "device",
		ApiKey:    s.serverApiKey,
		DeviceId:  s.DeviceId,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = s.roundTrip(data)
	return err
}

func (s *SonoffSwitch) update() error {
	req := types.UpdateMessage{
		Message: types.Message{
			Action:    types.Update,
			UserAgent: "device",
			ApiKey:    s.serverApiKey,
			DeviceId:  s.DeviceId,
		},
		Params: types.UpdateParams{
			Switch: s.State,
		},
	}

	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = s.roundTrip(data)
	return err
}

func (s *SonoffSwitch) query() error {
	req := types.QueryMessage{
		Message: types.Message{
			Action:    types.Query,
			UserAgent: "device",
			ApiKey:    s.serverApiKey,
			DeviceId:  s.DeviceId,
		},
		Params: []string{"timers"},
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = s.roundTrip(data)
	return err
}

func (s *SonoffSwitch) roundTrip(request []byte) ([]byte, error) {
	log.Printf("> %s", string(request))
	err := s.ws.WriteMessage(websocket.TextMessage, request)
	if err != nil {
		return nil, err
	}

	msgType, buf, err := s.ws.ReadMessage()
	if err != nil {
		return nil, err
	}
	if msgType != websocket.TextMessage {
		return nil, fmt.Errorf("Unexpected message type: %s", msgType)
	}

	log.Printf("< %s", string(buf))
	return buf, nil
}

func (s *SonoffSwitch) handleMessage(request []byte) ([]byte, error) {
	var msg types.Message
	err := json.Unmarshal(request, &msg)
	if err != nil {
		log.Printf("Failed to unmarshal request: %s", err)
		return nil, err
	}

	switch msg.Action {
	case types.Update:
		// TODO: use UpdateMessageFromApp type
		var updateMsg types.UpdateMessage
		err := json.Unmarshal(request, &updateMsg)
		if err != nil {
			log.Println("Failed to unmarshal update message: %s", err)
			return nil, err
		}
		s.mu.Lock()
		s.State = updateMsg.Params.Switch
		s.mu.Unlock()

		resp := struct {
			Err       int    `json:"error"`
			UserAgent string `json:"userAgent"`
			ApiKey    string `json:"apikey"`
			DeviceId  string `json:"deviceid"`
			Sequence  string `json:"sequence"`
		}{
			DeviceId:  s.DeviceId,
			UserAgent: "device",
			ApiKey:    s.serverApiKey,
			Err:       0,
			Sequence:  updateMsg.Sequence,
		}

		return json.Marshal(&resp)
	default:
		log.Printf("Unsupported action: %s", msg.Action)
	}

	return nil, nil
}

//websocketUrl return the websocket URL in the form of
// "wss://50.18.84.251:443/api/ws"
func (s *SonoffSwitch) websocketUrl() string {
	return fmt.Sprintf("wss://%s:%d/api/ws", s.websocketServerIp, s.websocketPort)
}
