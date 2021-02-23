package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cpucycle/astrotime"
	"github.com/d4l3k/go-sct"
	"github.com/d4l3k/go-sct/geoip"
)

var dayTemp = flag.Int("dayTemp", 6500, "The color temperature during the day.")
var nightTemp = flag.Int("nightTemp", 3000, "The color temperature during the night.")
var daemon = flag.Bool("d", true, "run app as a daemon")
var mode = flag.String("mode", "geoip", "Mode of daemon (geoip or timed). Timed mode uses specified sunrise-time, midday-time, and sunset-time.")
var sunriseTimeStr = flag.String("sunrise-time", "06:30", "Sunrise time (HH:MM)")
var sunsetTimeStr = flag.String("sunset-time", "21:00", "Sunset time (HH:MM)")
var middayTimeStr = flag.String("midday-time", "14:00", "Mid day (brightest) time (HH:MM)")
var midnight, sunriseTime, sunsetTime, middayTime time.Time

func monitorGeo() {
	log.Printf("Fetching location...")
	geo, err := geoip.LookupIP("")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" - City: %s, Lat: %f, Lon: %f", geo.City, geo.Latitude, geo.Longitude)
	log.Printf("Monitoring daylight settings...")
	var lastState *bool
	for {
		rise := astrotime.NextSunrise(time.Now(), geo.Latitude, -geo.Longitude)
		set := astrotime.NextSunset(time.Now(), geo.Latitude, -geo.Longitude)
		state := rise.Before(set)
		if lastState != nil && state == *lastState {
			time.Sleep(1 * time.Minute)
			continue
		}
		lastState = &state
		if state {
			log.Print("Good night!")
			if err := interpolateColorTemp(*nightTemp); err != nil {
				log.Fatalf("%+v", err)
			}
		} else {
			log.Print("Good morning!")
			if err := interpolateColorTemp(*dayTemp); err != nil {
				log.Fatalf("%+v", err)
			}
		}
	}
}

var (
	totalTime = 3 * time.Second
	stepEvery = 1 * time.Second / 60
)

func interpolateColorTemp(new int) error {
	old, err := getCurrentColorTemp()
	if err != nil {
		return err
	}

	steps := int(totalTime / stepEvery)
	stepSize := (new - old) / steps
	for i := 0; i < steps; i++ {
		c := time.After(stepEvery)
		sct.SetColorTemp(old + stepSize*i)
		<-c
	}
	sct.SetColorTemp(new)

	return saveCurrentColorTemp(new)
}

func tempFile() string {
	display := os.Getenv("DISPLAY")
	return filepath.Join(os.TempDir(), "sct-temp-"+display)
}

func saveCurrentColorTemp(temp int) error {
	return ioutil.WriteFile(tempFile(), []byte(strconv.Itoa(temp)), 0644)
}

func getCurrentColorTemp() (int, error) {
	body, err := ioutil.ReadFile(tempFile())
	if os.IsNotExist(err) {
		return *dayTemp, nil
	} else if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(body))
}

func monitorTime() {
	var monitorTemperature int
	monitorTemperature = 6500

	for {
		curTime := time.Now()
		midnight = time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.Local)

		// Advance the day?
		if midnight.After(sunsetTime) {
			sunriseTime = sunriseTime.AddDate(0, 0, 1)
			middayTime = middayTime.AddDate(0, 0, 1)
			sunsetTime = sunsetTime.AddDate(0, 0, 1)
		}

		if curTime.After(sunriseTime) && curTime.Before(sunsetTime) {
			if curTime.Before(middayTime) {
				elapsedDuration := curTime.Sub(sunriseTime)
				ratio := float64(elapsedDuration) / float64(middayTime.Sub(sunriseTime))
				monitorTemperature = int(float64(*nightTemp)*(1-ratio) + float64(*dayTemp)*ratio)
			} else {
				elapsedDuration := curTime.Sub(middayTime)
				ratio := float64(elapsedDuration) / float64(sunsetTime.Sub(middayTime))
				monitorTemperature = int(float64(*dayTemp)*(1-ratio) + float64(*nightTemp)*ratio)
			}
		} else {
			monitorTemperature = *nightTemp
		}
		sct.SetColorTemp(monitorTemperature)
		time.Sleep(10 * time.Minute)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		if temp, err := strconv.Atoi(args[0]); err == nil {
			if err := interpolateColorTemp(temp); err != nil {
				log.Fatalf("%+v", err)
			}
		}
	} else if len(args) == 0 {
		// Parse time arguments
		curTime := time.Now()
		midnight = time.Date(curTime.Year(), curTime.Month(), curTime.Day(), 0, 0, 0, 0, time.Local)
		var perr error
		if sunriseTime, perr = time.ParseInLocation("2006-01-02 15:04", midnight.Format("2006-01-02 ")+*sunriseTimeStr, time.Local); perr != nil {
			log.Fatal(perr)
		}
		if sunsetTime, perr = time.ParseInLocation("2006-01-02 15:04", midnight.Format("2006-01-02 ")+*sunsetTimeStr, time.Local); perr != nil {
			log.Fatal(perr)
		}
		if middayTime, perr = time.ParseInLocation("2006-01-02 15:04", midnight.Format("2006-01-02 ")+*middayTimeStr, time.Local); perr != nil {
			log.Fatal(perr)
		}

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
			switch *mode {
			case "timed":
				monitorTime()
			default:
				monitorGeo()
			}
		}
	}
}
