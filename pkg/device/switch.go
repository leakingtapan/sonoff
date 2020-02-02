package device

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/leakingtapan/sonoff/pkg/types"
)

type SonoffSwitch struct {
	ws            *websocket.Conn
	serverIp      string
	websocketPort int

	deviceId   string
	apiKey     string
	version    int
	romVersion string
	model      string

	serverApiKey string
	// the state of the switch as "on" or "off"
	state string
}

func NewSonoffSwitch(serverIp string, websocketPort int) *SonoffSwitch {
	return &SonoffSwitch{
		websocketUrl: "wss://50.18.84.251:443/api/ws",
		deviceId:     "",
		apiKey:       "",
		version:      2,
		romVersion:   "1.5.5",
		model:        "ITA-GZ1-GL",
		state:        "off",
	}
}

func (s *SonoffSwitch) Run(ctx context.Context) error {
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

	// TODO: close connection
	ws, _, err := dialer.Dial(s.websocketUrl(), headers)
	if err != nil {
		return err
	}

	s.ws = ws

	err = s.Register()
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.Date()
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.Update()
	if err != nil {
		log.Println(err)
		return err
	}

	err = s.Query()
	if err != nil {
		log.Println(err)
		return err
	}

	go s.loop(ctx)

	select {
	case <-ctx.Done():
		break
	}

	return nil
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
			break
		}

		log.Printf("< %s", string(buf))
		switch msgType {
		case websocket.TextMessage:
			resp, err := s.handleMessage(buf)
			if err != nil {
				log.Printf("Failed to handle message: %s", err)
				continue
			}
			log.Printf("> %s", string(resp))
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

//TODO: impl
func (s *SonoffSwitch) Dispatch() {
}

func (s *SonoffSwitch) Register() error {
	req := types.Message{
		Action:     types.Register,
		UserAgent:  "device",
		ApiKey:     s.apiKey,
		DeviceId:   s.deviceId,
		Version:    s.version,
		RomVersion: s.romVersion,
		Model:      s.model,
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

func (s *SonoffSwitch) Date() error {
	req := types.Message{
		Action:    types.Date,
		UserAgent: "device",
		ApiKey:    s.serverApiKey,
		DeviceId:  s.deviceId,
	}
	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = s.roundTrip(data)
	return err
}

func (s *SonoffSwitch) Update() error {
	req := types.UpdateMessage{
		Message: types.Message{
			Action:    types.Update,
			UserAgent: "device",
			ApiKey:    s.serverApiKey,
			DeviceId:  s.deviceId,
		},
		Params: types.UpdateParams{
			Switch: s.state,
		},
	}

	data, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = s.roundTrip(data)
	return err
}

func (s *SonoffSwitch) Query() error {
	req := types.QueryMessage{
		Message: types.Message{
			Action:    types.Query,
			UserAgent: "device",
			ApiKey:    s.serverApiKey,
			DeviceId:  s.deviceId,
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
		s.state = updateMsg.Params.Switch

		resp := struct {
			Err       int    `json:"error"`
			UserAgent string `json:"userAgent"`
			ApiKey    string `json:"apikey"`
			DeviceId  string `json:"deviceid"`
			Sequence  string `json:"sequence"`
		}{
			DeviceId:  s.deviceId,
			UserAgent: "device",
			ApiKey:    s.serverApiKey,
			Err:       0,
			Sequence:  updateMsg.Sequence,
		}

		return json.Marshal(&resp)
	default:
		log.Printf("Unsupported action: %s", msg.Action)
		return nil, err
	}

}

func (s *SonoffSwitch) websocketUrl() string {
	//websocketUrl: "wss://50.18.84.251:443/api/ws",
	return fmt.Sprintf("wss://%s:%d/api/ws", s.serverIp, s.websocketPort)
}
