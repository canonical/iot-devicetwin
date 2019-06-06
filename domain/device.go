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

package domain

import "time"

// Health update contains enough details to record a device
type Health struct {
	OrganizationID string    `json:"orgId"`
	DeviceID       string    `json:"deviceId"`
	Refresh        time.Time `json:"refresh"`
}

// Device holds the details of a device
type Device struct {
	OrganizationID string        `json:"orgId"`
	DeviceID       string        `json:"deviceId"`
	Brand          string        `json:"brand"`
	Model          string        `json:"model"`
	SerialNumber   string        `json:"serial"`
	StoreID        string        `json:"store"`
	DeviceKey      string        `json:"deviceKey"`
	Version        DeviceVersion `json:"version"`
	Created        time.Time     `json:"created"`
	LastRefresh    time.Time     `json:"lastRefresh"`
}
