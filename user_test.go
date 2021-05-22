package cowatchbenchmark

import (
	_ "net/http/pprof"
	"testing"
)

func TestRoomUnit_Chat(t *testing.T) {
	// create r witch server, https
	//r := NewRoom("http://cowatch_server", 25*time.Second, 45*time.Second, 300, 20, 1, "appid", 1)
	//if err := r.Request(); err != nil {
	//	t.Error("failed to finish Request")
	//}
	//log.SetFlags(0)
	//log.SetOutput(ioutil.Discard)
	//fmt.Println("room roomName = ", r.roomName)
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//ch := make(chan struct{})
	//
	//go r.UsersConnection(ch, true,ctx)
	//close(ch)
	//// pprof
	//go func() {
	//	fmt.Println("pprof start...")
	//	fmt.Println(http.ListenAndServe("127.0.0.1:9876", nil))
	//}()
}
