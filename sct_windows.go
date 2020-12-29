package sct

import (
	"syscall"
	"unsafe"
)

// setColorTemp changes the device gamma curve colors to reflect the specified color temperature.
func setColorTemp(temp int) {
	user32 := syscall.NewLazyDLL("User32.dll")
	gdi32 := syscall.NewLazyDLL("Gdi32.dll")
	procGetDC := user32.NewProc("GetDC")
	procSetDeviceGammaRamp := gdi32.NewProc("SetDeviceGammaRamp")

	if temp < 1000 || temp > 10000 {
		temp = 6500
	}

	temp -= 1000

	ratio := float64((temp-1000)%500) / 500.0
	point := whitepoints[temp/500]
	gammar := point.r*(1-ratio) + point.r*ratio
	gammag := point.g*(1-ratio) + point.g*ratio
	gammab := point.b*(1-ratio) + point.b*ratio

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
