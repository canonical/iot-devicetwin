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

package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/everactive/iot-devicetwin/config"
	"github.com/everactive/iot-devicetwin/service/controller"
	"github.com/everactive/iot-devicetwin/service/devicetwin"
	"github.com/everactive/iot-devicetwin/service/factory"
	"github.com/everactive/iot-devicetwin/service/mqtt"
	"github.com/everactive/iot-devicetwin/web"
)

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) > 0 {
		l, err := log.ParseLevel(logLevel)
		if err != nil {
			log.SetLevel(log.TraceLevel)
			log.Tracef("LOG_LEVEL environment variable is set to %s, could not parse to a valid log level. Using trace logging.", logLevel)
		} else {
			log.SetLevel(l)
			log.Infof("Using LOG_LEVEL %s", logLevel)
		}
	}
	// Set up the dependency chain
	settings := config.ParseArgs()
	db, err := factory.CreateDataStore(settings)
	if err != nil {
		log.Fatalf("Error connecting to data store: %v", err)
	}
	m, err := mqtt.GetConnection(settings)
	if err != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", err)
	}
	twin := devicetwin.NewService(settings, db)
	ctrl := controller.NewService(settings, m, twin)

	// Start the web API service
	w := web.NewService(settings, ctrl)
	log.Fatal(w.Run())
}
