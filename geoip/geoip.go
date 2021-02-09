// Package geoip is a packaged version of http://devdungeon.com/content/ip-geolocation-go
package geoip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GeoIP struct {
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	ContinentCode string  `json:"continent_code"`
	ContinentName string  `json:"continent_name"`
	CountryCode   string  `json:"country_code"`
	CountryName   string  `json:"country_name"`
	RegionCode    string  `json:"region_code"`
	RegionName    string  `json:"region_name"`
	City          string  `json:"city"`
	Zip           string  `json:"zip"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Location      struct {
		GeonameID int    `json:"geoname_id"`
		Capital   string `json:"capital"`
		Languages []struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			Native string `json:"native"`
		} `json:"languages"`
		CountryFlag             string `json:"country_flag"`
		CountryFlagEmoji        string `json:"country_flag_emoji"`
		CountryFlagEmojiUnicode string `json:"country_flag_emoji_unicode"`
		CallingCode             string `json:"calling_code"`
		IsEu                    bool   `json:"is_eu"`
	} `json:"location"`
}

// APIKey is the API key to use for looking up the ip addresses.
var APIKey = "92dd0c8e9da13e9d292af57c86ba348e"

var (
	address  string
	err      error
	geo      GeoIP
	response *http.Response
	body     []byte
)

// LookupIP looks up the geolocation information for the specified address ("" for current host).
func LookupIP(address string) (*GeoIP, error) {
	if address == "" {
		address = "check"
	}
	response, err = http.Get(fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s&output=json", address, APIKey))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// response.Body() is a reader type. We have
	// to use ioutil.ReadAll() to read the data
	// in to a byte slice(string)
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON byte slice to a GeoIP struct
	err = json.Unmarshal(body, &geo)
	if err != nil {
		return nil, err
	}

	return &geo, nil
}
