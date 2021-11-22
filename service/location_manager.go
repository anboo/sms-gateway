package service

import (
    "fmt"
    "github.com/oschwald/geoip2-golang"
    "github.com/ttacon/libphonenumber"
    "net"
)

type LocationManager struct {
    //local database pointer
    //can be nil if not called LocationManager.UseDatabase(path)
    //if nil, location manager not parse ip for detect iso code and use libphonenumber.UNKNOWN_REGION
    mixmind *geoip2.Reader
}

func (l *LocationManager) UseDatabase(databasePath string) error {
    var err error
    l.mixmind, err = geoip2.Open(databasePath)

    if err != nil {
        return err
    }

    return nil
}

func (l *LocationManager) ParsePhoneAndFormatE164(phoneNumber string, clientIp string) (string, string, error) {
    regionCodeByIp := libphonenumber.UNKNOWN_REGION

    ip := net.ParseIP(clientIp)
    if ip != nil && l.mixmind != nil {
        country, err := l.mixmind.Country(ip)
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
