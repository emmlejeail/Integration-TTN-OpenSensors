package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/TheThingsNetwork/ttn/core/types"
	TTNmqtt "github.com/TheThingsNetwork/ttn/mqtt"
	"net/http"
)

type Config struct {
	OSapiURL         string `json: "OSapiURL"`
	OSapiKey         string `json: "OSapiKey"`
	OSdevicePassword string `json: "OSdevicePassword"`
	OSdeviceID       string `json: "OSdeviceID"`
	OSuserName       string `json: "OSuserName"`
	OStopicName      string `json: "OStopicName"`
	TTNapplicationID string `json: "TTNapplicationID"`
	TTNdeviceID      string `json: "TTNdeviceID"`
	TTNaccessKey     string `json: "TTNaccessKey"`
}

func main() {
	/*var config Config
	config.OSapiKey = "872946ba-f2e1-4e08-b5d0-80de02966023"
	config.OSdeviceID = "6033"
	config.OSdevicePassword = "zN9mr4Pn"
	config.OStopicName = "celcius"
	config.OSuserName = "emmlej"
	config.TTNapplicationID = "office-app"
	config.TTNaccessKey = "ttn-account-v2.OfuuW9smtu33PjpPtVAs54Bmc2dcgHEOywtuAT1oqzk"
	config.TTNdeviceID = "office-hq"*/
	
	var config Config
	file, err := ioutil.ReadFile("configuration.json")
	if err != nil {
		fmt.Sprintf("%s", err.Error())
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Sprintf("%s", err.Error())
	}
	
	//apiURL complete link
	config.OSapiURL = "https://realtime.opensensors.io/v1/topics//users/" + config.OSuserName + "/" + config.OStopicName + "?client-id=" + config.OSdeviceID + "&password=" + config.OSdevicePassword

	//connection to mqtt client of the things network
	clientmqtt := TTNmqtt.NewClient(nil, "emmlej", config.TTNapplicationID, config.TTNaccessKey, "tcp://eu.thethings.network:1883")
	err = clientmqtt.Connect()
	if err != nil {
		fmt.Sprintf("error: connecting to the mqtt client %s", err.Error())
	}
	fmt.Printf("connected\n")

	//Handler using the function to post the message to OpenSensors
	handler := func(client TTNmqtt.Client, appID string, devID string, req types.UplinkMessage) {
		fmt.Printf("\n*******MESSAGE INCOMING*******\n")
		response, err := config.postMessage(req.PayloadFields)
		if err != nil || (response.StatusCode != 200 && response.StatusCode != 201 && response != nil) {
			if err != nil {
				fmt.Sprintf("Error: %s", err.Error())
			} else {
				buffer := new(bytes.Buffer)
				buffer.ReadFrom(response.Body)
				fmt.Sprintf("Error : %s", buffer.String())
			}
		} else {
			fmt.Printf("Your message was transmitted!")
		}
	}
	//Subscribing to the device of TTN
	token := clientmqtt.SubscribeDeviceUplink(config.TTNapplicationID, config.TTNdeviceID, handler)
	fmt.Printf("the subscription succeded\n...waiting for incoming messages...")
	token.Wait()
	if err := token.Error(); err != nil {
		fmt.Sprintf("No subscription made %s", err.Error())
	}
	//keeps the program running till a message arrives
	select {}
}

//function used to post the message to OpenSensors
func (config Config) postMessage(data map[string]interface{}) (*http.Response, error) {
	DataInString, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	DataTab := make(map[string]interface{})
	DataTab["data"] = string(DataInString[:])
	message, err := json.Marshal(DataTab)
	if err != nil {
		return nil, err
	}

	CliHTTP := &http.Client{}
	DataInBytes := bytes.NewReader(message)
	request, err := http.NewRequest("POST", config.OSapiURL, DataInBytes)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "api-key "+config.OSapiKey)
	return CliHTTP.Do(request)
}
