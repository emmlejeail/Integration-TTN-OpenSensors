package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/TheThingsNetwork/ttn/core/types"
	TTNmqtt "github.com/TheThingsNetwork/ttn/mqtt"
	"io/ioutil"
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
	TTNregion        string `json: "eu"`
}

func main() {

	var config Config
	file, err := ioutil.ReadFile("configuration.json")
	if err != nil {
		errRead := fmt.Sprintf("Error: %s", err.Error())
		fmt.Printf(errRead)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		errConv := fmt.Sprintf("Error: %s", err.Error())
		fmt.Printf(errConv)
	}

	//connection to mqtt client of the things network
	clientmqtt := TTNmqtt.NewClient(nil, "emmlej", config.TTNapplicationID, config.TTNaccessKey, "tcp://"+config.TTNregion+".thethings.network:1883")
	err = clientmqtt.Connect()
	if err != nil {
		err1 := fmt.Sprintf("%s", err.Error())
		fmt.Printf(err1)
	} else {
		fmt.Printf("connected\n")
	}

	//Handler using the function to post the message to OpenSensors
	handler := func(client TTNmqtt.Client, appID string, devID string, req types.UplinkMessage) {
		fmt.Printf("\n*******MESSAGE INCOMING*******\n")
		response, err := config.postMessage(req.PayloadFields)
		if err != nil || (response.StatusCode != 200 && response.StatusCode != 201 && response != nil) {
			if err != nil {
				errPost := fmt.Sprintf("Error: %s", err.Error())
				fmt.Printf(errPost)
			} else {
				buffer := new(bytes.Buffer)
				buffer.ReadFrom(response.Body)
				errResp := fmt.Sprintf("Error : %s", buffer.String())
				fmt.Printf(errResp)
			}
		} else {
			fmt.Printf("Your message was transmitted!")
		}
	}
	//Subscribing to the device of TTN
	token := clientmqtt.SubscribeDeviceUplink(config.TTNapplicationID, config.TTNdeviceID, handler)
	fmt.Printf("...waiting for incoming messages...")
	token.Wait()
	if err := token.Error(); err != nil {
		errSub := fmt.Sprintf("No subscription made %s", err.Error())
		fmt.Printf(errSub)
	}
	//keeps the program running till a message arrives
	select {}
}

//function used to get the apiURL for OpenSensors
func (config Config) getapiURL() string {
	return config.OSapiURL + "topics//users/" + config.OSuserName + "/" + config.OStopicName + "?client-id=" + config.OSdeviceID + "&password=" + config.OSdevicePassword
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
	request, err := http.NewRequest("POST", config.getapiURL(), DataInBytes)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "api-key "+config.OSapiKey)
	return CliHTTP.Do(request)
}
