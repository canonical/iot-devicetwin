// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Device Twin Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/CanonicalLtd/iot-identity/service/cert"
)

// Default settings
const (
	DefaultPort       = "8040"
	DefaultDriver     = "memory"
	DefaultDataSource = ""
	DefaultMQTTURL    = "mqtt.example.com"
	DefaultMQTTPort   = "8883"
	DefaultConfigPath = "certs"
	keyFilename       = ".secret"
	rootCA            = "ca.crt"
	clientCert        = "server.crt"
	clientKey         = "server.key"
	prefix            = "devicetwin"
)

var drivers = []string{"memory"}

// MQTTConnect holds the credentials for MQTT connection
type MQTTConnect struct {
	ClientID   string
	RootCA     []byte
	ClientCert []byte
	ClientKey  []byte
}

// Settings defines the application configuration
type Settings struct {
	Port        string
	Driver      string
	DataSource  string
	MQTTUrl     string
	MQTTPort    string
	KeySecret   string
	MQTTConnect MQTTConnect
}

// ParseArgs checks the command line arguments
func ParseArgs() *Settings {
	var (
		port       string
		driver     string
		datasource string
		mqttURL    string
		mqttPort   string
		configDir  string
	)
	flag.StringVar(&port, "port", DefaultPort, "The port the service listens on")
	flag.StringVar(&driver, "driver", DefaultDriver, "The data repository driver")
	flag.StringVar(&datasource, "datasource", DefaultDataSource, "The data repository data source")
	flag.StringVar(&mqttURL, "mqtturl", DefaultMQTTURL, "URL of the MQTT broker")
	flag.StringVar(&mqttPort, "mqttport", DefaultMQTTPort, "Port of the MQTT broker")
	flag.StringVar(&configDir, "configdir", DefaultConfigPath, "Directory path to the config file")
	flag.Parse()

	// Validate the driver
	found := false
	for i := range drivers {
		if drivers[i] == driver {
			found = true
			break
		}
	}
	if !found {
		log.Fatalf("The database driver must be one of: %s", strings.Join(drivers, ", "))
	}

	// Get/set the encryption secret
	p := path.Join(configDir, keyFilename)
	secret, err := getSecret(p)
	if err != nil {
		log.Fatalf("Error generating encryption secret: %v", err)
	}

	// Get the certificates for the MQTT broker
	m, err := readCerts(configDir)
	if err != nil {
		log.Fatalf("Error reading certificates: %v", err)
	}

	return &Settings{
		Port:        port,
		Driver:      driver,
		DataSource:  datasource,
		MQTTUrl:     mqttURL,
		MQTTPort:    mqttPort,
		KeySecret:   secret,
		MQTTConnect: m,
	}
}

func getSecret(p string) (string, error) {
	// Attempt to open the secrets file
	source, err := ioutil.ReadFile(p)
	if err == nil {
		return string(source), nil
	}

	// No secret file, so generate a secret
	s, err := cert.CreateSecret(32)
	if err != nil {
		return s, fmt.Errorf("error creating secret: %v", err)
	}

	err = ioutil.WriteFile(p, []byte(s), 0600)
	return s, err
}

// readCerts reads the certificates from the file system
func readCerts(configDir string) (MQTTConnect, error) {
	c := MQTTConnect{}
	rootCA, err := ioutil.ReadFile(path.Join(configDir, rootCA))
	if err != nil {
		return c, err
	}

	cert, err := ioutil.ReadFile(path.Join(configDir, clientCert))
	if err != nil {
		return c, err
	}

	key, err := ioutil.ReadFile(path.Join(configDir, clientKey))

	c.RootCA = rootCA
	c.ClientKey = key
	c.ClientCert = cert
	c.ClientID = generateClientID()

	return c, err
}

func generateClientID() string {
	// Generate a random string
	s, err := cert.CreateSecret(6)
	if err != nil {
		log.Printf("Error creating client ID: %v", err)
	}

	return fmt.Sprintf("%s-%s", prefix, s)
}
