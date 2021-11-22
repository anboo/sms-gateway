package provider

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "github.com/vonage/vonage-go-sdk"
    "os"
    "sync"
)

type VonageSmsProvider struct {
    client *vonage.VerifyClient
    cache  sync.Map
}

func (p *VonageSmsProvider) Init() {
    auth := vonage.CreateAuthFromKeySecret(os.Getenv("VONAGE_API_KEY"), os.Getenv("VONAGE_API_SECRET"))
    p.client = vonage.NewVerifyClient(auth)
}

func (p *VonageSmsProvider) GetProviderCode() string {
    return "vonage"
}

func (p *VonageSmsProvider) SupportPhoneNumber(phone string) bool {
    return true
}

func (p *VonageSmsProvider) SendVerificationCode(phone string) (ResVerifyReqId, error) {
    oldReqId, ok := p.cache.Load(phone)
    if ok {
        p.client.Cancel(oldReqId.(string))
    }

    res, errResp, err := p.client.Request(phone, os.Getenv("VONAGE_BRAND_NAME"), vonage.VerifyOpts{
        CodeLength: 4,
    })

    if res.Status != "0" {
        return "", fmt.Errorf("incorrect status from vonage: " + res.Status)
    }

    if errResp.Status != "" {
        return "", fmt.Errorf("error from vonage: %s %s", errResp.Status, errResp.ErrorText)
    }

    if err != nil {
        return "", err
    }

    reqId := ResVerifyReqId(res.RequestId)

    p.cache.Store(phone, reqId)

    return reqId, nil
}

func (p *VonageSmsProvider) CheckVerificationCode(phone string, code string) bool {
    l := logrus.WithField("phone", phone)

    reqId, ok := p.cache.Load(phone)
    if !ok {
        l.Warningln("not found requestId for phone in local storage")
        return false
    }

    res, errRes, err := p.client.Check(reqId.(string), code)
    if err != nil || errRes.Status != "" {
        l.Warningln("error response checking code is " + err.Error() + " " + errRes.Status)
        return false
    }

    return res.Status == "approved"
}
