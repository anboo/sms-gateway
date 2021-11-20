package provider

type ResVerifyReqId string

type SmsProvider interface {
    GetProviderCode() string
    SupportPhoneNumber(phoneNumber string) bool
    SendVerificationCode(phoneNumber string) (ResVerifyReqId, error)
    CheckVerificationCode(phoneNumber string, code string) bool
}
