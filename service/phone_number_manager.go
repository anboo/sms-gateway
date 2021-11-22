package service

import (
    "fmt"
    "github.com/anboo/sms-gateway/provider"
)

type PhoneNumberManager struct {
    Providers []provider.SmsProvider
}

func (m *PhoneNumberManager) GetProviderForPhoneNumber(phoneNumber string) (provider.SmsProvider, error) {
    for _, item := range m.Providers {
        p := item.(provider.SmsProvider)
        if p.SupportPhoneNumber(phoneNumber) {
            return p, nil
        }
    }

    return nil, fmt.Errorf("provider for phone number not found")
}
