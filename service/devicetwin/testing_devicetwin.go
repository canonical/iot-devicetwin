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
	"fmt"
	"github.com/CanonicalLtd/iot-devicetwin/domain"
)

// MockDeviceTwin mocks a device twin service
type MockDeviceTwin struct {
	Actions []string
}

// HealthHandler mocks the health handler
func (twin *MockDeviceTwin) HealthHandler(payload domain.Health) error {
	if payload.DeviceID == "invalid" || payload.DeviceID == "new-device" {
		return fmt.Errorf("MOCK error in health handler")
	}
	return nil
}

// ActionResponse mocks the action handler
func (twin *MockDeviceTwin) ActionResponse(clientID, actionID, action string, payload []byte) error {
	if action == "invalid" {
		return fmt.Errorf("MOCK error in action")
	}
	return nil
}

// DeviceSnaps mocks the snap list
func (twin *MockDeviceTwin) DeviceSnaps(orgID, clientID string) ([]domain.DeviceSnap, error) {
	if clientID == "invalid" {
		return nil, fmt.Errorf("MOCK snaps list")
	}
	return []domain.DeviceSnap{
		{Name: "example-snap", InstalledSize: 2000, Status: "active"},
	}, nil
}

// ActionCreate mocks the action log creation
func (twin *MockDeviceTwin) ActionCreate(orgID, deviceID string, act domain.SubscribeAction) error {
	if deviceID == "invalid" {
		return fmt.Errorf("MOCK action log create")
	}
	if twin.Actions == nil {
		twin.Actions = []string{}
	}
	twin.Actions = append(twin.Actions, act.ID)
	return nil
}

// ActionUpdate mocks the action log update
func (twin *MockDeviceTwin) ActionUpdate(actionID, status, message string) error {
	return nil
}

// DeviceGet mocks fetching a device
func (twin *MockDeviceTwin) DeviceGet(orgID, clientID string) (domain.Device, error) {
	if clientID == "invalid" {
		return domain.Device{}, fmt.Errorf("MOCK error device get")
	}
	return domain.Device{
		OrganizationID: "abc",
		DeviceID:       "c333",
		Brand:          "canonical",
		Model:          "ubuntu-core-18-amd64",
		SerialNumber:   "d75f7300-abbf-4c11-bf0a-8b7103038490",
		DeviceKey:      "CCCCCCCCC",
	}, nil
}

// DeviceList mocks fetching devices for an organization
func (twin *MockDeviceTwin) DeviceList(orgID string) ([]domain.Device, error) {
	if orgID == "invalid" {
		return nil, fmt.Errorf("MOCK error device list")
	}

	return []domain.Device{
		{OrganizationID: "abc",
			DeviceID:     "c333",
			Brand:        "canonical",
			Model:        "ubuntu-core-18-amd64",
			SerialNumber: "d75f7300-abbf-4c11-bf0a-8b7103038490",
			DeviceKey:    "CCCCCCCCC",
		},
	}, nil
}

// GroupCreate mocks creating a group
func (twin *MockDeviceTwin) GroupCreate(orgID, name string) error {
	if orgID == "invalid" {
		return fmt.Errorf("MOCK error group create")
	}
	return nil
}

// GroupList mocks listing groups
func (twin *MockDeviceTwin) GroupList(orgID string) ([]domain.Group, error) {
	if orgID == "invalid" {
		return nil, fmt.Errorf("MOCK error group list")
	}
	return []domain.Group{
		{OrganizationID: "abc", Name: "workshop"},
	}, nil
}

// GroupGet mocks fetching a group
func (twin *MockDeviceTwin) GroupGet(orgID, name string) (domain.Group, error) {
	if orgID == "invalid" || name == "invalid" {
		return domain.Group{}, fmt.Errorf("MOCK error group device unlink")
	}
	return domain.Group{
		OrganizationID: "abc", Name: "workshop",
	}, nil
}

// GroupLinkDevice mocks linking a device to a group
func (twin *MockDeviceTwin) GroupLinkDevice(orgID, name, clientID string) error {
	if orgID == "invalid" || name == "invalid" || clientID == "invalid" {
		return fmt.Errorf("MOCK error group device link")
	}
	return nil
}

// GroupUnlinkDevice mocks unlinking a device from a group
func (twin *MockDeviceTwin) GroupUnlinkDevice(orgID, name, clientID string) error {
	if orgID == "invalid" || name == "invalid" || clientID == "invalid" {
		return fmt.Errorf("MOCK error group device unlink")
	}
	return nil
}

// GroupGetDevices mocks retrieving the devices for a group
func (twin *MockDeviceTwin) GroupGetDevices(orgID, name string) ([]domain.Device, error) {
	if orgID == "invalid" || name == "invalid" {
		return nil, fmt.Errorf("MOCK error group device unlink")
	}
	return []domain.Device{
		{OrganizationID: "abc",
			DeviceID:     "c333",
			Brand:        "canonical",
			Model:        "ubuntu-core-18-amd64",
			SerialNumber: "d75f7300-abbf-4c11-bf0a-8b7103038490",
			DeviceKey:    "CCCCCCCCC",
		},
	}, nil
}
