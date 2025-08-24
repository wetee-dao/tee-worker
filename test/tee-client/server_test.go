package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/panjf2000/gnet/v2"
)

func TestTEEServer(t *testing.T) {
	server := new(TEEServer)
	go func() {
		err := gnet.Run(server, "tcp://:18883", gnet.WithMulticore(true))
		if err != nil {
			fmt.Println("server start failed: ", err.Error())
		}
	}()

	time.Sleep(1 * time.Second)

	client, err := NewTEEClient("127.0.0.1:18883")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err = client.Start()
		if err != nil {
			fmt.Println("Client start failed")
		}
	}()

	for range 100 {
		data, err := client.Invoke("/report", []byte("report"))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("data: ", string(data))
		time.Sleep(1 * time.Second)
	}

	err = server.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
