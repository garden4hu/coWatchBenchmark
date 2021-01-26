package cowatchbenchmark

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/// as is
func NewRoom(host string, httpTimeout time.Duration, wsTimeout time.Duration, maximumUsers, msgLength, frequency int) *RoomUnit {
	if frequency == 0 {
		frequency = 10
	}
	room := &RoomUnit{usersCap: maximumUsers, msgLength: msgLength, msgSendingInternal: time.Microsecond * time.Duration(1000*1000/frequency)}
	ur, _ := url.Parse(host)
	room.Schema = ur.Scheme
	room.Host = ur.Host
	// set initial ping interval
	room.PingInterval = 25000
	room.httpTimeout = httpTimeout
	room.wsTimeout = wsTimeout
	return room
}

func (p *RoomUnit) RequestServerRoom() error {
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
	fmt.Println("request body:", string(body))
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

func (p *RoomUnit) CreateUsers(when time.Time) error {
	// create Users
	for i := 0; i < p.usersCap; i++ {
		u := newUser()
		p.wg.Add(1)
		go userRun(p, u)
	}
	now := time.Now()
	if now.UnixNano() > when.UnixNano() {
		return errors.New("current time is newer than the schedule time. Operation of creating users will not be executed")
	}
	time.Sleep(time.Nanosecond * time.Duration(when.UnixNano()-now.UnixNano()))
	p.start = true
	p.wg.Wait()
	return nil
}

func userRun(r *RoomUnit, u *User) {
	defer r.wg.Done()
	for {
		if r.start == false {
			time.Sleep(time.Microsecond * 50)
			continue
		}
		break
	}
	// begin to connect the server
	log.Println("Users try to connect")
	_ = u.connect(r)
	// a clock for sending ping
}

// ---------------------------------------- for statistics --------------------------------------------
// return the average time consumption of all rooms which are created successfully
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
