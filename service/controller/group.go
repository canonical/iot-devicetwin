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

package controller

import "github.com/canonical/iot-devicetwin/domain"

// GroupCreate creates a device group
func (srv *Service) GroupCreate(orgID, name string) error {
	return srv.DeviceTwin.GroupCreate(orgID, name)
}

// GroupList lists the groups for an organization
func (srv *Service) GroupList(orgID string) ([]domain.Group, error) {
	return srv.DeviceTwin.GroupList(orgID)
}

// GroupGet retrieves a device group
func (srv *Service) GroupGet(orgID, name string) (domain.Group, error) {
	return srv.DeviceTwin.GroupGet(orgID, name)
}

// GroupLinkDevice links a device to a group
func (srv *Service) GroupLinkDevice(orgID, name, clientID string) error {
	return srv.DeviceTwin.GroupLinkDevice(orgID, name, clientID)
}

// GroupUnlinkDevice unlinks a device from a group
func (srv *Service) GroupUnlinkDevice(orgID, name, clientID string) error {
	return srv.DeviceTwin.GroupUnlinkDevice(orgID, name, clientID)
}

// GroupGetDevices retrieves the devices from a group
func (srv *Service) GroupGetDevices(orgID, name string) ([]domain.Device, error) {
	return srv.DeviceTwin.GroupGetDevices(orgID, name)
}

// GroupGetExcludedDevices retrieves the devices not in a group
func (srv *Service) GroupGetExcludedDevices(orgID, name string) ([]domain.Device, error) {
	return srv.DeviceTwin.GroupGetExcludedDevices(orgID, name)
}
