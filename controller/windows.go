package controller

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010

	KEYEVENTF_KEYDOWN = 0x0000
	KEYEVENTF_KEYUP   = 0x0002

	WH_KEYBOARD_LL = 13
)

var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetCursor        = user32.NewProc("GetCursorPos")
	procSetCursor        = user32.NewProc("SetCursorPos")
	procMouseEvent       = user32.NewProc("mouse_event")
	procKeybdEvent       = user32.NewProc("keybd_event")
	procGetAsyncKey      = user32.NewProc("GetAsyncKeyState") // 加載 GetAsyncKeyState
	procCallNextHookEx   = user32.NewProc("CallNextHookEx")
	procSetWindowsHookEx = user32.NewProc("SetWindowsHookExW")
)

type POINT struct {
	X, Y int32
}

// 取得當前滑鼠坐標
func getMousePos() (int, int, error) {
	var p POINT
	r1, _, err := procGetCursor.Call(uintptr(unsafe.Pointer(&p)))
	if r1 == 0 {
		return 0, 0, err
	}
	return int(p.X), int(p.Y), nil
}

// 移動滑鼠
func moveMouse(x, y int) error {
	_, _, err := procSetCursor.Call(uintptr(x), uintptr(y))
	if err != nil {
		return err
	}
	return nil
}

// 模擬滑鼠點擊
func mouseClick(button string) error {
	switch button {
	case "left":
		r1, _, err := procMouseEvent.Call(MOUSEEVENTF_LEFTDOWN, 0, 0, 0, 0)
		if r1 == 0 {
			return err
		}
		r1, _, err = procMouseEvent.Call(MOUSEEVENTF_LEFTUP, 0, 0, 0, 0)
		if r1 == 0 {
			return err
		}
	case "right":
		r1, _, err := procMouseEvent.Call(MOUSEEVENTF_RIGHTDOWN, 0, 0, 0, 0)
		if r1 == 0 {
			return err
		}
		r1, _, err = procMouseEvent.Call(MOUSEEVENTF_RIGHTUP, 0, 0, 0, 0)
		if r1 == 0 {
			return err
		}
	default:
		return fmt.Errorf("unsupported button: %s", button)
	}
	return nil
}

// 模擬鍵盤按鍵
func keyPress(key uint8) error {
	// 模擬按鍵按下
	r1, _, err := procKeybdEvent.Call(uintptr(key), 0, KEYEVENTF_KEYDOWN, 0)
	if r1 == 0 {
		return err
	}
	// 模擬按鍵放開
	r1, _, err = procKeybdEvent.Call(uintptr(key), 0, KEYEVENTF_KEYUP, 0)
	if r1 == 0 {
		return err
	}
	return nil
}

// 鍵盤監聽回傳已按的鍵
func ListenKeyPress(key uint8, callback func(uint8)) {
	for {
		// 假設鍵盤上的字母鍵 "A" 的虛擬鍵碼是 0x41
		// 當按下 "A" 鍵時觸發回調
		if isKeyPressed(key) {
			callback(key)
		}
	}
}

// 檢查指定的虛擬鍵碼是否被按下
func isKeyPressed(vkCode uint8) bool {
	// 使用 GetAsyncKeyState 函數檢查鍵盤按鍵是否按下
	state, _, _ := procGetAsyncKey.Call(uintptr(vkCode))
	return state&0x8000 != 0
}

type KBDLLHOOKSTRUCT struct {
	VkCode   uint32
	ScanCode uint32
	Flags    uint32
	Time     uint32
	DwExtra  uintptr
}

func Hook(key uint32, callaback func()) {
	// hook, _, err := procSetWindowsHookEx.Call(uintptr(WH_KEYBOARD_LL), syscall.NewCallback(keyboardProc), 0, 0)
	// if hook == 0 {
	// 	fmt.Println("Error setting hook:", err)
	// 	return
	// }

	// cannot use func(nCode int, wParam, lParam uintptr) uintptr {…}
	// (value of type func(nCode int, wParam uintptr, lParam uintptr) uintptr) as
	// uintptr value in argument to procSetWindowsHookEx.Call

	hook, _, err := procSetWindowsHookEx.Call(uintptr(WH_KEYBOARD_LL),
		syscall.NewCallback(func(nCode int, wParam, lParam uintptr) uintptr {
			if nCode == 0 {
				keyboard := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
				fmt.Printf("Key Pressed: %d\n", keyboard.VkCode)
				if keyboard.VkCode == key {
					callaback()
				}
			}
			// 傳遞事件
			ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
			return ret
		}), 0, 0)
	if hook == 0 {
		fmt.Println("Error setting hook:", err)
		return
	}
}

// func keyboardProc(nCode int, wParam, lParam uintptr) uintptr {
// 	if nCode == 0 {
// 		keyboard := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
// 		fmt.Printf("Key Pressed: %d\n", keyboard.VkCode)
// 	}
// 	// 傳遞事件
// 	ret, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
// 	return ret
// }
