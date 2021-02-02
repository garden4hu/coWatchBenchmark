package cowatchbenchmark

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestRoomUnit_Chat(t *testing.T) {
	// create r witch server, https
	r := NewRoom("http://cowatch_server", 25*time.Second, 45*time.Second, 300, 20, 1)
	if err := r.Request(); err != nil {
		t.Error("failed to finish Request")
	}
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	fmt.Println("room Id = ", r.Id)
	go r.UsersConnection(time.Now().Add(time.Second), 0)
	// pprof
	go func() {
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe("127.0.0.1:9876", nil))
	}()

	select {}
}
