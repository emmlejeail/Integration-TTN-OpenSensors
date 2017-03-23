package main

import (
  "encoding/json"
  "io/ioutils"
)

type ConfTTN struct {
  appID string
  accessKey string
  deviceID  string
}

type ConfOpenSensors struct {
  apiKey string
  deviceID string
  devicePassword string
  topicName string
}

type Configuration struct {
  opensensors ConfOpenSensors
  thethingsnetowrk ConfTTN
}

var config Configuration
config.thethingsnetowrk.appID:=office-app
config.thethingsnetowrk.accesskey:=ttn-account-v2.OfuuW9smtu33PjpPtVAs54Bmc2dcgHEOywtuAT1oqzk
config.thethingsnetowrk.deviceID:=office-test
config.opensensors.apiKey:=
config.opensensors.deviceID:=
config.opensensors.devicePassword:=
config.opensensors.topicName:=
