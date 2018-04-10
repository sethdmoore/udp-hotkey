package windows

import "syscall"

func init() {

	// It's safe to put everything in init as golang will panic if any of thse fail
	user32 := syscall.MustLoadDLL("user32")

	Windows.RegHotKey = user32.MustFindProc("RegisterHotKey")

	Windows.KeyState = user32.MustFindProc("GetAsyncKeyState")
	Windows.GetMSG = user32.MustFindProc("GetMessageW")
	Windows.TranslateMSG = user32.MustFindProc("TranslateMessage")
	Windows.DispatchMSG = user32.MustFindProc("DispatchMessageW")
	Windows.GetKeyNameText = user32.MustFindProc("GetKeyNameTextW")

	defer user32.Release()
}

func Get() WindowsCalls {
	return Windows
}
