package http

import (
    "github.com/anboo/sms-gateway/service"
    "github.com/gin-gonic/gin"
)

func CheckCodeHandler(c *gin.Context) {
    req := struct {
        PhoneNumber string `json:"phone_number"`
        Code        string `json:"code"`
    }{}

    err := c.BindJSON(&req)
    if err != nil {
        c.JSON(400, gin.H{"error": "incorrect_body"})
        return
    }

    if req.PhoneNumber == "" || req.Code == "" || len(req.Code) < 4 {
        c.JSON(400, gin.H{"error": "invalid_params"})
        return
    }

    ip := c.ClientIP()
    formattedPhoneNumber, _, err := service.GetObjectStorage().GetLocationManager().ParsePhoneAndFormatE164(req.PhoneNumber, ip)
    if err != nil {
        c.JSON(400, gin.H{"error": "incorrect_phone"})
        return
    }

    phm := service.GetObjectStorage().GetPhoneNumberManager()
    for _, p := range phm.Providers {
        //p := item.(provider.SmsProvider)
        //@todo need parallel checking of code in all available providers
        //@todo need save to storage last providers for checking it later
        if p.CheckVerificationCode(formattedPhoneNumber, req.Code) {
            c.JSON(200, gin.H{"response": "OK"})
            return
        }
    }

    c.JSON(400, gin.H{"error": "incorrect_code"})
    return
}
