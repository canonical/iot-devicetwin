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

	"github.com/CanonicalLtd/iot-devicetwin/config"
	"github.com/CanonicalLtd/iot-devicetwin/datastore/memory"
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
	p5 := []byte(`{"id":"a1", "action":"install", "success":true, "message":"", "result": "101"}`)
	p6 := []byte(`{"id":"a1", "action":"conf", "success":true, "message":"", "result": {"name":"abc", "status":"active", "version":"1.0", "config":"{\"title\": \"Jack\"}"}}`)
	p7 := []byte(`{"id":"a1", "action":"server", "success":true, "message":"", "result": {"deviceId":"a111", "osVersionId":"core-123", "series":"16", "kernelVersion":"kernel-123"}}`)

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

		{"valid-install", args{"a111", "install", p5}, false},
		{"install-empty-payload", args{"a111", "install", p1}, true},

		{"valid-conf", args{"a111", "conf", p6}, false},
		{"conf-no-device", args{"invalid", "conf", p6}, true},
		{"conf-empty-payload", args{"a111", "conf", p1}, true},

		{"valid-server", args{"a111", "server", p7}, false},
		{"server-no-device", args{"invalid", "server", p7}, true},
		{"server-empty-payload", args{"a111", "server", p1}, true},
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

func TestService_ActionCreate(t *testing.T) {
	a1 := domain.SubscribeAction{
		ID:     "aa1234",
		Action: "install",
		Snap:   "helloworld",
	}
	type args struct {
		orgID    string
		deviceID string
		action   domain.SubscribeAction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111", a1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.ActionCreate(tt.args.orgID, tt.args.deviceID, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("Service.ActionCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GroupWorkflow(t *testing.T) {
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
		{"invalid", args{"invalid", "test-group"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.GroupCreate(tt.args.orgID, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Service.GroupCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
