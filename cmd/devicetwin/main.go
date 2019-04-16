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
	"github.com/CanonicalLtd/iot-devicetwin/config"
	"github.com/CanonicalLtd/iot-devicetwin/datastore/memory"
	"github.com/CanonicalLtd/iot-devicetwin/service/controller"
	"github.com/CanonicalLtd/iot-devicetwin/service/devicetwin"
	"github.com/CanonicalLtd/iot-devicetwin/service/mqtt"
	"github.com/CanonicalLtd/iot-devicetwin/web"
	"log"
)

func main() {
	// Set up the dependency chain
	settings := config.ParseArgs()
	db := memory.NewStore()
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
