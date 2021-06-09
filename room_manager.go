package cowatchbenchmark

import (
	"errors"
	"log"
	"sync"
	"time"
)

// use "log.SetOutput(ioutil.Discard)" in main to disable log output

type RoomManager struct {
	Addr             string
	RoomSize         int
	UserSize         int
	MsgLen           int
	Frequency        int
	LckRoom          sync.Mutex
	Rooms            []*RoomUnit
	Start            bool
	HttpTimeout      time.Duration
	WSTimeout        time.Duration
	AppID            string
	SingleClientMode int
	ParallelRequest  bool
	NotifyUserAdd    <-chan int // chan 大小为 用户总数

	// for internal usage
	notifyUserAdd            chan int
	creatingRoomsOK          bool
	creatingUsersOK          bool
	finishedReqRoomRoutines  int
	finishedReqUsersRoutines int
}

// NewRoomManager will return a RoomManager
func NewRoomManager(addr string, room, user, msgLen, frequency, httpTimeout, webSocketTimeout int, appID string, singleClientMode int, parallel int) *RoomManager {
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
	rm := &RoomManager{Addr: addr, RoomSize: room, UserSize: user, MsgLen: msgLen, Frequency: frequency, Start: false, HttpTimeout: time.Second * time.Duration(httpTimeout), WSTimeout: time.Second * time.Duration(webSocketTimeout), AppID: appID, SingleClientMode: singleClientMode, ParallelRequest: parallel == 1}
	rm.creatingRoomsOK = false
	rm.creatingUsersOK = false
	rm.notifyUserAdd = make(chan int, room*user)
	rm.NotifyUserAdd = rm.notifyUserAdd
	rm.finishedReqRoomRoutines = 0
	rm.finishedReqUsersRoutines = 0
	return rm
}

func (p *RoomManager) Close() {
	close(p.notifyUserAdd)
}

func (p *RoomManager) CheckCreatingRoomsOK() bool {
	return p.creatingRoomsOK
}

func (p *RoomManager) CheckCreatingUsersOK() bool {
	return p.creatingUsersOK
}

// RequestRoom will request and create a room in server immediately.
// It will return a valid RoomUnit when the error returned is nil.
//func (p *RoomManager) RequestRoom() (*RoomUnit, error) {
//	r := NewRoom(p.Addr, p.HttpTimeout, p.WSTimeout, p.UserSize, p.MsgLen, p.Frequency, p.AppID, p.SingleClientMode)
//	err := r.Request()
//	if err != nil {
//		return nil, err
//	}
//	p.LckRoom.Lock()
//	p.Rooms = append(p.Rooms, r)
//	p.LckRoom.Unlock()
//	return r, nil
//}

// RequestAllRooms will request all the rooms from the server.
// param when is the start time for Request room from server concurrently [Only useful when parallel is true]
// param mode is the mode for Request room. true means parallel and false means serial
func (p *RoomManager) RequestAllRooms(when time.Time) error {
	var wg sync.WaitGroup
	start := make(chan struct{})

	// for serial request
	mtx := sync.Mutex{}
	leftGoroutine := p.RoomSize

	for i := 0; i < p.RoomSize; {
		// all goroutines will send request in the same time
		if p.ParallelRequest == true {
			wg.Add(1)
			go p.requestRoom(&wg, start)
			i++
		} else {
			//  线程创建，为了提高速度，一次创建 8 个
			for j := i; j < i+8 && j < p.RoomSize; j++ {
				// go p.RequestRoom()
				go func() {
					r := NewRoom(p.Addr, p.HttpTimeout, p.WSTimeout, p.UserSize, p.MsgLen, p.Frequency, p.AppID, p)
					_ = r.Request()
					mtx.Lock()
					leftGoroutine -= 1
					mtx.Unlock()
				}()
			}
			i += 8
			time.Sleep(20 * time.Millisecond)
		}
	}
	if p.ParallelRequest == true && p.SingleClientMode == 0 {
		if p.SingleClientMode == 0 { // 多台测试主机并发测试，需要等待特定时刻并发请求
			now := time.Now()
			if now.UnixNano() > when.UnixNano() {
				return errors.New("current time is newer than the schedule time. Operation of creating rooms will not be executed")
			}
			time.Sleep(time.Nanosecond * time.Duration(when.UnixNano()-now.UnixNano()))
		}
	}

	close(start) // 开始并发创建请求

	if p.ParallelRequest == true {
		wg.Wait()
	} else {
		for leftGoroutine != 0 {
			time.Sleep(1 * time.Second)
		}
	}
	p.creatingRoomsOK = true
	return nil
}

func (p *RoomManager) requestRoom(wg *sync.WaitGroup, start chan struct{}) {
	r := NewRoom(p.Addr, p.HttpTimeout, p.WSTimeout, p.UserSize, p.MsgLen, p.Frequency, p.AppID, p)
	if wg != nil {
		defer wg.Done()
	}
	if p.ParallelRequest {
		<-start // 需要等待
	}

	_ = r.Request()
}

func (p *RoomManager) RequestRoom() error {
	r := NewRoom(p.Addr, p.HttpTimeout, p.WSTimeout, p.UserSize, p.MsgLen, p.Frequency, p.AppID, p)
	err := r.Request()
	if err != nil {
		return err
	}
	return nil
}
