package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/appened/HTTPLogger"
	appendedGo "github.com/appened/clients/go-sdk"
	"github.com/gorilla/mux"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

// TODO: Check that X-Twilio-Signature header to ensure request is authentically from Twilio

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
	twilioClient := twilio.NewRestClientWithParams(twilio.RestClientParams{
		Username: config.AccountSid,
		Password: config.AuthToken,
	})

	// Init appended client
	appendedClient := appendedGo.New(config.AppenedToken, config.AppenedURL)

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("Error decoding request body:\n%v", err)
		}

		bodyMap, err := url.ParseQuery(string(body))
		if err != nil {
			log.Fatalf("Error converting body to map:\n%v", err)
		}

		// Get the income message and phone number of user
		incomingMsg := bodyMap["Body"][0]
		phoneNumber := bodyMap["From"][0]

		// Ensure it is the one whitelisted number
		if phoneNumber != config.ClientNumber {
			logger.Info("Incoming text from invalid number" + phoneNumber)
			return
		}

		// Create response to message
		msg, inputErr := messageResponse(incomingMsg, appendedClient)
		if inputErr != nil {
			if smsErr := sendSMS(inputErr.Error(), config, twilioClient); smsErr != nil {
				logger.Error(smsErr)
			}
			logger.Info("Error in message: " + inputErr.Error())
			return
		}

		if err = sendSMS(msg, config, twilioClient); err != nil {
			logger.Error(err)
		} else {
			logger.Info("Replied to SMS")
		}
	})
	logger.Info("Listening on port 8080")
	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting on server on ':8080':\n%v\n", err)
	}
}

func sendSMS(msg string, config Config, twilioClient *twilio.RestClient) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(config.ClientNumber)
	params.SetFrom(config.TwilioNumber)
	params.SetBody(msg)

	_, err := twilioClient.ApiV2010.CreateMessage(params)
	if err != nil {
		return err
	}
	return nil
}

func messageResponse(msg string, client *appendedGo.Client) (string, error) {
	words := strings.Split(strings.TrimSpace(msg), " ")
	if len(words) == 0 {
		return "", errors.New("Empty text message received")
	}

	cmd := strings.ToLower(words[0])

	// TODO: Find a tidier way to do this
	if len(words) == 1 {
		if cmd == "h" {
			msg := "h: this message"
			msg += "\nlf: list folios"
			msg += "\ncf <folioName>: create folio"
			msg += "\ndf <folioName>: delete folio"
			msg += "\nln <folioName>: list notes in folio"
			msg += "\nlna <folioName>: list all notes in folio, including done"
			msg += "\ndn <folioName> <number>: Toggle done on note at number"
			msg += "\na <folioName> <msg>: append note to folio"

			return msg, nil

		} else if cmd == "lf" {
			folioNames, err := client.GetFolios()
			if err != nil {
				return "", err
			}
			if len(folioNames) == 0 {
				return "No folios yet!", nil
			}
			return strings.Join(folioNames, "\n"), nil
		}
	} else if len(words) == 2 {
		folioName := words[1]
		if cmd == "cf" {
			// create folio
			if err := client.CreateFolio(folioName); err != nil {
				return "", err
			} else {
				return "Created folio with name " + folioName, nil
			}
		} else if cmd == "ln" {
			notes, err := client.GetNotes(folioName)
			if err != nil {
				return "", err
			}

			if len(notes) == 0 {
				return "No notes yet!", nil
			}

			// Remove trailing ✅
			filteredNotes := make([]string, 0)
			for _, note := range notes {
				if !strings.Contains(note, "✅") {
					filteredNotes = append(filteredNotes, note)
				}
			}

			return strings.Join(filteredNotes, "\n"), nil
		} else if cmd == "lna" {
			notes, err := client.GetNotes(folioName)
			if err != nil {
				return "", err
			}

			if len(notes) == 0 {
				return "No notes yet!", nil
			}
			return strings.Join(notes, "\n"), nil
		} else if cmd == "df" {
			if err := client.DeleteFolio(folioName); err != nil {
				return "", err
			}
			return "Deleted folio", nil
		}
	} else {
		folioName := words[1]
		note := strings.Join(words[2:], " ")
		if cmd == "a" {
			if err := client.AddNote(folioName, note); err != nil {
				return "", err
			}
			return "Appended", nil
		}
		if cmd == "dn" && len(words) == 3 {
			index, err := strconv.Atoi(words[2])
			if err != nil {
				return "", err
			}
			if err := client.ToggleDone(folioName, index-1); err != nil {
				return "", err
			}
			return "Toggled done", nil
		}
	}

	return "", errors.New("Invalid command")
}
