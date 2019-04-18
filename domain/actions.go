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

// SubscribeAction is the message format for the action topic
type SubscribeAction struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Snap   string `json:"snap"`
	Data   string `json:"data"`
}

// PublishResponse is the published message showing the result of an action
type PublishResponse struct {
	ID      string      `json:"id"`
	Action  string      `json:"action"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// ActionResponseBytes is the published message showing the result of an action, with the result as bytes
type ActionResponseBytes struct {
	ID      string `json:"id"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  []byte `json:"result"`
}

// PublishDevice is the published message showing the result of a device action
type PublishDevice struct {
	ID      string `json:"id"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		OrganizationID string `json:"orgId"`
		DeviceID       string `json:"deviceId"`
		Brand          string `json:"brand"`
		Model          string `json:"model"`
		SerialNumber   string `json:"serial"`
		StoreID        string `json:"store"`
		DeviceKey      string `json:"deviceKey"`
	} `json:"result"`
}

// PublishSnaps is the published message showing the result of a list snaps action
type PublishSnaps struct {
	ID      string       `json:"id"`
	Action  string       `json:"action"`
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Result  []DeviceSnap `json:"result"`
}

// DeviceSnap holds the details of snap on a device
type DeviceSnap struct {
	DeviceID      string    `json:"deviceId"`
	Name          string    `json:"name"`
	InstalledSize int64     `json:"installedSize"`
	InstalledDate time.Time `json:"installedDate"`
	Status        string    `json:"status"`
	Channel       string    `json:"channel"`
	Confinement   string    `json:"confinement"`
	Version       string    `json:"version"`
	Revision      int       `json:"revision"`
	Devmode       bool      `json:"devmode"`
	Config        string    `json:"config"`
}
