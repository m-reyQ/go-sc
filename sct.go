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
// int DefaultScreenMacro(Display * dpy) {
//   return DefaultScreen(dpy);
// }
// Window RootWindowMacro(Display * dpy, int scr) {
//   return RootWindow(dpy, scr);
// }
// RRCrtc crtcxid(RRCrtc * crtcs, int i) {
//	 return crtcs[i];
// }
// void ushortSet(ushort * s, int k, double v) {
//	 s[k] = (ushort)v;
// }
import "C"
import "unsafe"

type color struct {
	r, g, b float64
}

/* cribbed from redshift, but truncated with 500K steps */
var whitepoints = []color{
	{1.00000000, 0.18172716, 0.00000000}, /* 1000K */
	{1.00000000, 0.42322816, 0.00000000},
	{1.00000000, 0.54360078, 0.08679949},
	{1.00000000, 0.64373109, 0.28819679},
	{1.00000000, 0.71976951, 0.42860152},
	{1.00000000, 0.77987699, 0.54642268},
	{1.00000000, 0.82854786, 0.64816570},
	{1.00000000, 0.86860704, 0.73688797},
	{1.00000000, 0.90198230, 0.81465502},
	{1.00000000, 0.93853986, 0.88130458},
	{1.00000000, 0.97107439, 0.94305985},
	{1.00000000, 1.00000000, 1.00000000}, /* 6500K */
	{0.95160805, 0.96983355, 1.00000000},
	{0.91194747, 0.94470005, 1.00000000},
	{0.87906581, 0.92357340, 1.00000000},
	{0.85139976, 0.90559011, 1.00000000},
	{0.82782969, 0.89011714, 1.00000000},
	{0.80753191, 0.87667891, 1.00000000},
	{0.78988728, 0.86491137, 1.00000000}, /* 10000K */
	{0.77442176, 0.85453121, 1.00000000},
}

func SetColorTemp(temp int) {
	dpy := C.XOpenDisplay(nil)
	screen := C.DefaultScreenMacro(dpy)
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
			C.ushortSet(crtc_gamma.red, i, C.double(g*gammar))
			C.ushortSet(crtc_gamma.green, i, C.double(g*gammag))
			C.ushortSet(crtc_gamma.blue, i, C.double(g*gammab))
		}
		C.XRRSetCrtcGamma(dpy, crtcxid, crtc_gamma)
		C.XFree(unsafe.Pointer(crtc_gamma))
	}
}
