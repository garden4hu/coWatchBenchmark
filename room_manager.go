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

// RoomManager is using for manage rooms.
// Host, format for Host: Schema://server:port or Schema://ip:port
func (p *RoomManager) RequestRoomsFromServer(when time.Time) error {
	var wg sync.WaitGroup
	start := false
	for i := 0; i < p.RoomSize; i++ {
		wg.Add(1)
		go requestRoomSync(&wg, p, &start)
	}
	now := time.Now()
	if now.UnixNano() > when.UnixNano() {
		return errors.New("current time is newer than the schedule time. Operation of creating rooms will not be executed")
	}
	time.Sleep(time.Nanosecond * time.Duration(when.UnixNano()-now.UnixNano()))
	start = true
	wg.Wait()
	return nil
}

func requestRoomSync(wg *sync.WaitGroup, rm *RoomManager, start *bool) {
	defer wg.Done()
	r := NewRoom(rm.Addr, rm.HttpTimeout, rm.WSTimeout, rm.UserSize, rm.MsgLen, rm.Frequency)
	for {
		if *start == false {
			time.Sleep(time.Millisecond * 10)
			continue
		}
		break
	}
	err := r.RequestServerRoom()
	if err != nil {
		return
	}
	rm.LckRoom.Lock()
	rm.Rooms = append(rm.Rooms, r)
	rm.LckRoom.Unlock()
}

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
	return &RoomManager{Addr: addr, RoomSize: room, UserSize: user, MsgLen: msgLen, Frequency: frequency, Start: false, HttpTimeout: time.Second * time.Duration(httpTimeout), WSTimeout: time.Second * time.Duration(webSocketTimeout)}
}

// return the average time consumption of all rooms which are created successfully
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
