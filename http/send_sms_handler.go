package http

import (
    "github.com/anboo/sms-gateway/service"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

type SendSmsHandler struct {
    PhoneNumberManager *service.PhoneNumberManager
    LocationManager    *service.LocationManager
}

func (e *SendSmsHandler) Handler(c *gin.Context) {
    var err error

    req := struct {
        PhoneNumber string `json:"phone_number"`
    }{}

    err = c.BindJSON(&req)
    if err != nil {
        logrus.Warningln("cannot unmarshal body " + err.Error())
        c.JSON(400, gin.H{"error": "incorrect_body"})
        return
    }

    formattedPhoneNumber, regionCodeByIp, err := e.LocationManager.ParsePhoneAndFormatE164(
        req.PhoneNumber,
        c.ClientIP(),
    )
    if err != nil {
        logrus.Infoln("incorrect phone " + err.Error())
        c.JSON(400, gin.H{"error": "incorrect_phone"})
        return
    }

    //todo add checking number is libphonenumber.FIXED_LINE_OR_MOBILE
    //if libphonenumber.GetNumberType(formattedPhoneNumber) == libphonenumber.FIXED_LINE_OR_MOBILE {
    //}

    p, err := e.PhoneNumberManager.GetProviderForPhoneNumber(formattedPhoneNumber)
    if err != nil {
        logrus.Error("cannot find provider " + err.Error() + " for " + formattedPhoneNumber)
        c.JSON(500, gin.H{"error": "cannot_find_provider"})
        return
    }

    reqId, err := p.SendVerificationCode(formattedPhoneNumber)
    if err != nil {
        logrus.Error("error send verification code" + err.Error() + " for " + formattedPhoneNumber)
        c.JSON(500, gin.H{"error": "internal_error_at_sending_code"})
        return
    }

    c.JSON(201, gin.H{
        "reqId":                reqId,
        "formattedPhoneNumber": formattedPhoneNumber,
        "ipCountryCode":        regionCodeByIp,
    })
}
