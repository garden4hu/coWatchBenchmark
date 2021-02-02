package cowatchbenchmark

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func newUser() *User {
	return &User{client: getUUID(), isConnected: false, readyForMsg: false}
}

func (p *User) joinRoom(r *RoomUnit, mode int, ok chan bool) {
	// if Request user concurrently
	if mode == 0 {
		for {
			if r.start == false {
				time.Sleep(time.Microsecond * 50)
				continue
			}
			break
		}
	}

	v := url.Values{}
	v.Add("clientId", p.client)
	v.Add("Password", r.Password)
	v.Add("EIO", "3")
	v.Add("transport", "websocket")
	u := url.URL{Host: r.Host, Path: "/socket.io/", ForceQuery: true, RawQuery: v.Encode()}
	switch r.Schema {
	case "http":
		u.Scheme = "ws"
		break
	case "https":
		u.Scheme = "wss"
		break
	default:
		u.Scheme = "wss"
		break
	}
	start := time.Now()
	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: r.wsTimeout,
	}
	conn, _, err := dialer.Dial(u.String(), nil)
	// this
	if mode != 0 {
		ok <- true // notify that this goroutine has finish the websocket Request
	}
	if err != nil {
		log.Println("failed to dial websocket:", err)
		return
	}
	// add user to RoomUnit
	r.muxUsers.Lock()
	r.Users = append(r.Users, p)
	r.muxUsers.Unlock()
	p.ConnectionDuration = time.Since(start)
	defer conn.Close()
	p.isConnected = true
	done := make(chan struct{})
	go process(conn, r, p, done)
	defer close(done)
	pingTicker := time.NewTicker(time.Millisecond * time.Duration(r.PingInterval))
	log.Println("ping ticker duration:", r.PingInterval)
	defer pingTicker.Stop()
	sendMsgTicker := time.NewTicker(r.msgSendingInternal)
	defer sendMsgTicker.Stop()
	log.Println("sending MSG  ticker duration:", r.msgSendingInternal.String())

	for {
		select {
		case <-done:
			return
		case _ = <-pingTicker.C:
			// reset pingTicker and send ping
			p.Lw.Lock()
			err := conn.WriteMessage(websocket.TextMessage, []byte("2"))
			p.Lw.Unlock()
			if err != nil {
				log.Println("write:", err)
				return
			}
			log.Println("sending Ping MSG")
			pingTicker.Reset(time.Millisecond * time.Duration(r.PingInterval))
			break
			// sending msg
		case _ = <-sendMsgTicker.C:
			if p.readyForMsg {
				msg := generateMessage(r)
				p.Lw.Lock()
				_ = conn.WriteMessage(websocket.TextMessage, msg)
				log.Println("sending msg Frequency:", msg)
				p.Lw.Unlock()
			}
			sendMsgTicker.Reset(r.msgSendingInternal)
			break
		}
	}
}

// generate msg randomly
func generateMessage(r *RoomUnit) []byte {
	msg := "42/" + r.Id + ",[\"CMD:chat\",\"" + randStringBytes(r.msgLength) + "\"]"
	return []byte(msg)
}

func getUUID() string {
	return uuid.New().String()
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
