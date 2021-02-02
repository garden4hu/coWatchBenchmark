package cowatchbenchmark

import (
	"bytes"
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

// NewRoom return a RoomUnit object
func NewRoom(host string, httpTimeout time.Duration, wsTimeout time.Duration, maximumUsers, msgLength, frequency int) *RoomUnit {
	if frequency == 0 {
		frequency = 10
	}
	room := &RoomUnit{usersCap: maximumUsers, msgLength: msgLength, msgSendingInternal: time.Microsecond * time.Duration(60*1000*1000/frequency)}
	ur, _ := url.Parse(host)
	room.Schema = ur.Scheme
	room.Host = ur.Host
	// set initial ping interval
	room.PingInterval = 25000
	room.httpTimeout = httpTimeout
	room.wsTimeout = wsTimeout
	return room
}

func (p *RoomUnit) Request() error {
	strings.TrimSuffix(p.Host, "/")
	uri := p.Schema + "://" + p.Host + "/" + "createRoom"
	tr := func() *http.Transport {
		return &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}()
	client := &http.Client{Transport: tr, Timeout: p.httpTimeout}
	start := time.Now()
	// construct body
	body, _ := json.Marshal(CreateRoomReqBody{HostUid: getUUID()})
	bodyReader := bytes.NewReader(body)
	resp, err := client.Post(uri, "application/json", bodyReader)
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

// UsersConnection try to connect to the server and exchange message.
// param when is the time for requesting of websocket concurrently
// param mode is the mode for requesting. 0 means parallel and 1 means serial
func (p *RoomUnit) UsersConnection(when time.Time, mode int) error {
	// create Users
	finish := make(chan bool)
	defer close(finish)
	for i := 0; i < p.usersCap; i++ {
		u := newUser()
		if mode == 1 {
			p.start = true
		}
		go u.joinRoom(p, mode, finish)
		if mode == 1 {
			_ = <-finish
		}
	}
	if mode == 0 {
		now := time.Now()
		if now.UnixNano() > when.UnixNano() {
			return errors.New("current time is newer than the schedule time. Operation of creating users will not be executed")
		}
		time.Sleep(time.Nanosecond * time.Duration(when.UnixNano()-now.UnixNano()))
		p.start = true
	}
	return nil
}

// ---------------------------------------- for statistics --------------------------------------------

// GetUsersAvgConnectionDuration return the average time consumption of all rooms which are created successfully
func (p *RoomUnit) GetUsersAvgConnectionDuration() time.Duration {
	var totalDuration time.Duration = 0
	usersSize := len(p.Users)
	if usersSize == 0 {
		return totalDuration
	}
	for i := 0; i < usersSize; i++ {
		totalDuration += p.Users[i].ConnectionDuration
	}
	return time.Duration(int64(totalDuration) / int64(usersSize))
}

type CreateRoomReqBody struct {
	HostUid string `json:"hostUid"`
}
