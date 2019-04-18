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
	"github.com/CanonicalLtd/iot-devicetwin/datastore"
	"github.com/CanonicalLtd/iot-devicetwin/domain"
	"log"
)

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

// actionInstall process the install snap response
func (srv *Service) actionInstall(clientID string, payload []byte) error {
	return fmt.Errorf("not implemented")
}
