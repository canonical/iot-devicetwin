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

import "github.com/CanonicalLtd/iot-devicetwin/domain"

// DeviceSnaps gets the device's snaps from the database cache
func (srv *Service) DeviceSnaps(orgID, clientID string) ([]domain.DeviceSnap, error) {
	return srv.DeviceTwin.DeviceSnaps(orgID, clientID)
}

// DeviceSnapInstall triggers installing a snap on a device
func (srv *Service) DeviceSnapInstall(orgID, clientID, snap string) error {
	// Validate the org and device ID
	device, err := srv.DeviceTwin.DeviceGet(orgID, clientID)
	if err != nil {
		return err
	}

	// Trigger the install action on the device
	action := domain.SubscribeAction{
		Action: "install",
		Snap:   snap,
	}
	return srv.triggerActionOnDevice(device.OrganizationID, device.DeviceID, action)
}
