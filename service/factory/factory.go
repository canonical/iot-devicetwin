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

package factory

import (
	"fmt"
	"github.com/canonical/iot-devicetwin/config"
	"github.com/canonical/iot-devicetwin/datastore"
	"github.com/canonical/iot-devicetwin/datastore/memory"
	"github.com/canonical/iot-devicetwin/datastore/postgres"
)

// CreateDataStore is the factory method to create a data store
func CreateDataStore(settings *config.Settings) (datastore.DataStore, error) {
	var db datastore.DataStore
	switch settings.Driver {
	case "memory":
		db = memory.NewStore()
	case "postgres":
		db = postgres.OpenDataStore(settings.Driver, settings.DataSource)
	default:
		return nil, fmt.Errorf("unknown data store driver: %v", settings.Driver)
	}

	return db, nil
}
