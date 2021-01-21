package coWatchBenchmark

import (
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func newUser() *User {
	return &User{client: getUUID(), isConnected: false, readyForMsg: false}
}

func (p *User) connect(room *RoomUnit) error {
	v := url.Values{}
	v.Add("clientId", p.client)
	v.Add("Password", room.Password)
	v.Add("EIO", "3")
	v.Add("transport", "websocket")
	u := url.URL{Host: room.Host, Path: "/socket.io/", ForceQuery: true, RawQuery: v.Encode()}
	switch room.Schema {
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
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("failed to dial websocket:", err)
		return err
	}
	p.ConnectionDuration = time.Since(start)
	defer conn.Close()
	p.isConnected = true
	done := make(chan struct{})
	go process(conn, room, p, done)
	defer close(done)
	pingTicker := time.NewTicker(time.Millisecond * time.Duration(room.PingInterval))
	log.Println("ping ticker duration:", room.PingInterval)
	defer pingTicker.Stop()
	sendMsgTicker := time.NewTicker(room.msgSendingInternal)
	defer sendMsgTicker.Stop()
	log.Println("sending MSG  ticker duration:", room.msgSendingInternal.String())

	for {
		select {
		case <-done:
			return nil
		case _ = <-pingTicker.C:
			// reset pingTicker and send ping
			p.Lw.Lock()
			err := conn.WriteMessage(websocket.TextMessage, []byte("2"))
			p.Lw.Unlock()
			if err != nil {
				log.Println("write:", err)
				return err
			}
			log.Println("sending Ping MSG")
			pingTicker.Reset(time.Millisecond * time.Duration(room.PingInterval))
			break
			// sending msg
		case _ = <-sendMsgTicker.C:
			if p.readyForMsg {
				msg := generateMessage(room)
				p.Lw.Lock()
				_ = conn.WriteMessage(websocket.TextMessage, msg)
				log.Println("sending msg Frequency:", msg)
				p.Lw.Unlock()
			}
			sendMsgTicker.Reset(room.msgSendingInternal)
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
