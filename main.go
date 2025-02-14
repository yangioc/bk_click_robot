package main

import (
	"bk_click_robot/controller"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
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

	// 設定快捷鍵
	if err := controller.SetHotKey(1, controller.MOD_NOREPEAT, 0x73); err != nil { // F4
		fmt.Println(err)
		return
	}
	defer controller.UnsetHotKey(1)

	// 訊號處理 (允許 Ctrl+C 退出)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 啟動熱鍵監聽
	go controller.RunHotKeyListener(func(id uintptr) {
		switch id {
		case 1:
			fmt.Println("Click")
			un := !lock.Load()
			lock.Store(un)
		}
	})

	// 等待系統訊號
	<-sigChan
	fmt.Println("程式結束")
}
