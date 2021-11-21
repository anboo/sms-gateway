package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lesson/provider"
)

func main() {
	p := provider.TwilioSmsProvider{}
	p.Init()

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

		res := p.CheckVerificationCode(formattedPhoneNumber, req.Code)
		if res {
			c.JSON(200, gin.H{"response": "OK"})
			return
		} else {
			c.JSON(400, gin.H{"error": "incorrect_code"})
			return
		}
	})

	r.Run()

	//reqId, err := p.SendVerificationCode("+79636417683")
	//fmt.Println(reqId)
	//fmt.Println(err)

	check := p.CheckVerificationCode("+79636417683", "5067")
	fmt.Println(check)
}
