package main

import (
	"fmt"
	"os"
	"github.com/TheThingsNetwork/ttn/core/types"
	TTNmqtt "github.com/TheThingsNetwork/ttn/mqtt"
	"github.com/TheThingsNetwork/go-utils/log"
	"github.com/TheThingsNetwork/go-utils/log/apex"

)


type ConfigOS struct {
	apiKey         string 
	apiURL         string 
	deviceID       string 
	devicePassword string 
	topicName      string 
	userName       string 
}

type ConfigTTN struct {
	accessKey     string 
	applicationID string 
	deviceID      string 
	//region        string 
}

type Conf struct {
	OS ConfigOS 
	TTN ConfigTTN         
}

func main(){
	var config Conf
	config.OS.apiKey="872946ba-f2e1-4e08-b5d0-80de02966023"
	config.OS.apiURL="https://realtime.opensensors.io/v1/"
	config.OS.deviceID="6033"
	config.OS.devicePassword="zN9mr4Pn"
	config.OS.topicName="/users/emmlej/ttopic"
	config.OS.userName="emmlej"

	config.TTN.applicationID="office-app"
	config.TTN.accessKey="ttn-account-v2.OfuuW9smtu33PjpPtVAs54Bmc2dcgHEOywtuAT1oqzk"
	config.TTN.deviceID="office-test"
	//config.TTN.region="eu"
	fmt.Printf("ok")
	//connection to mqtt client of the things network
	ctx := apex.Stdout().WithField("Example", "Go Client")
	log.Set(ctx)
	clientmqtt:=TTNmqtt.NewClient(ctx, "emmlej", config.TTN.applicationID, config.TTN.accessKey, "tcp://eu.thethings.network:1883")
	err:=clientmqtt.Connect()
	if err!=nil{
		fmt.Println("error: connecting to the mqtt client"+err.Error())
		os.Exit(0)
	}
	fmt.Println("connected")

	//connection to open sensors client
	var clientOpenSensor ConfigOS
	clientOpenSensor=config.OS

	handler:=func(client TTNmqtt.Client, appID string, devID string, req types.UplinkMessage)

	token := clientmqtt.SubscribeDeviceUplink(config.TTN.applicationID, config.TTN.deviceID, func(client TTNmqtt.Client, appID string, devID string, req types.UplinkMessage) {
		// Do something with the uplink message
	})
	token.Wait()
	if err := token.Error(); err != nil {
		fmt.Printf("could not subscribe")
	}

}

