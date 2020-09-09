package netio

import (
	"os"
	"syscall"
)

func nioServer() (int, error) {
	syscall.ForkLock.RLock()
	//创建socket
	s, err := syscall.Socket(1,1,1)
	if err == nil {
		//关闭从父线程拷贝过来的文件描述符后，再执行子线程程序
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return -1, os.NewSyscallError("socket", err)
	}
	//设置socket为非阻塞
	if err = syscall.SetNonblock(s, true); err != nil {
		syscall.Close(s)
		return -1, os.NewSyscallError("setnonblock", err)
	}
	return s, nil
}
