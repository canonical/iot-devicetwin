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

import "time"

// Action is the log of an action request
type Action struct {
	ID             int64
	Created        time.Time
	Modified       time.Time
	OrganizationID string
	DeviceID       string
	ActionID       string
	Action         string
	Status         string
	Message        string
}

// Device the repository definition of a device
type Device struct {
	ID             int64
	Created        time.Time
	LastRefresh    time.Time
	OrganisationID string
	DeviceID       string
	Brand          string
	Model          string
	SerialNumber   string
	DeviceKey      string
	StoreID        string
	Active         bool
}

// DeviceSnap holds the details of snap on a device
type DeviceSnap struct {
	ID            int
	Created       time.Time
	Modified      time.Time
	DeviceID      int64
	Name          string
	InstalledSize int64
	InstalledDate time.Time
	Status        string
	Channel       string
	Confinement   string
	Version       string
	Revision      int
	Devmode       bool
	Config        string
}

// DeviceVersion holds the details of the OS details on the device
type DeviceVersion struct {
	ID            int64
	DeviceID      int64
	Version       string
	Series        string
	OSID          string
	OSVersionID   string
	OnClassic     bool
	KernelVersion string
}

// Group is the record for grouping devices
type Group struct {
	ID             int64
	Created        time.Time
	Modified       time.Time
	OrganisationID string
	Name           string
}

// GroupDeviceLink is the record for linking devices to groups
type GroupDeviceLink struct {
	ID             int64
	Created        time.Time
	OrganisationID string
	GroupID        int64
	DeviceID       int64
}
