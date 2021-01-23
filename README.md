# go-sct [![GoDoc](https://godoc.org/github.com/d4l3k/go-sct?status.svg)](https://godoc.org/github.com/d4l3k/go-sct)

A color temperature setting library and CLI that operates in a similar way to f.lux and Redshift.

The command line app automatically determines your location using GeoIP and adjusts the color temperature depending on time.

```sh
$ go get -u github.com/d4l3k/go-sct/sct

$ sct # Launch in background
$ sct 3000 # One time temperature change. Temperature must be 1000-10000.
```
This requires Go and (the Xrandr library or Windows).

## Windows
By default, the lowest color temperature allowed is around 4500K. More
information is available [here](http://jonls.dk/2010/09/windows-gamma-adjustments/)

There is a workaround to allow all possible adjustments by alterting the registry.

```
Windows Registry Editor Version 5.00

[HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\ICM]
"GdiIcmGammaRange"=dword:00000100
```
Save the above as a file with a ".reg" extension and double click to apply.

## Credit
Setting the color temperature uses a port of [sct](http://www.tedunangst.com/flak/post/sct-set-color-temperature) in Go. Credit goes to him for figuring out how to do this.

go-sct also provides the `geoip` package which is a packaged version of
http://devdungeon.com/content/ip-geolocation-go

## License
go-sct is licensed under the MIT license. `geoip` and `sct` are copyrighted by their respective owners.

Made by [Tristan Rice](https://fn.lc).
