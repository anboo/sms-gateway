package main

import (
	"fmt"
	"lesson/provider"
)

func main() {
	p := provider.TwilioSmsProvider{}
	p.Init()

	//reqId, err := p.SendVerificationCode("+79636417683")
	//fmt.Println(reqId)
	//fmt.Println(err)

	check := p.CheckVerificationCode("+79636417683", "5067")
	fmt.Println(check)
}
