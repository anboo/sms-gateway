package main

import (
    "fmt"
    "lesson/provider"
)

func main() {
    p := provider.VonageSmsProvider{}
    p.Init()

    reqId := p.CheckVerificationCode("+79636417683", "8775")
    fmt.Println(reqId)
}
