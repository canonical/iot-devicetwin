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

package devicetwin

import (
	"encoding/json"
	"fmt"
	"github.com/CanonicalLtd/iot-devicetwin/config"
	"github.com/CanonicalLtd/iot-devicetwin/datastore"
	"github.com/CanonicalLtd/iot-devicetwin/domain"
	"log"
)

// DeviceTwin interface for the service
type DeviceTwin interface {
	HealthHandler(payload domain.Health) error
	ActionResponse(clientID, action string, payload []byte) error // process a response from a device

	DeviceSnaps(clientID string) ([]domain.DeviceSnap, error)
}

// Service implementation of the identity use cases
type Service struct {
	Settings *config.Settings
	DB       datastore.DataStore
}

// NewService creates an implementation of the device twin use cases
func NewService(settings *config.Settings, db datastore.DataStore) *Service {
	return &Service{
		Settings: settings,
		DB:       db,
	}
}

// HealthHandler handles a health update from a device
func (srv *Service) HealthHandler(payload domain.Health) error {
	// Check that we have the device
	_, err := srv.DB.DeviceGet(payload.DeviceID)
	if err != nil {
		// Request the device details to be published as we don't have it
		return err
	}

	// Update the last refresh on the device
	return srv.DB.DevicePing(payload.DeviceID, payload.Refresh)
}

// ActionResponse handles action response from a device
func (srv *Service) ActionResponse(clientID, action string, payload []byte) error {
	// Act based on the message action
	switch action {
	case "device":
		return srv.actionDevice(payload)
	case "list":
		return srv.actionList(clientID, payload)
	//case "install":
	//case "remove":
	//case "refresh":
	//case "revert":
	//case "enable":
	//case "disable":
	//case "conf":
	//case "setconf":
	//case "info":
	//case "ack":
	//case "server":
	default:
		return fmt.Errorf("error unhandled action `%s`", action)
	}
}

// actionDevice process the device info received from a device
func (srv *Service) actionDevice(payload []byte) error {
	// Parse the payload
	d := domain.PublishDevice{}
	if err := json.Unmarshal(payload, &d); err != nil {
		log.Printf("Error in device action message: %v", err)
		return fmt.Errorf("error in device action message: %v", err)
	}

	// Get the device details and create/update the device
	_, err := srv.DB.DeviceGet(d.Result.DeviceID)
	if err == nil {
		return fmt.Errorf("error in device action: device already exists")
	}

	// Device does not exit, so create
	device := datastore.Device{
		OrganisationID: d.Result.OrganizationID,
		DeviceID:       d.Result.DeviceID,
		Brand:          d.Result.Brand,
		Model:          d.Result.Model,
		SerialNumber:   d.Result.SerialNumber,
		DeviceKey:      d.Result.DeviceKey,
		StoreID:        d.Result.StoreID,
	}
	_, err = srv.DB.DeviceCreate(device)
	return err
}

// actionList process the list of snaps received from a device
func (srv *Service) actionList(clientID string, payload []byte) error {
	// Parse the payload
	p := domain.PublishSnaps{}
	if err := json.Unmarshal(payload, &p); err != nil {
		log.Printf("Error in list action message: %v", err)
		return fmt.Errorf("error in list action message: %v", err)
	}

	// Get the device details
	device, err := srv.DB.DeviceGet(clientID)
	if err != nil {
		return fmt.Errorf("cannot find device with ID `%s`", clientID)
	}

	// Delete the existing snap list for the device
	if err := srv.DB.DeviceSnapDelete(device.ID); err != nil {
		return fmt.Errorf("error deleting snap records: %v", err)
	}

	// Add the installed snaps
	for _, s := range p.Result {
		snap := datastore.DeviceSnap{
			//Created       time.Time
			//Modified      time.Time
			DeviceID:      device.ID,
			Name:          s.Name,
			InstalledSize: s.InstalledSize,
			InstalledDate: s.InstalledDate,
			Status:        s.Status,
			Channel:       s.Channel,
			Confinement:   s.Confinement,
			Version:       s.Version,
			Revision:      s.Revision,
			Devmode:       s.Devmode,
			Config:        s.Config,
		}

		if err := srv.DB.DeviceSnapUpsert(snap); err != nil {
			return err
		}
	}

	return nil
}
