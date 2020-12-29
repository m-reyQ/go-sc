package sct

// #cgo LDFLAGS: -lX11 -lXrandr
// #include <stdio.h>
// #include <strings.h>
// #include <string.h>
// #include <stdlib.h>
// #include <stdint.h>
// #include <inttypes.h>
// #include <stdarg.h>
// #include <math.h>
// #include <X11/Xlib.h>
// #include <X11/Xlibint.h>
// #include <X11/Xproto.h>
// #include <X11/Xatom.h>
// #include <X11/extensions/Xrandr.h>
// #include <X11/extensions/Xrender.h>
// Window RootWindowMacro(Display * dpy, int scr) {
//   return RootWindow(dpy, scr);
// }
// RRCrtc crtcxid(RRCrtc * crtcs, int i) {
//	 return crtcs[i];
// }
// void ushortSet(ushort * s, int k, ushort v) {
//	 s[k] = (ushort)v;
// }
// int screenCount(Display * dpy) {
//   return XScreenCount(dpy);
// }
import "C"
import "unsafe"

// setColorTemp changes the Xrandr colors to reflect the specified color temperature.
func setColorTemp(temp int) {
	dpy := C.XOpenDisplay(nil)
	screenCount := C.screenCount(dpy)
	for screen := C.int(0); screen < screenCount; screen++ {
		root := C.RootWindowMacro(dpy, screen)

		res := C.XRRGetScreenResourcesCurrent(dpy, root)

		if temp < 1000 || temp > 10000 {
			temp = 6500
		}
		temp -= 1000
		ratio := float64((temp-1000)%500) / 500.0
		point := whitepoints[temp/500]
		gammar := point.r*(1-ratio) + point.r*ratio
		gammag := point.g*(1-ratio) + point.g*ratio
		gammab := point.b*(1-ratio) + point.b*ratio

		for c := C.int(0); c < res.ncrtc; c++ {
			crtcxid := C.crtcxid(res.crtcs, c)

			size := C.XRRGetCrtcGammaSize(dpy, crtcxid)
			crtc_gamma := C.XRRAllocGamma(size)
			for i := C.int(0); i < size; i++ {
				g := 65535.0 * float64(i) / float64(size)
				C.ushortSet(crtc_gamma.red, i, C.ushort(g*gammar))
				C.ushortSet(crtc_gamma.green, i, C.ushort(g*gammag))
				C.ushortSet(crtc_gamma.blue, i, C.ushort(g*gammab))
			}
			C.XRRSetCrtcGamma(dpy, crtcxid, crtc_gamma)
			C.XFree(unsafe.Pointer(crtc_gamma))
		}
	}
}
