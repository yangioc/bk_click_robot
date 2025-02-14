package main

import (
	"bk_click_robot/controller"
	"fmt"
	"sync/atomic"
	"time"
)

// 請根據實際路徑修改
var lock atomic.Bool

func main() {
	lock.Store(false)
	go func() {
		for {
			if lock.Load() {
				controller.MouseClick(controller.Mouse_Button_Left)
				time.Sleep(time.Millisecond)
			}
		}
	}()

	controller.Hook(0x73, func() {
		fmt.Println("Click")
		un := !lock.Load()
		lock.Store(un)
	})
	fmt.Println("finish")
	select {}
}
