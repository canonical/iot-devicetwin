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
	"github.com/CanonicalLtd/iot-devicetwin/datastore/memory"
	"testing"

	"github.com/CanonicalLtd/iot-devicetwin/config"
	"github.com/CanonicalLtd/iot-devicetwin/domain"
)

func TestService_HealthHandler(t *testing.T) {
	h1 := domain.Health{OrganizationID: "abc", DeviceID: "a111"}
	h2 := domain.Health{OrganizationID: "abc", DeviceID: "invalid"}
	type args struct {
		payload domain.Health
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{h1}, false},
		{"invalid", args{h2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.HealthHandler(tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Service.HealthHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ActionResponse(t *testing.T) {
	p1 := []byte("")
	p2 := []byte(`{"id":"a1", "action":"device", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"d444", "brand":"example", "model":"drone-1000", "serial":"d444"}}`)
	p3 := []byte(`{"id":"a1", "action":"device", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"a111", "brand":"example", "model":"drone-1000", "serial":"d444"}}`)
	p4 := []byte(`{"id":"a1", "action":"list", "success":true, "message":"", "result": [{"name":"abc", "status":"active", "version":"1.0"}, {"name":"alpaca", "status":"active", "version":"2.3"}]}`)

	type args struct {
		clientID string
		action   string
		payload  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty-payload", args{"", "device", p1}, true},
		{"valid-device", args{"d444", "device", p2}, false},
		{"invalid-action", args{"d444", "invalid", p2}, true},
		{"device-exists", args{"a111", "device", p3}, true},

		{"valid-list", args{"a111", "list", p4}, false},
		{"list-empty-payload", args{"", "list", p1}, true},
		{"list-no-device", args{"invalid", "list", p4}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.ActionResponse(tt.args.clientID, "a1", tt.args.action, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("Service.ActionResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
