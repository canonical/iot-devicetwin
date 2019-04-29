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

import (
	"fmt"
	"github.com/CanonicalLtd/iot-devicetwin/domain"
)

// DeviceSnaps gets the device's snaps from the database cache
func (srv *Service) DeviceSnaps(orgID, clientID string) ([]domain.DeviceSnap, error) {
	return srv.DeviceTwin.DeviceSnaps(orgID, clientID)
}

// DeviceSnapInstall triggers installing a snap on a device
func (srv *Service) DeviceSnapInstall(orgID, clientID, snap string) error {
	act := domain.SubscribeAction{
		Action: "install",
		Snap:   snap,
	}
	return srv.deviceSnapAction(orgID, clientID, snap, act)
}

// DeviceSnapRemove triggers uninstalling a snap on a device
func (srv *Service) DeviceSnapRemove(orgID, clientID, snap string) error {
	act := domain.SubscribeAction{
		Action: "remove",
		Snap:   snap,
	}
	return srv.deviceSnapAction(orgID, clientID, snap, act)
}

// DeviceSnapUpdate triggers a snap update on a device
func (srv *Service) DeviceSnapUpdate(orgID, clientID, snap, action string) error {
	switch action {
	case "enable", "disable", "refresh":
		act := domain.SubscribeAction{
			Action: action,
			Snap:   snap,
		}
		return srv.deviceSnapAction(orgID, clientID, snap, act)
	default:
		return fmt.Errorf("invalid update action `%s`", action)
	}
}

// DeviceSnapConf triggers a snap settings update on a device
func (srv *Service) DeviceSnapConf(orgID, clientID, snap, settings string) error {
	// Trigger the update settings action on the device
	act := domain.SubscribeAction{
		Action: "setconf",
		Snap:   snap,
		Data:   settings,
	}
	return srv.deviceSnapAction(orgID, clientID, snap, act)
}

// deviceSnapAction triggers a snap action on a device
func (srv *Service) deviceSnapAction(orgID, clientID, snap string, action domain.SubscribeAction) error {
	// Validate the org and device ID
	device, err := srv.DeviceTwin.DeviceGet(orgID, clientID)
	if err != nil {
		return err
	}

	// Trigger the install action on the device
	return srv.triggerActionOnDevice(device.OrganizationID, device.DeviceID, action)
}
