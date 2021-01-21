package cowatchbenchmark

import (
	"fmt"
	"testing"
)

func TestRoomManager_RequestRoomsFromServer(t *testing.T) {
	rm := NewRoomManager("https://cowatch_server", 3000, 1, 20, 1)
	rm.RequestRoomsFromServer("")
	fmt.Println("target room size :=", rm.RoomSize, "real room size = : ", len(rm.Rooms))
}
