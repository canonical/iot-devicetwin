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
	"testing"

	"github.com/CanonicalLtd/iot-devicetwin/service/devicetwin"
	"github.com/CanonicalLtd/iot-devicetwin/service/mqtt"
)

func TestService_GroupCreate(t *testing.T) {
	type args struct {
		orgID string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "test-group"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.GroupCreate(tt.args.orgID, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GroupList(t *testing.T) {
	type args struct {
		orgID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			got, err := srv.GroupList(tt.args.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Service.GroupList() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestService_GroupGet(t *testing.T) {
	type args struct {
		orgID string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "workshop"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			got, err := srv.GroupGet(tt.args.orgID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Name != tt.args.name {
					t.Errorf("Service.GroupGet() name = %v, wantErr %v", got.Name, tt.args.name)
				}
			}
		})
	}
}

func TestService_GroupLinkDevice(t *testing.T) {
	type args struct {
		orgID    string
		name     string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "workshop", "c333"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.GroupLinkDevice(tt.args.orgID, tt.args.name, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupLinkDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GroupUnlinkDevice(t *testing.T) {
	type args struct {
		orgID    string
		name     string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "workshop", "c333"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.GroupUnlinkDevice(tt.args.orgID, tt.args.name, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupUnlinkDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GroupGetDevices(t *testing.T) {
	type args struct {
		orgID string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc", "workshop"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			got, err := srv.GroupGetDevices(tt.args.orgID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupGetDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Service.GroupGetDevices() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestService_GroupGetExcludedDevices(t *testing.T) {
	type args struct {
		orgID string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc", "workshop"}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			got, err := srv.GroupGetExcludedDevices(tt.args.orgID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupGetExcludedDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Service.GroupGetExcludedDevices() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
