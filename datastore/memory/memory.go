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

package memory

import (
	"fmt"
	"github.com/CanonicalLtd/iot-devicetwin/datastore"
	"sync"
	"time"
)

// Store implements an in-memory store for testing
type Store struct {
	Devices []datastore.Device
	lock    sync.RWMutex
}

// NewStore creates a new memory store
func NewStore() *Store {
	d1 := datastore.Device{ID: 1, OrganisationID: "abc", DeviceID: "a111", Brand: "example", Model: "drone-1000", SerialNumber: "DR1000A111", DeviceKey: "AAAAAAAAA", StoreID: "example-store", Active: true}
	d2 := datastore.Device{ID: 2, OrganisationID: "abc", DeviceID: "b222", Brand: "example", Model: "drone-1000", SerialNumber: "DR1000B222", DeviceKey: "BBBBBBBBB", StoreID: "example-store", Active: true}
	d3 := datastore.Device{ID: 3, OrganisationID: "abc", DeviceID: "c333", Brand: "canonical", Model: "ubuntu-core-18-amd64", SerialNumber: "d75f7300-abbf-4c11-bf0a-8b7103038490", DeviceKey: "CCCCCCCCC", Active: true}

	return &Store{
		Devices: []datastore.Device{d1, d2, d3},
	}
}

// DeviceGet fetches an existing device
func (mem *Store) DeviceGet(id string) (datastore.Device, error) {
	mem.lock.RLock()
	defer mem.lock.RUnlock()

	for _, d := range mem.Devices {
		if d.DeviceID == id {
			return d, nil
		}
	}
	return datastore.Device{}, fmt.Errorf("device with ID `%s` not found", id)
}

// DevicePing updates a device to indicate its health
func (mem *Store) DevicePing(id string, refresh time.Time) error {
	device, err := mem.DeviceGet(id)
	if err != nil {
		return err
	}

	mem.lock.Lock()
	defer mem.lock.Unlock()
	device.LastRefresh = refresh

	for i := range mem.Devices {
		if mem.Devices[i].DeviceID == id {
			mem.Devices[i] = device
		}
	}
	return nil
}

// DeviceCreate creates a new device
func (mem *Store) DeviceCreate(device datastore.Device) (int64, error) {
	// Check the device does not exist
	if _, err := mem.DeviceGet(device.DeviceID); err == nil {
		return 0, fmt.Errorf("device with ID `%s` already exists", device.DeviceID)
	}

	mem.lock.Lock()
	defer mem.lock.Unlock()

	device.Created = time.Now()
	device.LastRefresh = time.Now()

	device.ID = int64(len(mem.Devices) + 1)
	mem.Devices = append(mem.Devices, device)
	return device.ID, nil
}
