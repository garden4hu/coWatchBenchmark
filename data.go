package cowatchbenchmark

import (
	"sync"
	"time"
)

type RoomUnit struct {
	Host         string // server+port
	Schema       string
	Id           string // room ID
	Password     string
	PingInterval int // ws/wss
	PingTimeout  int // ws/wss
	RtcToken     string
	Users        []User //users in room
	UserManager  *userManager
	// for statistics
	ConnectionDuration time.Duration

	// for internal usage
	chanStop           chan bool
	start              bool          // start to concurrent request
	usersCap           int           // users in this room
	usersOnline        int           // online users
	msgLength          int           // length of message
	msgSendingInternal time.Duration // Microsecond as the unit
}

type User struct {
	client             string     //uuid
	sid                string     // correspond with client Id
	Lw                 sync.Mutex // lock for writing
	isConnected        bool
	readyForMsg        bool
	ConnectionDuration time.Duration
}

type RequestedUserInfo struct {
	Sid          string `json:"sid"`
	Upgrades     []int  `json:"upgrades"`
	PingInterval int    `json:"PingInterval"`
	PingTimeOut  int    `json:"PingTimeout"`
}

type Room struct {
	Name string `json:"name"`
}
