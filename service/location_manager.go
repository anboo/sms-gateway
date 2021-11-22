package service

import (
    "fmt"
    "github.com/oschwald/geoip2-golang"
    "github.com/ttacon/libphonenumber"
    "net"
    "os"
    "sync"
)

type LocationManager struct {
    sync.Mutex
    mixmind *geoip2.Reader
}

func (l *LocationManager) ParsePhoneAndFormatE164(phoneNumber string, clientIp string) (string, string, error) {
    mixmindDatabasePath := os.Getenv("MIXMIND_DATABASE_PATH")

    l.Lock()
    if mixmindDatabasePath != "" && l.mixmind == nil {
        var err error
        l.mixmind, err = geoip2.Open(mixmindDatabasePath)
        if err != nil {
            panic(err)
        }
    }
    l.Unlock()

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
