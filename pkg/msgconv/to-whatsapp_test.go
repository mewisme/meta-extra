package msgconv

import (
	"math"
	"testing"

	"go.mau.fi/whatsmeow/proto/waArmadilloXMA"
)

func TestParseGeoURI(t *testing.T) {
	tests := []struct {
		name  string
		uri   string
		lat   float64
		long  float64
		valid bool
	}{
		{name: "coordinates", uri: "geo:10.123456,106.654321", lat: 10.123456, long: 106.654321, valid: true},
		{name: "parameters", uri: "geo:-90,-180;u=10", lat: -90, long: -180, valid: true},
		{name: "missing_prefix", uri: "10,106"},
		{name: "missing_longitude", uri: "geo:10"},
		{name: "invalid_latitude", uri: "geo:x,106"},
		{name: "invalid_longitude", uri: "geo:10,x"},
		{name: "latitude_range", uri: "geo:90.1,106"},
		{name: "longitude_range", uri: "geo:10,-180.1"},
		{name: "latitude_nan", uri: "geo:NaN,106"},
		{name: "longitude_inf", uri: "geo:10,+Inf"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lat, long, err := parseGeoURI(test.uri)
			if test.valid {
				if err != nil || math.Abs(lat-test.lat) > 1e-9 || math.Abs(long-test.long) > 1e-9 {
					t.Fatalf("parseGeoURI(%q) = %f,%f,%v", test.uri, lat, long, err)
				}
			} else if err == nil {
				t.Fatalf("parseGeoURI(%q) unexpectedly succeeded", test.uri)
			}
		})
	}
}

func TestLocationXMA(t *testing.T) {
	xma := locationXMA(10.123456, 106.654321)
	if xma.GetTargetType() != waArmadilloXMA.ExtendedContentMessage_MSG_LOCATION_SHARING_V2 || xma.GetXmaLayoutType() != waArmadilloXMA.ExtendedContentMessage_SINGLE {
		t.Fatalf("unexpected location XMA type: %v %v", xma.GetTargetType(), xma.GetXmaLayoutType())
	}
	if len(xma.GetCtas()) != 1 || xma.GetCtas()[0].GetButtonType() != waArmadilloXMA.ExtendedContentMessage_OPEN_NATIVE {
		t.Fatalf("unexpected location CTA: %+v", xma.GetCtas())
	}
	if got := xma.GetCtas()[0].GetNativeURL(); got != "messenger://location_share?lat=10.123456&long=106.654321" {
		t.Fatalf("unexpected location URL: %q", got)
	}
}
