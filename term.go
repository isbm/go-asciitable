package asciitable

import (
	"syscall"
	"unsafe"
)

type termSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getTerminalWidth() uint {
	size := &termSize{}
	code, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(size)))

	if int(code) == -1 {
		panic(err)
	}
	return uint(size.Col)
}
