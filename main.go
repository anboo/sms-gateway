package main

import (
	"fmt"
	"github.com/ttacon/libphonenumber"
	"sms-gateway/provider"
	"github.com/gin-gonic/gin"
)

func main() {
	p := provider.TwilioSmsProvider{}
	p.Init()

	r := gin.Default()
	r.POST("/v1/sms/verification", func(c *gin.Context) {
		req := struct {
			PhoneNumber string `json:"phone_number"`
		}{}

		err := c.BindJSON(req); if err != nil {
			c.JSON(400, gin.H{"error": "incorrect_body"})
		}

		phone, err := libphonenumber.Parse(req.PhoneNumber, ""); if err != nil {
			c.JSON(400, gin.H{"error": "incorrect_phone"})
		}

		reqId, err := p.SendVerificationCode(libphonenumber.Format(phone, libphonenumber.E164)); if err != nil {
			c.JSON(500, gin.H{"error": "internal_error_at_sending_code"})
		}

		c.JSON(201, gin.H{
			"reqId": reqId,
			"regionCode": "",
			"formattedPhoneNumber": "",
			"ipCountryCode": "",
			"phoneNumberCountryCode": phone.CountryCodeSource.String(),
		})
	})

	r.Run()

	//reqId, err := p.SendVerificationCode("+79636417683")
	//fmt.Println(reqId)
	//fmt.Println(err)

	check := p.CheckVerificationCode("+79636417683", "5067")
	fmt.Println(check)
}
