package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/appened/HTTPLogger"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Config struct {
	AccountSid   string `json:"accountSid"`
	AuthToken    string `json:"authToken"`
	ClientNumber string `json:"clientNumber"`
	TwilioNumber string `json:"twilioNumber"`
	AppenedToken string `json:"appenedToken"`
	AppenedURL   string `json:"appenedURL"`
}

func main() {
	// Init Logger
	logger := HTTPLogger.New(os.Stdout, HTTPLogger.LOG_ALL)

	// Init Config
	configJSON, err := os.Open("config.json")
	if err != nil {
		logger.Error(err)
	}
	defer configJSON.Close()

	configBytes, err := ioutil.ReadAll(configJSON)
	if err != nil {
		logger.Error(err)
	}
	config := Config{}
	json.Unmarshal(configBytes, &config)
	logger.Info("Loaded config successfully")

	// Init twilio client
	client := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: config.AccountSid,
		Password: config.AuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(config.ClientNumber)
	params.SetFrom(config.TwilioNumber)
	params.SetBody("Hello from Go!")

	resp, err := client.ApiV2010.CreateMessage(params)
	if err != nil {
		logger.Error(err)
	} else {
		logger.Info("Message Sid: " + *resp.Sid)
	}
}
