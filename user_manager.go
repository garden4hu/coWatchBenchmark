package coWatchBenchmark

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type userManager struct {
	r     *RoomUnit
	users []*User
	start bool
	wg    sync.WaitGroup
}

func (p *RoomUnit) CreateUsers() {
	// create Users
	for i := 0; i < p.usersCap; i++ {
		u := newUser()
		p.UserManager.users = append(p.UserManager.users, u)
		p.UserManager.wg.Add(1)
		go userRun(p, u)
	}
	p.UserManager.start = true
	log.Println("[INFO] Begin to send message")
	p.UserManager.wg.Wait()
}

func userRun(r *RoomUnit, u *User) {
	defer r.UserManager.wg.Done()
	for {
		if r.UserManager.start == false {
			time.Sleep(time.Microsecond * 50)
			continue
		}
		break
	}
	// begin to connect the server
	fmt.Println("Users try to connect")
	_ = u.connect(r)
	// a clock for sending ping
}

func NewUserManager(room *RoomUnit) *userManager {
	return &userManager{r: room, start: false}
}

// return the average time consumption of all rooms which are created successfully
func (p *RoomUnit) GetUsersAvgConnectionDuration() time.Duration {
	var totalDuration time.Duration = 0
	usersSize := len(p.UserManager.users)
	for i := 0; i < usersSize; i++ {
		totalDuration += p.UserManager.users[i].ConnectionDuration
	}
	return time.Duration(int64(totalDuration) / int64(usersSize))
}
