package provider

import (
    "github.com/twilio/twilio-go"
    openapi "github.com/twilio/twilio-go/rest/verify/v2"
    "os"
)

type TwilioSmsProvider struct {
    client *twilio.RestClient
}

func (p *TwilioSmsProvider) Init() {
    p.client = twilio.NewRestClientWithParams(twilio.RestClientParams{
        AccountSid: os.Getenv("TWILIO_ACCOUNT_SID"),
        Password:   os.Getenv("TWILIO_AUTH_TOKEN"),
    })
}

func (p *TwilioSmsProvider) GetProviderCode() string {
    return "twilio"
}

func (p *TwilioSmsProvider) SupportPhoneNumber(phone string) bool {
    return true
}

func (p *TwilioSmsProvider) SendVerificationCode(phone string) (string, error) {
    locale := "en"
    channel := "sms"

    res, err := p.client.VerifyV2.CreateVerification(p.provideServiceSid(), &openapi.CreateVerificationParams{
        Locale:  &locale,
        To:      &phone,
        Channel: &channel,
    })

    if err != nil {
        return "", err
    }

    return *res.Sid, nil
}

func (p *TwilioSmsProvider) CheckVerificationCode(phone string, code string) bool {
    res, err := p.client.VerifyV2.CreateVerificationCheck(p.provideServiceSid(), &openapi.CreateVerificationCheckParams{
        To:   &phone,
        Code: &code,
    })

    if err != nil {
        return false
    }

    return *res.Status == "approved"
}

func (p *TwilioSmsProvider) provideServiceSid() string {
    return os.Getenv("TWILIO_SERVICE_SID")
}
