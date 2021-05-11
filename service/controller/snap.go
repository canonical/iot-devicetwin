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

// Package controller is the interface to handling with REST API actions
package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/everactive/iot-devicetwin/pkg/actions"
	"github.com/everactive/iot-devicetwin/pkg/messages"
)

const (
	snapActionDelay = 10 * time.Second
)

// DeviceSnaps gets the device's snaps from the database cache
func (srv *Service) DeviceSnaps(orgID, clientID string) ([]messages.DeviceSnap, error) {
	return srv.DeviceTwin.DeviceSnaps(orgID, clientID)
}

// DeviceSnapList triggers listing snaps on a device
func (srv *Service) DeviceSnapList(orgID, clientID string) error {
	act := messages.SubscribeAction{
		Action: actions.List,
	}
	return srv.deviceSnapAction(orgID, clientID, act)
}

// DeviceSnapServiceAction triggers stop,`start, or restart for a snap on a device
func (srv *Service) DeviceSnapServiceAction(orgID, clientID, snap, action string, services *messages.SnapServiceMessage) error {
	switch action {
	case actions.Start, actions.Stop, actions.Restart:
	default:
		return fmt.Errorf("invalid snap service action `%s`", action)
	}

	jsonBytes, err := json.Marshal(services)
	if err != nil {
		return err
	}

	act := messages.SubscribeAction{
		Action: action,
		Snap:   snap,
		Data:   string(jsonBytes),
	}
	return srv.deviceSnapAction(orgID, clientID, act)
}

// DeviceSnapInstall triggers installing a snap on a device
func (srv *Service) DeviceSnapInstall(orgID, clientID, snap string) error {
	act := messages.SubscribeAction{
		Action: actions.Install,
		Snap:   snap,
	}
	return srv.deviceSnapAction(orgID, clientID, act)
}

// DeviceSnapRemove triggers uninstalling a snap on a device
func (srv *Service) DeviceSnapRemove(orgID, clientID, snap string) error {
	act := messages.SubscribeAction{
		Action: actions.Remove,
		Snap:   snap,
	}
	return srv.deviceSnapAction(orgID, clientID, act)
}

// DeviceSnapUpdate triggers a snap update on a device
func (srv *Service) DeviceSnapUpdate(orgID, clientID, snap, action string) error {
	switch action {
	case actions.Enable, actions.Disable, actions.Refresh:
		act := messages.SubscribeAction{
			Action: action,
			Snap:   snap,
		}
		return srv.deviceSnapAction(orgID, clientID, act)
	default:
		return fmt.Errorf("invalid update action `%s`", action)
	}
}

// DeviceSnapConf triggers a snap settings update on a device
func (srv *Service) DeviceSnapConf(orgID, clientID, snap, settings string) error {
	// Trigger the update settings action on the device
	act := messages.SubscribeAction{
		Action: actions.SetConf,
		Snap:   snap,
		Data:   settings,
	}
	return srv.deviceSnapAction(orgID, clientID, act)
}

// deviceSnapAction triggers a snap action on a device
func (srv *Service) deviceSnapAction(orgID, clientID string, action messages.SubscribeAction) error {
	// Validate the org and device ID
	device, err := srv.DeviceTwin.DeviceGet(orgID, clientID)
	if err != nil {
		return err
	}

	// Trigger the action on the device
	err = srv.triggerActionOnDevice(device.OrgId, device.DeviceId, action)
	if err != nil {
		return err
	}

	// State of the snaps has changed, so request a snap list
	if action.Action != actions.List {
		// Request the list action after a few seconds
		time.AfterFunc(snapActionDelay, func() {
			_ = srv.DeviceSnapList(orgID, clientID)
		})
	}
	return err
}
