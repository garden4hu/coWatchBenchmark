package cowatchbenchmark

import (
	"fmt"
	"testing"
	"time"
)

func TestRoomManager_RequestRoomsFromServer(t *testing.T) {
	rm := NewRoomManager("https://cowatch_server", 3000, 1, 20, 1, 25000, 25000)
	err := rm.RequestRoomsFromServer(time.Now())
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Println("target room size :=", rm.RoomSize, "real room size = : ", len(rm.Rooms))
}
