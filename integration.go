package main

import (
	"fmt"
	"os"
	"net/http"
	"bytes"
	"encoding/json"
	"github.com/TheThingsNetwork/ttn/core/types"
	TTNmqtt "github.com/TheThingsNetwork/ttn/mqtt"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"

)


type ConfigOS struct {
	apiURL	string
	apiKey  string
	devicePassword string 
	deviceID string 
	userName string 
	topicName string 
}

type ConfigTTN struct {
	applicationID string
	deviceID string
	accessKey string 
}

//topic to get info from "office-app/devices/office-test/up"

func main(){
	var confOS ConfigOS
	var confTTN ConfigTTN
	//config on OpenSensors's side
	confOS.apiKey="872946ba-f2e1-4e08-b5d0-80de02966023"
	confOS.deviceID="6033"
	confOS.devicePassword="zN9mr4Pn"
	confOS.topicName="celcius"
	confOS.userName="emmlej"
	fmt.Printf(" config TTN ok ")
	//config on TTN's side
	confTTN.applicationID="office-app"
	confTTN.accessKey="ttn-account-v2.OfuuW9smtu33PjpPtVAs54Bmc2dcgHEOywtuAT1oqzk"
	confTTN.deviceID="office-hq"
	fmt.Printf(" config OS ok ")

	//apiURL complete link
	confOS.apiURL="https://realtime.opensensors.io/v1/topics//users/"+confOS.userName+"/"+confOS.topicName+"?client-id="+confOS.deviceID+"&password="+confOS.devicePassword

	//connection to mqtt client of the things network
	ctx := apex.Stdout().WithField("Integration", "TTN uplink")
	log.Set(ctx)
	clientmqtt:=TTNmqtt.NewClient(ctx, "emmlej", confTTN.applicationID, confTTN.accessKey, "tcp://eu.thethings.network:1883")
	err:=clientmqtt.Connect()
	if err!=nil{
		fmt.Println("error: connecting to the mqtt client"+err.Error())
		os.Exit(0)
	}
	fmt.Println("connected")

	//Handler using the function to post the message to OpenSensors
	handler := func(client TTNmqtt.Client, appID string, devID string, req types.UplinkMessage) {
		fmt.Print("\n*******MESSAGE INCOMING*******\n")
		response, err := confOS.postMessage(req.PayloadFields); 
		if err != nil || (response.StatusCode!=200 && response.StatusCode!=201 && response==nil) {
			fmt.Println("Error while transmitting the message")
		} else {
			fmt.Println("Your message was transmitted!")
		}
	}
	//Subscribing to the device of TTN
	token := clientmqtt.SubscribeDeviceUplink(confTTN.applicationID, confTTN.deviceID, handler)
	fmt.Println("...waiting for incoming messages...")
	token.Wait()
	if err := token.Error(); err != nil {
		fmt.Println("No subscription made " + err.Error())
		os.Exit(0)
	}
	//keeps the program running till a message arrives
	select {}
}
//function used to post the message to OpenSensors
func (confOS ConfigOS) postMessage(data map[string]interface{}) (*http.Response, error){
	DataInString, err:=json.Marshal(data)
	if err!=nil {
		return nil, err
	}

	DataTab:=make(map[string]interface{})
	DataTab["data"]=string(DataInString[:])
	message, err:=json.Marshal(DataTab)
	if err!=nil{
		return nil, err
	}

	CliHTTP := &http.Client{}
	DataInBytes := bytes.NewReader(message)
	request, err := http.NewRequest("POST", confOS.apiURL, DataInBytes)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "api-key "+ confOS.apiKey)
	return CliHTTP.Do(request)
}
