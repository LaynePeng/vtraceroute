package main

import (
	"github.com/fiorix/freegeoip"
	"log"
	"net"
	"time"
)

const maxmindFile = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"

type customQuery struct {
	Country struct {
		ISOCode string            `maxminddb:"iso_code"`
		Names   map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}

func lookupIP(ip string) (float64, float64) {
	updateInterval := 7 * 24 * time.Hour // week
	maxRetryInterval := time.Hour
	db, err := freegeoip.OpenURL(maxmindFile, updateInterval, maxRetryInterval)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	select {
	case <-db.NotifyOpen():
		log.Println("Reading from the geoip file")
	case err := <-db.NotifyError():
		log.Fatal(err)
	}
	var result customQuery
	err = db.Lookup(net.ParseIP(ip), &result)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Approximate location is %v, %v\n", result.Location.Latitude, result.Location.Longitude)
	return result.Location.Latitude, result.Location.Longitude
}
