package main

import (
	"fmt"
	"github.com/anboo/sms-gateway/provider"
	"github.com/gin-gonic/gin"
)

var providers map[string]interface{}

func getProviderForPhoneNumber(phoneNumber string) (provider.SmsProvider, error) {
	for _, item := range providers {
		p := item.(provider.SmsProvider)
		if p.SupportPhoneNumber(phoneNumber) {
			return p, nil
		}
	}

	return nil, fmt.Errorf("provider for phone number not found")
}

func main() {
	providers = map[string]interface{}{
		"vonage": &provider.VonageSmsProvider{},
		"twilio": &provider.TwilioSmsProvider{},
	}

	for _, item := range providers {
		p := item.(provider.SmsProvider)
		p.Init()
	}

	r := gin.Default()
	r.POST("/v1/sms/verification/send", func(c *gin.Context) {
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

		formattedPhoneNumber, regionCodeByIp, err := ParsePhoneAndFormatE164(req.PhoneNumber, c.ClientIP())
		if err != nil {
			c.JSON(400, gin.H{"error": "incorrect_phone"})
			return
		}

		p, err := getProviderForPhoneNumber(formattedPhoneNumber)
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
	})

	r.POST("/v1/sms/verification/check", func(c *gin.Context) {
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

		formattedPhoneNumber, _, err := ParsePhoneAndFormatE164(req.PhoneNumber, c.ClientIP())
		if err != nil {
			c.JSON(400, gin.H{"error": "incorrect_phone"})
			return
		}

		for _, item := range providers {
			p := item.(provider.SmsProvider)
			//@todo need parallel checking of code in all available providers
			//@todo need save to storage last providers for checking it later
			if p.CheckVerificationCode(formattedPhoneNumber, req.Code) {
				c.JSON(200, gin.H{"response": "OK"})
				return
			}
		}

		c.JSON(400, gin.H{"error": "incorrect_code"})
		return
	})

	err := r.Run(); if err != nil {
		panic(err)
	}
}
