package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/cpucycle/astrotime"
	"github.com/d4l3k/go-sct"
	"github.com/d4l3k/go-sct/geoip"
)

var dayTemp = flag.Int("dayTemp", 6500, "The color temperature during the day.")
var nightTemp = flag.Int("nightTemp", 3000, "The color temperature during the day.")
var daemon = flag.Bool("d", true, "run app as a daemon")

func monitorTime() {
	log.Printf("Fetching location...")
	geo, err := geoip.LookupIP("")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" - City: %s, Lat: %f, Lon: %f", geo.City, geo.Lat, geo.Lon)
	log.Printf("Monitoring daylight settings...")
	var lastState *bool
	for {
		rise := astrotime.NextSunrise(time.Now(), geo.Lat, -geo.Lon)
		set := astrotime.NextSunset(time.Now(), geo.Lat, -geo.Lon)
		state := rise.Before(set)
		if lastState != nil && state == *lastState {
			time.Sleep(1 * time.Minute)
			continue
		}
		lastState = &state
		if state {
			log.Print("Good night!")
			sct.SetColorTemp(*nightTemp)
		} else {
			log.Print("Good morning!")
			sct.SetColorTemp(*dayTemp)
		}
		time.Sleep(1 * time.Minute)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		if temp, err := strconv.Atoi(args[0]); err == nil {
			sct.SetColorTemp(temp)
		}
	} else if len(args) == 0 {
		if *daemon {
			args := os.Args[1:]
			for i := 0; i < len(args); i++ {
				if strings.HasPrefix(args[i], "-d") {
					args[i] = "-d=false"
					break
				}
			}
			args = append(args, "-d=false")
			cmd := exec.Command(os.Args[0], args...)
			cmd.Start()
			log.Println("Launched background process... pid", cmd.Process.Pid)
			os.Exit(0)
		} else {
			monitorTime()
		}
	}
}
