package cowatchbenchmark

import (
	"github.com/google/uuid"
	"hash/crc32"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

// generate a uuid v4
func getUUID() string {
	return uuid.New().String()
}

// randStringBytes
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// generate a tens digit as hostUid
func getHostId() int {
	customHash := func(s string) int {
		v := int(crc32.ChecksumIEEE([]byte(s)))
		if v >= 0 {
			return v
		}
		if -v >= 0 {
			return -v
		}
		// v == MinInt
		return 0
	}
	return customHash(getUUID())
}

func generateUserName(length int) string {
	return randStringBytes(length)
}

// generate text message randomly for user
func generateMessage(r *RoomUnit) []byte {
	msg := "42/" + r.roomName + ",[\"CMD:chat\",\"" + randStringBytes(r.msgLength) + "\"]"
	return []byte(msg)
}
