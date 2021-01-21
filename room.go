package cowatchbenchmark

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (p *RoomUnit) RequestServerRoom() error {
	strings.TrimSuffix(p.Host, "/")
	uri := p.Schema + "://" + p.Host + "/" + "createRoom"
	tr := func() *http.Transport {
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}()
	client := &http.Client{Transport: tr, Timeout: time.Second * 50}
	start := time.Now()
	resp, err := client.Post(uri, "application/json", nil)
	if err != nil {
		fmt.Println("Failed to post, err = ", err)
		return fmt.Errorf("failed to post :%s err:%s", uri, err.Error())
	}
	defer resp.Body.Close()
	p.ConnectionDuration = time.Since(start)
	roomRaw, _ := ioutil.ReadAll(resp.Body)
	// unmarshal
	room := new(Room)
	err = json.Unmarshal(roomRaw, room)
	if err != nil {
		return errors.New("failed to parse json of room")
	}
	p.Id = room.Name
	return nil
}

/// as is
func NewRoom(host string, maximumUsers, msgLength, frequency int) *RoomUnit {
	if frequency == 0 {
		frequency = 10
	}
	room := &RoomUnit{usersCap: maximumUsers, msgLength: msgLength, msgSendingInternal: time.Microsecond * time.Duration(1000*1000/frequency)}
	ur, _ := url.Parse(host)
	room.Schema = ur.Scheme
	room.Host = ur.Host
	// set initial ping interval
	room.PingInterval = 25000
	room.UserManager = NewUserManager(room)
	return room
}
