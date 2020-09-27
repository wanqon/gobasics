package netio

import (
	"fmt"
	"net"
	"os"
	"time"
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
		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 *time.Minute))
	defer conn.Close()
	request := make([]byte, 1024)
	for {
		readLen, err := conn.Read(request)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if readLen == 0 {
			break	//connection already closed by client
		} else {
			conn.Write([]byte("读取成功"))
			fmt.Println(string(request))
		}
		request = make([]byte, 1024)
	}

}

func checkError(err error)  {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}


func EpollCreate1(flag int) (fd int, err error) {
	r0, _, e1 := RawSyscall(SYS_EPOLL_CREATE1, uintptr(flag), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}