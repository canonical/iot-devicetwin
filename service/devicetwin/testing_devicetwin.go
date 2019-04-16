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
type MockDeviceTwin struct{}

// HealthHandler mocks the health handler
func (twin *MockDeviceTwin) HealthHandler(payload domain.Health) error {
	if payload.DeviceID == "invalid" {
		return fmt.Errorf("MOCK error in health handler")
	}
	return nil
}

// ActionResponse mocks the action handler
func (twin *MockDeviceTwin) ActionResponse(action string, payload []byte) error {
	if action == "invalid" {
		return fmt.Errorf("MOCK error in action")
	}
	return nil
}
