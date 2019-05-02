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

import "github.com/CanonicalLtd/iot-devicetwin/domain"

// GroupCreate creates a device group
func (srv *Service) GroupCreate(orgID, name string) error {
	_, err := srv.DB.GroupCreate(orgID, name)
	return err
}

// GroupList lists groups for an organization
func (srv *Service) GroupList(orgID string) ([]domain.Group, error) {
	gg, err := srv.DB.GroupList(orgID)
	if err != nil {
		return nil, err
	}

	groups := []domain.Group{}
	for _, g := range gg {
		groups = append(groups, domain.Group{
			OrganizationID: g.OrganisationID,
			Name:           g.Name,
		})
	}
	return groups, nil
}

// GroupGet retrieves a device group
func (srv *Service) GroupGet(orgID, name string) (domain.Group, error) {
	g, err := srv.DB.GroupGet(orgID, name)
	if err != nil {
		return domain.Group{}, err
	}

	return domain.Group{
		OrganizationID: g.OrganisationID,
		Name:           g.Name,
	}, nil
}

// GroupLinkDevice links a device to a group
func (srv *Service) GroupLinkDevice(orgID, name, clientID string) error {
	return srv.DB.GroupLinkDevice(orgID, name, clientID)
}

// GroupUnlinkDevice unlinks a device from a group
func (srv *Service) GroupUnlinkDevice(orgID, name, clientID string) error {
	return srv.DB.GroupUnlinkDevice(orgID, name, clientID)
}

// GroupGetDevices retrieves the devices from a group
func (srv *Service) GroupGetDevices(orgID, name string) ([]domain.Device, error) {
	dd, err := srv.DB.GroupGetDevices(orgID, name)
	if err != nil {
		return nil, err
	}

	devices := []domain.Device{}
	for _, d := range dd {
		devices = append(devices, dataToDomainDevice(d))
	}
	return devices, nil
}
