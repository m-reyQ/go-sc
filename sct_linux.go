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
func setColorTemp(gammar, gammag, gammab float64) {
	dpy := C.XOpenDisplay(nil)
	screenCount := C.screenCount(dpy)
	for screen := C.int(0); screen < screenCount; screen++ {
		root := C.RootWindowMacro(dpy, screen)

		res := C.XRRGetScreenResourcesCurrent(dpy, root)

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
