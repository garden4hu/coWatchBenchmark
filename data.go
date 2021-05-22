package cowatchbenchmark

import (
	"sync"
	"time"
)

type RoomUnit struct {
	Address      string // server+port
	Schema       string
	roomName     string // room name
	RoomId       int    // room ID
	Password     string
	httpTimeout  time.Duration
	wsTimeout    time.Duration
	PingInterval int // ws/wss
	PingTimeout  int // ws/wss
	RtcToken     string
	muxUsers     sync.Mutex
	Users        []*User // valid users in room
	AppId        string
	ExpireTime   int
	SdkVersion   string

	condMutex *sync.Mutex // used for conditional waiting
	cond      *sync.Cond

	// for statistics
	ConnectionDuration time.Duration

	roomManager *RoomManager

	// for internal usage
	chanStop           chan bool
	wg                 sync.WaitGroup
	start              bool          // start to concurrent Request
	usersCap           int           // users cap in this room
	usersOnline        int           // online users
	msgLength          int           // length of message
	msgSendingInternal time.Duration // Microsecond as the unit
}

type User struct {
	name               string     //uuid
	sid                string     // correspond with name roomName
	uid                int        // digital id
	Lw                 sync.Mutex // lock for writing
	connected          bool
	readyForMsg        bool
	ConnectionDuration time.Duration
	hostCoWatch        bool // only the user who create the room can be the host
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
