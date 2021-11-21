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
    p.client = twilio.NewRestClientWithParams(twilio.RestClientParams{AccountSid: os.Getenv("TWILIO_ACCOUNT_SID")})
}

func (p *TwilioSmsProvider) GetProviderCode() string {
    return "twilio"
}

func (p *TwilioSmsProvider) SupportPhoneNumber(phone string) bool {
    return true
}

func (p *TwilioSmsProvider) SendVerificationCode(phone string) (ResVerifyReqId, error) {
    locale := "en"

    res, err := p.client.VerifyV2.CreateVerification(p.provideServiceSid(), &openapi.CreateVerificationParams{
        Locale: &locale,
        To:     &phone,
    })

    if err != nil {
        return ResVerifyReqId(""), err
    }

    return ResVerifyReqId(*res.Sid), nil
}

func (p *TwilioSmsProvider) CheckVerificationCode(phone string, code string) bool {
    res, err := p.client.VerifyV2.CreateVerificationCheck(p.provideServiceSid(), &openapi.CreateVerificationCheckParams{
        To:   &phone,
        Code: &code,
    })

    if err != nil {
        return false
    }

    expectedStatus := "approved"

    return res.Status == &expectedStatus
}

func (p *TwilioSmsProvider) provideServiceSid() string {
    return os.Getenv("TWILIO_SERVICE_SID")
}
