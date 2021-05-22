package cowatchbenchmark

import "time"

// ---------------------------------------- for statistics --------------------------------------------

func (p *RoomManager) GetCreatedRooms() int {
	return len(p.Rooms)
}

func (p *RoomManager) GetTotalUsers() int {
	total := 0
	for i := 0; i < len(p.Rooms); i++ {
		total += len(p.Rooms[i].Users)
	}
	return total
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
		totalDuration += p.Rooms[i].usersAvgConnectionDuration()
	}
	return time.Duration(int64(totalDuration) / int64(len(p.Rooms)))
}

// usersAvgConnectionDuration return the average time consumption of all rooms which are created successfully
func (p *RoomUnit) usersAvgConnectionDuration() time.Duration {
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
