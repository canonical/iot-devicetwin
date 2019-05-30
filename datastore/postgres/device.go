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
	"github.com/CanonicalLtd/iot-devicetwin/datastore"
	"log"
	"time"
)

// CreateDeviceTable creates the database table for devices with its indexes.
func (db *DataStore) createDeviceTable() error {
	_, err := db.Exec(createDeviceTableSQL)
	return err
}

// DeviceCreate adds a new record to device database table, returning the record ID
func (db *DataStore) DeviceCreate(device datastore.Device) (int64, error) {
	var id int64
	err := db.QueryRow(createDeviceSQL, device.OrganisationID, device.DeviceID, device.Brand, device.Model, device.SerialNumber, device.StoreID, device.DeviceKey).Scan(&id)
	if err != nil {
		log.Printf("Error creating device %s/%s: %v\n", device.Brand, device.Model, err)
	}

	return id, err
}

// DeviceGet fetches a device from the database
func (db *DataStore) DeviceGet(deviceID string) (datastore.Device, error) {
	item := datastore.Device{}
	row := db.QueryRow(getDeviceSQL, deviceID)
	err := row.Scan(&item.ID, &item.Created, &item.LastRefresh, &item.OrganisationID, &item.DeviceID, &item.Brand, &item.Model, &item.SerialNumber, &item.StoreID, &item.DeviceKey, &item.Active)
	if err != nil {
		log.Printf("Error retrieving device %s: %v\n", deviceID, err)
	}
	return item, err
}

// DevicePing updates the last ping time from a device
func (db *DataStore) DevicePing(deviceID string, refresh time.Time) error {
	_, err := db.Exec(pingDeviceSQL, deviceID, refresh)
	if err != nil {
		log.Printf("Error updating the device: %v\n", err)
	}

	return err
}

// DeviceList fetches the devices for an organization from the database
func (db *DataStore) DeviceList(orgID string) ([]datastore.Device, error) {
	rows, err := db.Query(listDeviceSQL, orgID)
	if err != nil {
		log.Printf("Error retrieving devices: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	devices := []datastore.Device{}
	for rows.Next() {
		item := datastore.Device{}
		err := rows.Scan(&item.ID, &item.Created, &item.LastRefresh, &item.OrganisationID, &item.DeviceID, &item.Brand, &item.Model, &item.SerialNumber, &item.StoreID, &item.DeviceKey, &item.Active)
		if err != nil {
			return nil, err
		}
		devices = append(devices, item)
	}

	return devices, nil
}
