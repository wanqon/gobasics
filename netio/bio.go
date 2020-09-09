package netio

import (
	"fmt"
	"net"
	"os"
)

func bioServer()  {
	//socket() bind() listen()
	listener, err := net.Listen("tcp",":9001")
	checkError(err)
	for {
		//阻塞，strace log:accept(3,
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		conn.Close()
	}

}

func checkError(err error)  {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
