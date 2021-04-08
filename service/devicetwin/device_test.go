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
	"testing"

	"github.com/canonical/iot-devicetwin/config"
	"github.com/canonical/iot-devicetwin/datastore/memory"
)

func TestService_DeviceGet(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111"}, false},
		{"valid-with-server", args{"abc", "c333"}, false},
		{"invalid", args{"abc", "invalid"}, true},
		{"invalid-orgid", args{"invalid", "a111"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			got, err := srv.DeviceGet(tt.args.orgID, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.DeviceID != tt.args.clientID {
				t.Errorf("Service.DeviceGet() = %v, want %v", got.DeviceID, tt.args.clientID)
			}
		})
	}
}

func TestService_DeviceList(t *testing.T) {
	type args struct {
		orgID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc"}, 3, false},
		{"valid-no-devices", args{"none"}, 0, false},
		{"invalid", args{"invalid"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			got, err := srv.DeviceList(tt.args.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Service.DeviceList() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
