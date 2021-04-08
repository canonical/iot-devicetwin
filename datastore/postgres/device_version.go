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

package postgres

import (
	"github.com/canonical/iot-devicetwin/datastore"
	"log"
)

// createDeviceVersionTable creates the database table for a device's OS details
func (db *DataStore) createDeviceVersionTable() error {
	_, err := db.Exec(createDeviceVersionTableSQL)
	return err
}

// DeviceVersionGet fetches a device version details from the database
func (db *DataStore) DeviceVersionGet(deviceID int64) (datastore.DeviceVersion, error) {
	item := datastore.DeviceVersion{}
	row := db.QueryRow(getDeviceVersionSQL, deviceID)
	err := row.Scan(&item.ID, &item.DeviceID, &item.Version, &item.Series, &item.OSID, &item.OSVersionID, &item.OnClassic, &item.KernelVersion)
	if err != nil {
		log.Printf("Error retrieving device version: %v\n", err)
	}
	return item, err
}

// DeviceVersionUpsert creates or updates a device version record
func (db *DataStore) DeviceVersionUpsert(dv datastore.DeviceVersion) error {
	var id int64
	err := db.QueryRow(upsertDeviceVersionSQL, dv.DeviceID, dv.Version, dv.Series, dv.OSID, dv.OSVersionID, dv.OnClassic, dv.KernelVersion).Scan(&id)
	if err != nil {
		log.Printf("Error creating device version: %v\n", err)
	}

	return err
}

// DeviceVersionDelete removes a device version
func (db *DataStore) DeviceVersionDelete(id int64) error {
	_, err := db.Exec(deleteDeviceVersionSQL, id)
	if err != nil {
		log.Printf("Error updating the device snap: %v\n", err)
	}

	return err
}
