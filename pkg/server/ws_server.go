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
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("read: %s", err)
			break
		}
		log.Print("REQ | WS | DEV | %s", message)
		switch messageType {
		case websocket.TextMessage:
			var msg Message
			err := json.Unmarshal([]byte(message), &msg)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}
			log.Println(msg)
			err = ws.handleMessage(&msg, c)
			if err != nil {
				log.Printf("Failed to handle message: %s", err)
				continue
			}
		default:
			log.Printf("Non-supported message type: %d", messageType)
		}
	}
}

func (ws *WsServer) handleMessage(message *Message, conn *websocket.Conn) error {
	switch message.Action {
	case Register:
		err := ws.Register(message, conn)
		if err != nil {
			return err
		}
	case Update:
		ws.Update(message, conn)
	case Query:
		ws.Query(message, conn)
	case Date:
		ws.Date(message, conn)
	default:
		log.Println("Unsupported message action: %s", message.Action)
	}
	return nil
}

func (ws *WsServer) Register(msg *Message, conn *websocket.Conn) error {
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

	payload, err := json.Marshal(&resp)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, payload)
	if err != nil {
		return err
	}

	log.Printf("INFO | WS | Device %s registered", device.Id)
	return nil
}

func (ws *WsServer) Update(msg *Message, conn *websocket.Conn) {
	log.Printf("INFO | WS | DEV %s", msg)
}

func (ws *WsServer) Query(msg *Message, conn *websocket.Conn) {
	log.Printf("INFO | WS | DEV %s", msg)
}

func (ws *WsServer) Date(msg *Message, conn *websocket.Conn) {
	log.Printf("INFO | WS | DEV %s", msg)
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
