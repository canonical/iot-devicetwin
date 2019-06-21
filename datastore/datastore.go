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

package datastore

import (
	"time"
)

// DataStore is the interfaces for the data repository
type DataStore interface {
	DeviceList(orgID string) ([]Device, error)
	DeviceGet(id string) (Device, error)
	DevicePing(id string, refresh time.Time) error
	DeviceCreate(Device) (int64, error)

	DeviceSnapList(id int64) ([]DeviceSnap, error)
	DeviceSnapDelete(id int64) error
	DeviceSnapUpsert(ds DeviceSnap) error

	ActionCreate(act Action) (int64, error)
	ActionUpdate(actionID, status, message string) error
	ActionListForDevice(deviceID string) ([]Action, error)

	DeviceVersionGet(deviceID int64) (DeviceVersion, error)
	DeviceVersionUpsert(dv DeviceVersion) error
	DeviceVersionDelete(id int64) error

	GroupCreate(orgID, name string) (int64, error)
	GroupList(orgID string) ([]Group, error)
	GroupGet(orgID, name string) (Group, error)
	GroupLinkDevice(orgID, name, deviceID string) error
	GroupUnlinkDevice(orgID, name, deviceID string) error
	GroupGetDevices(orgID, name string) ([]Device, error)
	GroupGetExcludedDevices(orgID, name string) ([]Device, error)
}
