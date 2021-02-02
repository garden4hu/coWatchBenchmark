package cowatchbenchmark

import (
	"errors"
	"log"
	"sync"
	"time"
)

// use "log.SetOutput(ioutil.Discard)" in main to disable log output

type RoomManager struct {
	Addr        string
	RoomSize    int
	UserSize    int
	MsgLen      int
	Frequency   int
	LckRoom     sync.Mutex
	Rooms       []*RoomUnit
	Start       bool
	HttpTimeout time.Duration
	WSTimeout   time.Duration
}

// NewRoomManager will return a RoomManager
func NewRoomManager(addr string, room, user, msgLen, frequency, httpTimeout, webSocketTimeout int) *RoomManager {
	if room < 0 || user < 0 || frequency <= 0 {
		log.Fatalln("Invalid param")
		return nil
	}
	if httpTimeout > 60 || httpTimeout < 0 {
		httpTimeout = 60
	}
	if webSocketTimeout > 60 || webSocketTimeout < 0 {
		webSocketTimeout = 45
	}
	if frequency <= 0 {
		frequency = 1
	}
	return &RoomManager{Addr: addr, RoomSize: room, UserSize: user, MsgLen: msgLen, Frequency: frequency, Start: false, HttpTimeout: time.Second * time.Duration(httpTimeout), WSTimeout: time.Second * time.Duration(webSocketTimeout)}
}

// RequestRoom will request and create a room in server immediately.
// It will return a valid RoomUnit when the error returned is nil.
func (p *RoomManager) RequestRoom() (*RoomUnit, error) {
	r := NewRoom(p.Addr, p.HttpTimeout, p.WSTimeout, p.UserSize, p.MsgLen, p.Frequency)
	err := r.Request()
	if err != nil {
		return nil, err
	}
	p.LckRoom.Lock()
	p.Rooms = append(p.Rooms, r)
	p.LckRoom.Unlock()
	return r, nil
}

// RequestAllRooms will request all the rooms from the server.
// param when is the start time for Request room from server concurrently.
// param mode is the mode for Request room. 0 means parallel and 1 means serial
func (p *RoomManager) RequestAllRooms(when time.Time, mode int) error {
	var wg sync.WaitGroup
	start := false
	for i := 0; i < p.RoomSize; i++ {
		wg.Add(1)
		if mode == 0 {
			go requestRoom(&wg, p, &start)
		} else {
			start = true
			requestRoom(&wg, p, &start)
		}
	}
	if mode == 0 {
		now := time.Now()
		if now.UnixNano() > when.UnixNano() {
			return errors.New("current time is newer than the schedule time. Operation of creating rooms will not be executed")
		}
		time.Sleep(time.Nanosecond * time.Duration(when.UnixNano()-now.UnixNano()))
		start = true
	}
	wg.Wait()
	return nil
}

// GetCreatingRoomAvgDuration return the average time consumption of all rooms which are created successfully
func (p *RoomManager) GetCreatingRoomAvgDuration() time.Duration {
	if len(p.Rooms) == 0 {
		return time.Duration(0)
	}
	var totalDuration time.Duration = 0
	for i := 0; i < len(p.Rooms); i++ {
		totalDuration += p.Rooms[i].ConnectionDuration
	}
	return time.Duration(int64(totalDuration) / int64(len(p.Rooms)))
}

func (p *RoomManager) GetCreatingUsersAvgDuration() time.Duration {
	if len(p.Rooms) == 0 {
		return time.Duration(0)
	}
	var totalDuration time.Duration = 0
	for i := 0; i < len(p.Rooms); i++ {
		totalDuration += p.Rooms[i].GetUsersAvgConnectionDuration()
	}
	return time.Duration(int64(totalDuration) / int64(len(p.Rooms)))
}

func requestRoom(wg *sync.WaitGroup, rm *RoomManager, start *bool) {
	defer wg.Done()
	r := NewRoom(rm.Addr, rm.HttpTimeout, rm.WSTimeout, rm.UserSize, rm.MsgLen, rm.Frequency)
	for {
		if *start == false {
			time.Sleep(time.Millisecond * 20)
			continue
		}
		break
	}
	err := r.Request()
	if err != nil {
		return
	}
	rm.LckRoom.Lock()
	rm.Rooms = append(rm.Rooms, r)
	rm.LckRoom.Unlock()
}
