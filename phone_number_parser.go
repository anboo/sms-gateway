package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/ttacon/libphonenumber"
	"net"
	"os"
)

var mixmind *geoip2.Reader

func ParsePhoneAndFormatE164(phoneNumber string, clientIp string) (string, string, error) {
	mixmindDatabasePath := os.Getenv("MIXMIND_DATABASE_PATH")

	if mixmindDatabasePath != "" {
		var err error
		mixmind, err = geoip2.Open(mixmindDatabasePath)
		if err != nil {
			panic(err)
		}
	}

	regionCodeByIp := libphonenumber.UNKNOWN_REGION

	ip := net.ParseIP(clientIp)
	if ip != nil && mixmind != nil {
		country, err := mixmind.Country(ip)
		if err != nil {
			fmt.Println("Cannot detect country for ip " + ip.String())
		} else {
			regionCodeByIp = country.Country.IsoCode
		}
	}

	parsed, err := libphonenumber.Parse(phoneNumber, regionCodeByIp)
	if err != nil {
		return "", "", err
	}

	return libphonenumber.Format(parsed, libphonenumber.E164), regionCodeByIp, nil
}
