# go-sct [![GoDoc](https://godoc.org/github.com/d4l3k/go-sct?status.svg)](https://godoc.org/github.com/d4l3k/go-sct)

A color temperature setting library and CLI that operates in a similar way to f.lux and Redshift.

The command line app automatically determines your location using GeoIP and adjusts the color temperature depending on time.

```sh
go install github.com/d4l3k/go-sct/sct

sct 3000 // Temperature must be 1000-10000.
```
This requires Go and the Xrandr library.

## Credit
Setting the color temperature uses a port of [sct](http://www.tedunangst.com/flak/post/sct-set-color-temperature) in Go. Credit goes to him for figuring out how to do this.

go-sct also provides the `geoip` package which is a packaged version of
http://devdungeon.com/content/ip-geolocation-go

## License
go-sct is licensed under the MIT license. `geoip` and `sct` are copyrighted by their respective owners.

Made by [Tristan Rice](https://fn.lc).
