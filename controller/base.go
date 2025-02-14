package controller

const (
	Mouse_Button_Left  string = "left"
	Mouse_Button_Right string = "right"
)

func GetMousePos() (int, int, error) {
	return getMousePos()
}

func MoveMouse(x, y int) error {
	return moveMouse(x, y)
}
func MouseClick(button string) error {
	return mouseClick(button)
}
func KeyPress(key uint8) error {
	return keyPress(key)
}
