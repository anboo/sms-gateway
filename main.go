package main

import (
    "github.com/anboo/sms-gateway/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.POST("/v1/sms/verification/send", http.SendSmsHandler)
    r.POST("/v1/sms/verification/check", http.CheckCodeHandler)

    err := r.Run()
    if err != nil {
        panic(err)
    }
}
