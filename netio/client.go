package netio

import (
	"fmt"
	"net"
)

func client()  {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("Dial err:", err.Error())
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Are you ready?"))
	if err != nil {
		fmt.Println("Write err:", err.Error())
		return
	}
}
