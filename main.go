package main

import (
	"github.com/anboo/sms-gateway/http"
	"github.com/anboo/sms-gateway/provider"
	"github.com/anboo/sms-gateway/service"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	twilioProvider *provider.TwilioSmsProvider
	vonageProvider *provider.VonageSmsProvider
	phm            *service.PhoneNumberManager
	lm             *service.LocationManager
)

func main() {
	r := gin.Default()

	twilioProvider = &provider.TwilioSmsProvider{}
	twilioProvider.Init()

	vonageProvider = &provider.VonageSmsProvider{}
	vonageProvider.Init()

	phm = &service.PhoneNumberManager{
		Providers: []provider.SmsProvider{
			vonageProvider,
			twilioProvider,
		},
	}

	lm = &service.LocationManager{}
	mixmindDatabasePath := os.Getenv("MIXMIND_DATABASE_PATH")
	if mixmindDatabasePath != "" {
		err := lm.UseDatabase(mixmindDatabasePath)
		if err != nil {
			panic(err)
		}
	}

	log.SetOutput(os.Stdout)

	sendSmsHandler := http.SendSmsHandler{PhoneNumberManager: phm, LocationManager: lm}
	r.POST("/v1/sms/verification/send", sendSmsHandler.Handler)
	checkCodeHandler := http.CheckCodeHandler{PhoneNumberManager: phm, LocationManager: lm}
	r.POST("/v1/sms/verification/check", checkCodeHandler.Handler)
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			Timeout:             time.Millisecond * 500,
			MaxRequestsInFlight: 100,
		},
	)))
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
