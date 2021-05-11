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

	"github.com/everactive/iot-devicetwin/pkg/messages"
	"github.com/everactive/iot-devicetwin/service/devicetwin"
	"github.com/everactive/iot-devicetwin/service/mqtt"
)

func TestService_DeviceSnaps(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc", "a111"}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			got, err := srv.DeviceSnaps(tt.args.orgID, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnaps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Service.DeviceSnaps() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestService_DeviceSnapInstall(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
		snap     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111", "helloworld"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.DeviceSnapInstall(tt.args.orgID, tt.args.clientID, tt.args.snap); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnapInstall() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DeviceSnapRemove(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
		snap     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111", "helloworld"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.DeviceSnapRemove(tt.args.orgID, tt.args.clientID, tt.args.snap); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnapRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DeviceSnapUpdate(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
		snap     string
		action   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111", "helloworld", "enable"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.DeviceSnapUpdate(tt.args.orgID, tt.args.clientID, tt.args.snap, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnapUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DeviceSnapConf(t *testing.T) {
	type args struct {
		orgID    string
		clientID string
		snap     string
		settings string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111", "helloworld", `{"title": "Hello World!"}`}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.DeviceSnapConf(tt.args.orgID, tt.args.clientID, tt.args.snap, tt.args.settings); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnapConf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DeviceSnapList(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.DeviceSnapList(tt.args.orgID, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnapList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_DeviceSnapServiceAction(t *testing.T) {

	type args struct {
		orgID    string
		clientID string
		snap     string
		action   string
		services *messages.SnapServiceMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid start", args{"abc", "a111", "helloworld", "start", &messages.SnapServiceMessage{Services: []string{"hello_service"}}}, false},
		{"valid start multiple snap services", args{"abc", "a111", "helloworld", "start", &messages.SnapServiceMessage{Services: []string{"hello_service", "goodbye_service"}}}, false},
		{"valid stop", args{"abc", "a111", "helloworld", "stop", &messages.SnapServiceMessage{Services: []string{"hello_service"}}}, false},
		{"valid stop multiple snap services", args{"abc", "a111", "helloworld", "stop", &messages.SnapServiceMessage{Services: []string{"hello_service", "goodbyte service"}}}, false},
		{"valid restart", args{"abc", "a111", "helloworld", "restart", &messages.SnapServiceMessage{Services: []string{"hello_service"}}}, false},
		{"valid restart multiple snap services", args{"abc", "a111", "helloworld", "restart", &messages.SnapServiceMessage{Services: []string{"hello_service", "goodbyte service"}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{})
			if err := srv.DeviceSnapServiceAction(tt.args.orgID, tt.args.clientID, tt.args.snap, tt.args.action, tt.args.services); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeviceSnapServiceAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
