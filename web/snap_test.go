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

package web

import (
	"testing"

	"github.com/CanonicalLtd/iot-devicetwin/config"
	"github.com/CanonicalLtd/iot-devicetwin/service/controller"
	"github.com/CanonicalLtd/iot-devicetwin/service/devicetwin"
	"github.com/CanonicalLtd/iot-devicetwin/service/mqtt"
)

func testController() controller.Controller {
	return controller.NewService(config.TestConfig(), &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
}

func TestService_SnapList(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		code   int
		result string
	}{
		{"valid", "/v1/device/abc/a111/snaps", 200, ""},
		{"invalid", "/v1/device/abc/invalid/snaps", 400, "SnapList"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())

			w := sendRequest("GET", tt.url, nil, wb)
			if w.Code != tt.code {
				t.Errorf("Web.SnapList() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseSnapsResponse(w.Body)
			if err != nil {
				t.Errorf("Web.SnapList() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.SnapList() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_SnapInstall(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		code   int
		result string
	}{
		{"valid", "/v1/device/abc/a111/snaps/helloworld", 200, ""},
		{"invalid", "/v1/device/abc/invalid/snaps/helloworld", 400, "SnapInstall"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest("POST", tt.url, nil, wb)
			if w.Code != tt.code {
				t.Errorf("Web.SnapInstall() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.SnapInstall() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.SnapInstall() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}
