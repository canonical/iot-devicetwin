[![Build Status][travis-image]][travis-url]
[![Go Report Card][goreportcard-image]][goreportcard-url]
[![codecov][codecov-image]][codecov-url]
# IoT Device Twin Service

The Device Twin service holds information about the state of the connected devices. Once a 
device is registered in the IoT Identity Service, a device twin it created to hold
the current and desired state of the device.

The device twin record holds:
 - Summary of the device
 - Last known heartbeat from the device
 - The groups that the device belongs to
 - Details about the device's state
 - Details about the device's desired state
 
 The service provides a cache so the devices can be monitored by the [IoT Management](https://github.com/CanonicalLtd/iot-management) 
 Service, and relays actions to the device [IoT agent](https://github.com/CanonicalLtd/iot-agent) e.g. to install a new application.
 
 ## Design
 ![IoT Management Solution Overview](docs/IoTManagement.svg)
 
 ## Build
 The project uses vendorized dependencies using `govendor`. Development has been done on minimum Go version 1.12.1.
 ```bash
 $ go get github.com/CanonicalLtd/iot-devicetwin
 $ cd iot-devicetwin
 $ ./get-deps.sh
 $ go build ./...
 ```
 
 ## Run
 ```bash
 go run cmd/devicetwin/main.go -help
  -configdir string
        Directory path to the config file (default "certs")
  -datasource string
        The data repository data source
  -driver string
        The data repository driver (default "memory")
  -mqttport string
        Port of the MQTT broker (default "8883")
  -mqtturl string
        URL of the MQTT broker (default "mqtt.example.com")
  -port string
        The port the service listens on (default "8040")
 ```
 
 The service connects to the MQTT Broker using the certificates in the `configdir` (named `ca.crt`, `server.crt` and `server.key`).
 
 ## Contributing
 Before contributing you should sign [Canonical's contributor agreement](https://www.ubuntu.com/legal/contributors), it’s the easiest way for you to give us permission to use your contributions.

[travis-image]: https://travis-ci.org/CanonicalLtd/iot-devicetwin.svg?branch=master
[travis-url]: https://travis-ci.org/CanonicalLtd/iot-devicetwin
[goreportcard-image]: https://goreportcard.com/badge/github.com/CanonicalLtd/iot-devicetwin
[goreportcard-url]: https://goreportcard.com/report/github.com/CanonicalLtd/iot-devicetwin
[codecov-url]: https://codecov.io/gh/CanonicalLtd/iot-devicetwin
[codecov-image]: https://codecov.io/gh/CanonicalLtd/iot-devicetwin/branch/master/graph/badge.svg
