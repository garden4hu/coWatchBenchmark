package cowatchbenchmark

import (
	"fmt"
	"testing"
	"time"
)

func TestRoomUnit_GetRoomID(t *testing.T) {
	room := NewRoom("http://cowatch_server", 25*time.Second, 45*time.Second, 1, 20, 1)
	if err := room.Request(); err != nil {
		t.Error("failed to finish Request")
	}
	fmt.Println("room Id = ", room.Id)
}
