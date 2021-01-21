package coWatchBenchmark

import (
	"fmt"
	"testing"
)

func TestRoomUnit_Chat(t *testing.T) {
	// create room witch server, https
	room := NewRoom("http://cowatch_server", 1, 20, 1)
	if err := room.RequestServerRoom(); err != nil {
		t.Error("failed to finish RequestServerRoom")
	}
	fmt.Println("room Id = ", room.Id)
	room.CreateUsers()
	fmt.Println(room.Id)
}
