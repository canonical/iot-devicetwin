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
)

// createDeviceSnapTable creates the database table and index for snaps
func (db *DataStore) createDeviceSnapTable() error {
	_, err := db.Exec(createDeviceSnapTableSQL)
	if err != nil {
		return err
	}
	_, err = db.Exec(createDeviceSnapIndexSQL)
	return err
}

// DeviceSnapUpsert creates or updates a device snap record
func (db *DataStore) DeviceSnapUpsert(ds datastore.DeviceSnap) error {
	var id int64
	err := db.QueryRow(upsertDeviceSnapSQL, ds.DeviceID, ds.Name, ds.InstalledSize, ds.InstalledDate, ds.Status, ds.Channel, ds.Confinement, ds.Version, ds.Revision, ds.Devmode, ds.Config).Scan(&id)
	if err != nil {
		log.Printf("Error creating device snap %s: %v\n", ds.Name, err)
	}

	return err
}

// DeviceSnapList lists the snaps for a device
func (db *DataStore) DeviceSnapList(deviceID int64) ([]datastore.DeviceSnap, error) {
	rows, err := db.Query(listDeviceSnapSQL, deviceID)
	if err != nil {
		log.Printf("Error retrieving device snaps: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	snaps := []datastore.DeviceSnap{}
	for rows.Next() {
		item := datastore.DeviceSnap{}
		err := rows.Scan(&item.ID, &item.Created, &item.Modified, &item.DeviceID, &item.Name, &item.InstalledSize, &item.InstalledDate, &item.Status, &item.Channel, &item.Confinement, &item.Version, &item.Revision, &item.Devmode, &item.Config)
		if err != nil {
			return nil, err
		}
		snaps = append(snaps, item)
	}

	return snaps, nil
}

// DeviceSnapDelete removes a snap for a device
func (db *DataStore) DeviceSnapDelete(id int64) error {
	_, err := db.Exec(deleteDeviceSnapSQL, id)
	if err != nil {
		log.Printf("Error updating the device snap: %v\n", err)
	}

	return err
}
