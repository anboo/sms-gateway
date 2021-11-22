package service

import (
    "github.com/anboo/sms-gateway/provider"
    "sync"
)

var globalObjectStorage *ObjectStorage

// ObjectStorage @todo need make di
type ObjectStorage struct {
    sync.Mutex
    twilioSmsProvider  *provider.TwilioSmsProvider
    vonageSmsProvider  *provider.VonageSmsProvider
    locationManager    *LocationManager
    phoneNumberManager *PhoneNumberManager
}

func (s *ObjectStorage) GetTwilioSmsProvider() *provider.TwilioSmsProvider {
    if s.twilioSmsProvider == nil {
        s.twilioSmsProvider = &provider.TwilioSmsProvider{}
        s.twilioSmsProvider.Init()
    }

    return s.twilioSmsProvider
}

func (s *ObjectStorage) GetVonageSmsProvider() *provider.VonageSmsProvider {
    if s.vonageSmsProvider == nil {
        s.vonageSmsProvider = &provider.VonageSmsProvider{}
        s.vonageSmsProvider.Init()
    }

    return s.vonageSmsProvider
}

func (s *ObjectStorage) GetLocationManager() *LocationManager {
    if s.locationManager == nil {
        s.locationManager = &LocationManager{}
    }

    return s.locationManager
}

func (s *ObjectStorage) GetPhoneNumberManager() *PhoneNumberManager {
    if s.phoneNumberManager == nil {
        s.phoneNumberManager = &PhoneNumberManager{
            Providers: []provider.SmsProvider{
                s.GetVonageSmsProvider(),
                s.GetTwilioSmsProvider(),
            },
        }
    }

    return s.phoneNumberManager
}

func GetObjectStorage() *ObjectStorage {
    if globalObjectStorage == nil {
        globalObjectStorage = &ObjectStorage{}
    }

    return globalObjectStorage
}
