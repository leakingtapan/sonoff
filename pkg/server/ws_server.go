package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

const (
	dateLayout = "2006-01-02T15:04:05.000Z"
)

type Action string

const (
	Register Action = "register"
	Date     Action = "date"
	Query    Action = "query"
	Update   Action = "update"
)

type Message struct {
	Action     Action `json:"action"`
	DeviceId   string `json:"deviceId"`
	UserAgent  string `json:"userAgent"`
	ApiKey     string `json:"apiKey"`
	Version    int    `json:"version"`
	RomVersion string `json:"romVersion"`
	Model      string `json:"model"`
	Ts         int    `json:"ts"`
}

type UpdateMessage struct {
	Message
	Params struct {
		Switch string `json:"switch"`
	} `json:"params"`
}

type QueryMessage struct {
	Message
	Params []string `json:"params"`
}

var upgrader = websocket.Upgrader{} // use default options

type WsServer struct {
	port    int
	devices *Devices
}

func NewWsServer(port int, devices *Devices) *WsServer {
	return &WsServer{
		port:    port,
		devices: devices,
	}
}

func (ws *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %s", err)
		return
	}
	defer c.Close()

	for {
		messageType, payload, err := c.ReadMessage()
		if err != nil {
			log.Printf("read: %s", err)
			break
		}
		log.Printf("REQ | WS | DEV | %s", string(payload))
		switch messageType {
		case websocket.TextMessage:
			var msg Message
			err := json.Unmarshal(payload, &msg)
			if err != nil {
				log.Printf("Failed to unmarshal payload: %s", err)
				continue
			}

			err = ws.handleMessage(&msg, payload, c)
			if err != nil {
				log.Printf("Failed to handle message: %s", err)
				continue
			}
		default:
			log.Printf("Non-supported message type: %d", messageType)
		}
	}
}

func (ws *WsServer) handleMessage(message *Message, payload []byte, conn *websocket.Conn) error {
	var resp []byte
	var err error
	switch message.Action {
	case Register:
		resp, err = ws.Register(message, conn)
	case Update:
		resp, err = ws.Update(payload, conn)
	case Query:
		resp, err = ws.Query(payload)
	case Date:
		resp, err = ws.Date(message)
	default:
		log.Printf("Unsupported message action: %s", message.Action)
		resp, err = ws.Ack(message)
	}
	if err != nil {
		return err
	}

	log.Printf("RES | WS | %s", string(resp))
	return conn.WriteMessage(websocket.TextMessage, resp)
}

func (ws *WsServer) Register(msg *Message, conn *websocket.Conn) ([]byte, error) {
	device := Device{
		Id:      msg.DeviceId,
		Version: msg.Version,
		Model:   msg.Model,
		Conn:    conn,
	}

	ws.devices.AddOrUpdateDevice(&device)

	resp := struct {
		Err      int    `json:"error"`
		DeviceId string `json:"deviceid"`
		ApiKey   string `json:"apikey"`
	}{
		Err:      0,
		DeviceId: device.Id,
		ApiKey:   "111111111-1111-1111-1111-111111111111",
	}

	log.Printf("INFO | WS | Device %s registered", device.Id)
	return json.Marshal(&resp)
}

func (ws *WsServer) Update(payload []byte, conn *websocket.Conn) ([]byte, error) {
	var msg UpdateMessage
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return nil, err
	}

	d, found := ws.devices.Get(msg.DeviceId)
	if !found {
		log.Printf("ERR | WS | Unknown device %s", msg.DeviceId)
	} else {
		d.State = msg.Params.Switch
		d.Conn = conn
		ws.devices.AddOrUpdateDevice(d)
	}

	resp := struct {
		Err      int    `json:"error"`
		DeviceId string `json:"deviceid"`
		ApiKey   string `json:"apikey"`
	}{
		Err:      0,
		DeviceId: msg.DeviceId,
		ApiKey:   "111111111-1111-1111-1111-111111111111",
	}

	return json.Marshal(&resp)
}

func (ws *WsServer) Query(payload []byte) ([]byte, error) {
	var msg QueryMessage
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		return nil, err
	}

	resp := struct {
		Err      int    `json:"error"`
		DeviceId string `json:"deviceid"`
		ApiKey   string `json:"apikey"`
	}{
		Err:      0,
		DeviceId: msg.DeviceId,
		ApiKey:   "111111111-1111-1111-1111-111111111111",
	}

	return json.Marshal(&resp)
}

func (ws *WsServer) Date(msg *Message) ([]byte, error) {
	resp := struct {
		Err      int    `json:"error"`
		DeviceId string `json:"deviceid"`
		ApiKey   string `json:"apikey"`
		Date     string `json:"date"`
	}{
		Err:      0,
		DeviceId: msg.DeviceId,
		ApiKey:   "111111111-1111-1111-1111-111111111111",
		Date:     time.Now().Format(dateLayout),
	}

	return json.Marshal(&resp)
}

func (ws *WsServer) Ack(msg *Message) ([]byte, error) {
	resp := struct {
		Err      int    `json:"error"`
		DeviceId string `json:"deviceid"`
		ApiKey   string `json:"apikey"`
	}{
		Err:      0,
		DeviceId: msg.DeviceId,
		ApiKey:   "111111111-1111-1111-1111-111111111111",
	}

	return json.Marshal(&resp)
}

func (ws *WsServer) Serve() {
	addr := fmt.Sprintf(":%d", ws.port)
	svr := http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      ws,
	}
	log.Fatal(svr.ListenAndServeTLS("./certs/server.crt", "./certs/server.key"))
}

func Echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read: %s", err)
			break
		}
		log.Printf("recv: %s mt: %d", message, mt)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Printf("write: %s", err)
			break
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "wss://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
