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

	"github.com/everactive/iot-devicetwin/config"
	"github.com/everactive/iot-devicetwin/datastore/memory"
	"github.com/everactive/iot-devicetwin/pkg/messages"
)

func TestService_HealthHandler(t *testing.T) {
	h1 := messages.Health{OrgId: "abc", DeviceId: "a111"}
	h2 := messages.Health{OrgId: "abc", DeviceId: "invalid"}
	type args struct {
		payload messages.Health
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
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.HealthHandler(localtt.args.payload); (err != nil) != localtt.wantErr {
				t.Errorf("Service.HealthHandler() error = %v, wantErr %v", err, localtt.wantErr)
			}
		})
	}
}

func TestService_ActionResponse(t *testing.T) {
	p1 := []byte("")
	p2 := []byte(`{"id":"a1", "action":"device", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"d444", "brand":"example", "model":"drone-1000", "serial":"d444"}}`)
	p2a := []byte(`{"id":"a1", "action":"device", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"d444", "brand":"example", "model":"drone-1000", "serial":"d444", "version":{"deviceId":"d444", "version":"1.0"}}}`)
	p3 := []byte(`{"id":"a1", "action":"device", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"a111", "brand":"example", "model":"drone-1000", "serial":"d444"}}`)
	p4 := []byte(`{"id":"a1", "action":"list", "success":true, "message":"", "result": [{"name":"abc", "status":"active", "version":"1.0"}, {"name":"alpaca", "status":"active", "version":"2.3"}]}`)
	p5 := []byte(`{"id":"a1", "action":"install", "success":true, "message":"", "result": "101"}`)
	p6 := []byte(`{"id":"a1", "action":"conf", "success":true, "message":"", "result": {"name":"abc", "status":"active", "version":"1.0", "config":"{\"title\": \"Jack\"}"}}`)
	p7 := []byte(`{"id":"a1", "action":"server", "success":true, "message":"", "result": {"deviceId":"a111", "osVersionId":"core-123", "series":"16", "kernelVersion":"kernel-123"}}`)
	p8 := []byte(`{"id":"a1", "action":"unregister", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"d444", "brand":"example", "model":"drone-1000", "serial":"d444"}}`)

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
		{"valid-device-version", args{"d444", "device", p2a}, false},
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

		{"valid-unregister", args{"a111", "unregister", p8}, false},
		{"server-no-device", args{"invalid", "unregister", p8}, true},
		{"server-empty-payload", args{"a111", "unregister", p1}, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.ActionResponse(localtt.args.clientID, "a1", localtt.args.action, localtt.args.payload); (err != nil) != localtt.wantErr {
				t.Errorf("Service.ActionResponse() error = %v, wantErr %v", err, localtt.wantErr)
			}
		})
	}
}

func TestService_ActionCreate(t *testing.T) {
	a1 := messages.SubscribeAction{
		Id:     "aa1234",
		Action: "install",
		Snap:   "helloworld",
	}
	type args struct {
		orgID    string
		deviceID string
		action   messages.SubscribeAction
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"abc", "a111", a1}, false},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.ActionCreate(localtt.args.orgID, localtt.args.deviceID, localtt.args.action); (err != nil) != localtt.wantErr {
				t.Errorf("Service.ActionCreate() error = %v, wantErr %v", err, localtt.wantErr)
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
		count   int
		wantErr bool
	}{
		{"valid", args{"abc", "test-group"}, 2, false},
		{"invalid", args{"invalid", "test-group"}, 0, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			if err := srv.GroupCreate(localtt.args.orgID, localtt.args.name); (err != nil) != localtt.wantErr {
				t.Errorf("Service.GroupCreate() error = %v, wantErr %v", err, localtt.wantErr)
			}

			groups, err := srv.GroupList(localtt.args.orgID)
			if (err != nil) != localtt.wantErr {
				t.Errorf("Service.GroupList() error = %v, wantErr %v", err, localtt.wantErr)
			}
			if len(groups) != localtt.count {
				t.Errorf("Service.GroupList() count = %v, wantErr %v", len(groups), localtt.count)
			}
		})
	}
}

func TestService_GroupLinkWorkflow(t *testing.T) {
	type args struct {
		orgID    string
		name     string
		deviceID string
	}
	tests := []struct {
		name    string
		args    args
		count   int
		wantErr bool
	}{
		{"valid", args{"abc", "workshop", "c333"}, 2, false},
		{"invalid", args{"invalid", "workshop", "c333"}, 0, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())

			// Get a group
			group, err := srv.GroupGet(localtt.args.orgID, localtt.args.name)
			if (err != nil) != localtt.wantErr {
				t.Errorf("Service.GroupGet() error = %v, wantErr %v", err, localtt.wantErr)
			}
			if !localtt.wantErr {
				if group.Name != localtt.args.name {
					t.Errorf("Service.GroupGet() name = %v, wantErr %v", group.Name, localtt.args.name)
				}
			}

			// Link a device to a group
			if err = srv.GroupLinkDevice(localtt.args.orgID, localtt.args.name, localtt.args.deviceID); (err != nil) != localtt.wantErr {
				t.Errorf("Service.GroupLinkDevice() error = %v, wantErr %v", err, localtt.wantErr)
			}

			// Get the devices for the group
			devices, err2 := srv.GroupGetDevices(localtt.args.orgID, localtt.args.name)
			if (err2 != nil) != localtt.wantErr {
				t.Errorf("Service.GroupGetDevices() error = %v, wantErr %v", err2, localtt.wantErr)
			}
			if len(devices) != localtt.count {
				t.Errorf("Service.GroupGetDevices() count = %v, wantErr %v", len(devices), localtt.count)
			}

			// Unlink a device from a group
			if err = srv.GroupUnlinkDevice(localtt.args.orgID, localtt.args.name, localtt.args.deviceID); (err != nil) != localtt.wantErr {
				t.Errorf("Service.GroupUnlinkDevice() error = %v, wantErr %v", err, localtt.wantErr)
			}

			// Get the devices for the group
			devices2, err3 := srv.GroupGetDevices(localtt.args.orgID, localtt.args.name)
			if (err3 != nil) != localtt.wantErr {
				t.Errorf("Service.GroupGetDevices() error = %v, wantErr %v", err3, localtt.wantErr)
			}
			if localtt.count > 0 {
				localtt.count--
			}
			if len(devices2) != localtt.count {
				t.Errorf("Service.GroupGetDevices() count = %v, wantErr %v", len(devices2), localtt.count)
			}
		})
	}
}
