package geoip

import "testing"

func TestLookupIP(t *testing.T) {
	geo, err := LookupIP("")
	if err != nil {
		t.Fatal(err)
	}
	if geo.Latitude == 0 {
		t.Fatalf("expected non empty Latitude: %+v", geo)
	}
	if geo.Longitude == 0 {
		t.Fatalf("expected non empty Longitude: %+v", geo)
	}
}
