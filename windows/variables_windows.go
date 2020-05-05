package windows

import "syscall"

const (
	ModAlt = 1 << iota
	ModCtrl
	ModShift
	ModWin
	ModNoRepeat = 16384
)

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	POINT  struct{ X, Y int64 }
}

type KeyState struct {
	KeyCode uint8
}

type WindowsCalls struct {
	GetMSG         *syscall.Proc
	KeyState       *syscall.Proc
	RegHotKey      *syscall.Proc
	TranslateMSG   *syscall.Proc
	DispatchMSG    *syscall.Proc
	GetKeyNameText *syscall.Proc
}

type KeyNameText struct {
	LPARAM   int
	LPSTRING uintptr
	CCHSIZE  int
}

var Windows WindowsCalls
