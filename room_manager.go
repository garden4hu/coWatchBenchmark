package coWatchBenchmark

import (
	"log"
	"sync"
	"time"
)

// use "log.SetOutput(ioutil.Discard)" in main to disable log output

type RoomManager struct {
	Addr      string
	RoomSize  int
	UserSize  int
	MsgLen    int
	Frequency int
	LckRoom   sync.Mutex
	Rooms     []*RoomUnit
	Start     bool
}

// RoomManager is using for manage Rooms.
// Host, format for Host: Schema://server:port or Schema://ip:port
func (p *RoomManager) RequestRoomsFromServer(when string) {
	var wg sync.WaitGroup
	start := false
	for i := 0; i < p.RoomSize; i++ {
		wg.Add(1)
		go requestRoomSync(&wg, p, p.Addr, p.UserSize, p.MsgLen, p.Frequency, &start)
	}
	now := time.Now()
	schedule, err := time.Parse(time.RFC3339, when)
	if err == nil && now.UnixNano() < schedule.UnixNano() {
		time.Sleep(time.Nanosecond * time.Duration(schedule.UnixNano()-now.UnixNano()))
	} else {
		log.Println("[WARN] failed to parsed start_time, this machine starting without sync :", time.Now().Format(time.RFC3339))
	}
	start = true
	wg.Wait()
}

func requestRoomSync(wg *sync.WaitGroup, roomManger *RoomManager, addr string, user, msgLength, frequency int, start *bool) {
	defer wg.Done()
	r := NewRoom(addr, user, msgLength, frequency)
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
	roomManger.LckRoom.Lock()
	roomManger.Rooms = append(roomManger.Rooms, r)
	roomManger.LckRoom.Unlock()
}

func NewRoomManager(addr string, room, user, msgLen, frequency int) *RoomManager {
	return &RoomManager{Addr: addr, RoomSize: room, UserSize: user, MsgLen: msgLen, Frequency: frequency, Start: false}
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
