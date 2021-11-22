package provider

type ResVerifyReqId string

type SmsProvider interface {
    Init()
    GetProviderCode() string
    SupportPhoneNumber(phoneNumber string) bool
    SendVerificationCode(phoneNumber string) (string, error)
    CheckVerificationCode(phoneNumber string, code string) bool
}
