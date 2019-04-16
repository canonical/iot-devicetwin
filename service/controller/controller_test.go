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

	"github.com/CanonicalLtd/iot-devicetwin/config"
	"github.com/CanonicalLtd/iot-devicetwin/service/devicetwin"
	"github.com/CanonicalLtd/iot-devicetwin/service/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var settings = config.TestConfig()

func TestService_SubscribeToActions(t *testing.T) {
	type fields struct {
		Settings   *config.Settings
		MQTT       mqtt.Connect
		DeviceTwin devicetwin.DeviceTwin
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(tt.fields.Settings, tt.fields.MQTT, tt.fields.DeviceTwin)
			if err := srv.SubscribeToActions(); (err != nil) != tt.wantErr {
				t.Errorf("Service.SubscribeToActions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ActionHandler(t *testing.T) {
	m1 := []byte(`{"success": false, "message": "MOCK error"}`)
	m2 := []byte(`{"success": true, "action": "invalid"}`)

	type fields struct {
		Settings   *config.Settings
		MQTT       mqtt.Connect
		DeviceTwin devicetwin.DeviceTwin
	}
	type args struct {
		client MQTT.Client
		msg    MQTT.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"valid", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{}}},
		{"error-response", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m1}}},
		{"invalid-action", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(tt.fields.Settings, tt.fields.MQTT, tt.fields.DeviceTwin)
			srv.ActionHandler(tt.args.client, tt.args.msg)
		})
	}
}

func TestService_HealthHandler(t *testing.T) {
	m1 := []byte(`{"orgId": "abc", "deviceId": "aa111"}`)
	m2 := []byte(`{"orgId": "abc", "deviceId": "invalid"}`)

	type fields struct {
		Settings   *config.Settings
		MQTT       mqtt.Connect
		DeviceTwin devicetwin.DeviceTwin
	}
	type args struct {
		client MQTT.Client
		msg    MQTT.Message
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"valid", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m1}}},
		{"invalid-message", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{}}},
		{"invalid-clientID", fields{settings, &mqtt.MockConnect{}, &devicetwin.MockDeviceTwin{}}, args{&mqtt.MockClient{}, &mqtt.MockMessage{Message: m2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(tt.fields.Settings, tt.fields.MQTT, tt.fields.DeviceTwin)
			srv.HealthHandler(tt.args.client, tt.args.msg)
		})
	}
}

func Test_getClientID(t *testing.T) {
	type args struct {
		msg MQTT.Message
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"valid", args{&mqtt.MockMessage{}}, "aa111"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getClientID(tt.args.msg); got != tt.want {
				t.Errorf("getClientID() = %v, want %v", got, tt.want)
			}
		})
	}
}
