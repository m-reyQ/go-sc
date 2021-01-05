package sct

import (
	"syscall"
	"unsafe"
)

// setColorTemp changes the device gamma curve colors to reflect the specified color temperature.
func setColorTemp(gammar, gammag, gammab float64) {
	user32 := syscall.NewLazyDLL("User32.dll")
	gdi32 := syscall.NewLazyDLL("Gdi32.dll")
	procGetDC := user32.NewProc("GetDC")
	procSetDeviceGammaRamp := gdi32.NewProc("SetDeviceGammaRamp")

	var gamma [3][256]uint16

	for i := 0; i < 256; i++ {
		g := 65535.0 * float64(i) / float64(256)
		gamma[0][i] = uint16(g * gammar)
		gamma[1][i] = uint16(g * gammag)
		gamma[2][i] = uint16(g * gammab)
	}

	hdc, _, _ := procGetDC.Call(uintptr(0))

	procSetDeviceGammaRamp.Call(hdc, uintptr(unsafe.Pointer(&gamma)))
}
