package http

import (
    "fmt"
    "github.com/anboo/sms-gateway/service"
    "github.com/gin-gonic/gin"
)

func SendSmsHandler(c *gin.Context) {
    var err error

    req := struct {
        PhoneNumber string `json:"phone_number"`
    }{}

    err = c.BindJSON(&req)
    if err != nil {
        fmt.Println("cannot unmarshal body " + err.Error())
        c.JSON(400, gin.H{"error": "incorrect_body"})
        return
    }

    formattedPhoneNumber, regionCodeByIp, err := service.GetObjectStorage().GetLocationManager().ParsePhoneAndFormatE164(
        req.PhoneNumber,
        c.ClientIP(),
    )
    if err != nil {
        c.JSON(400, gin.H{"error": "incorrect_phone"})
        return
    }

    p, err := service.GetObjectStorage().GetPhoneNumberManager().GetProviderForPhoneNumber(formattedPhoneNumber)
    if err != nil {
        c.JSON(500, gin.H{"error": "cannot_find_provider"})
        return
    }

    reqId, err := p.SendVerificationCode(formattedPhoneNumber)
    if err != nil {
        c.JSON(500, gin.H{"error": "internal_error_at_sending_code"})
        return
    }

    c.JSON(201, gin.H{
        "reqId":                reqId,
        "formattedPhoneNumber": formattedPhoneNumber,
        "ipCountryCode":        regionCodeByIp,
    })
}
