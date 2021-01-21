package coWatchBenchmark

import (
	"fmt"
	"testing"
)

func TestRoomUnit_GetRoomID(t *testing.T) {
	room := NewRoom("http://cowatch_server", 1, 20, 1)
	if err := room.RequestServerRoom(); err != nil {
		t.Error("failed to finish RequestServerRoom")
	}
	fmt.Println("room Id = ", room.Id)
}
