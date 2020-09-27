package main

import (
	"fmt"
	"github.com/wanqon/gobasics/pool"
	"time"
)

func main()  {
	t := pool.NewTask(func() error {
		fmt.Println(time.Now())
		return nil
	})
	 p := pool.NewPool(3)
	 go func() {
	 	for {
	 		p.EntryChannel <- t
		}
	 }()
	 p.Run()
}
