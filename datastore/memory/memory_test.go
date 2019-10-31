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

package memory

import (
	"testing"
	"time"

	"github.com/CanonicalLtd/iot-devicetwin/datastore"
)

func TestStore_DeviceGet(t *testing.T) {

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"b222"}, false},
		{"invalid", args{"does-not-exist"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			got, err := mem.DeviceGet(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.DeviceGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.DeviceID != tt.args.id {
					t.Errorf("Store.DeviceGet() device ID = %v, wantErr %v", got.DeviceID, tt.args.id)
				}
			}
		})
	}
}

func TestStore_DevicePing(t *testing.T) {
	type args struct {
		id      string
		refresh time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{"a111", time.Now()}, false},
		{"invalid", args{"does-not-exist", time.Now()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			if err := mem.DevicePing(tt.args.id, tt.args.refresh); (err != nil) != tt.wantErr {
				t.Errorf("Store.DevicePing() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Check the device is updated
			got, err := mem.DeviceGet(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.DeviceGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.LastRefresh != tt.args.refresh {
				t.Errorf("Store.DeviceGet() device refresh = %v, wantErr %v", got.LastRefresh, tt.args.refresh)
			}
		})
	}
}

func TestStore_DeviceCreate(t *testing.T) {
	d1 := datastore.Device{
		OrganisationID: "abc",
		DeviceID:       "d444",
		Brand:          "example",
		Model:          "drone-1000",
		SerialNumber:   "D44444444",
	}
	d2 := datastore.Device{
		OrganisationID: "abc",
		DeviceID:       "a111",
		Brand:          "example",
		Model:          "drone-1000",
		SerialNumber:   "D44444444",
	}
	type args struct {
		device datastore.Device
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"valid", args{d1}, 4, false},
		{"duplicate", args{d2}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			got, err := mem.DeviceCreate(tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.DeviceCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got != tt.want {
					t.Errorf("Store.DeviceCreate() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestStore_ActionWorkflow(t *testing.T) {
	type args struct {
		act datastore.Action
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"valid", args{datastore.Action{ActionID: "a1", Action: "device", DeviceID: "a111"}}, 3, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			got, err := mem.ActionCreate(tt.args.act)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.ActionCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Store.ActionCreate() = %v, want %v", got, tt.want)
			}

			err = mem.ActionUpdate(tt.args.act.ActionID, "complete", "Done")
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.ActionUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actions, err := mem.ActionListForDevice("abc", tt.args.act.DeviceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.ActionListForDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(actions) != 1 {
				t.Errorf("Store.ActionListForDevice() = %v, want %v", len(actions), 1)
			}

			a := actions[0]
			if a.ActionID != tt.args.act.ActionID && a.DeviceID != tt.args.act.DeviceID && a.Action != tt.args.act.Action && a.Status != "complete" {
				t.Error("Store.ActionListForDevice() = store action is invalid")
			}
		})
	}
}

func TestStore_DeviceVersionWorkflow(t *testing.T) {
	type args struct {
		dv datastore.DeviceVersion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{datastore.DeviceVersion{DeviceID: 1, OSID: "123", KernelVersion: "kernel-123", OSVersionID: "core-123"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			if err := mem.DeviceVersionUpsert(tt.args.dv); (err != nil) != tt.wantErr {
				t.Errorf("Store.DeviceVersionUpsert() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			// Get the created record
			dv, err := mem.DeviceVersionGet(tt.args.dv.DeviceID)
			if err != nil {
				t.Errorf("Store.DeviceVersionGet() error = %v", err)
			}
			if dv.OSVersionID != tt.args.dv.OSVersionID {
				t.Errorf("Store.DeviceVersionGet() OS version = %v, want %v", dv.OSVersionID, tt.args.dv.OSVersionID)
			}

			// Update the record
			tt.args.dv.OSVersionID = "changed"
			if err := mem.DeviceVersionUpsert(tt.args.dv); err != nil {
				t.Errorf("Store.DeviceVersionUpsert() error update = %v", err)
			}

			// Get the updated record
			dv2, err := mem.DeviceVersionGet(tt.args.dv.DeviceID)
			if err != nil {
				t.Errorf("Store.DeviceVersionGet() error = %v", err)
			}
			if dv2.OSVersionID != tt.args.dv.OSVersionID {
				t.Errorf("Store.DeviceVersionGet() OS version updated = %v, want %v", dv2.OSVersionID, tt.args.dv.OSVersionID)
			}

			// Delete the record
			if err := mem.DeviceVersionDelete(dv2.ID); err != nil {
				t.Errorf("Store.DeviceVersionDelete() error update = %v", err)
			}

			// Check the record is deleted
			if _, err := mem.DeviceVersionGet(tt.args.dv.DeviceID); err == nil {
				t.Error("Store.DeviceVersionDelete() error delete check failed")
			}
		})
	}
}

func TestStore_DeviceVersionDelete(t *testing.T) {
	dv1 := datastore.DeviceVersion{DeviceID: 1, OSID: "123", KernelVersion: "kernel-123", OSVersionID: "core-123"}
	type args struct {
		dv datastore.DeviceVersion
		id int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid-delete", args{dv1, 1}, false},
		{"invalid-delete", args{dv1, 999}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()

			if err := mem.DeviceVersionUpsert(tt.args.dv); err != nil {
				t.Errorf("Store.DeviceVersionUpsert() error = %v", err)
			}

			if err := mem.DeviceVersionDelete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Store.DeviceVersionDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStore_DeviceList(t *testing.T) {
	tests := []struct {
		name    string
		orgID   string
		want    int
		wantErr bool
	}{
		{"valid", "abc", 3, false},
		{"valid-no-devices", "none", 0, false},
		{"valid", "invalid", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			got, err := mem.DeviceList(tt.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.DeviceList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Store.DeviceList() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestStore_GroupCreate(t *testing.T) {
	type args struct {
		orgID string
		name  string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"valid", args{"abc", "test-group"}, 2, false},
		{"valid-exists", args{"abc", "workshop"}, 0, true},
		{"invalid", args{"invalid", "test-group"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			got, err := mem.GroupCreate(tt.args.orgID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GroupCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Store.GroupCreate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStore_GroupDeviceWorkflow(t *testing.T) {
	type args struct {
		orgID  string
		name   string
		device string
	}
	tests := []struct {
		name    string
		args    args
		count   int
		wantErr bool
	}{
		{"valid", args{"abc", "workshop", "c333"}, 2, false},
		{"valid-exists", args{"abc", "workshop", "a111"}, 1, false},
		{"invalid-org", args{"invalid", "workshop", "c333"}, 0, true},
		{"invalid-group", args{"abc", "invalid", "c333"}, 0, true},
		{"invalid-device", args{"abc", "workshop", "invalid"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			// Get the group
			_, err := mem.GroupGet(tt.args.orgID, tt.args.name)
			if (err != nil) != tt.wantErr && tt.args.device != "invalid" {
				t.Errorf("Store.GroupGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Link device to the group
			err = mem.GroupLinkDevice(tt.args.orgID, tt.args.name, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GroupLinkDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Get devices for the group
			if tt.args.device != "invalid" {
				devices, err := mem.GroupGetDevices(tt.args.orgID, tt.args.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("Store.GroupGetDevices() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if len(devices) != tt.count {
					t.Errorf("Store.GroupGetDevices() count = %v, expected %v", len(devices), tt.count)
					return
				}
			}

			// Unlink device from the group
			err = mem.GroupUnlinkDevice(tt.args.orgID, tt.args.name, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GroupUnlinkDevice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Get devices for the group
			if tt.args.device != "invalid" {
				devices2, err := mem.GroupGetDevices(tt.args.orgID, tt.args.name)
				if (err != nil) != tt.wantErr {
					t.Errorf("Store.GroupGetDevices() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.count > 0 {
					tt.count = tt.count - 1
				}
				if len(devices2) != tt.count {
					t.Errorf("Store.GroupGetDevices() count = %v, expected %v", len(devices2), tt.count)
					return
				}
			}
		})
	}
}

func TestStore_GroupGetExcludedDevices(t *testing.T) {
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
		{"invalid-not-exists", args{"abc", "does-not-exist"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := NewStore()
			got, err := mem.GroupGetExcludedDevices(tt.args.orgID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Store.GroupGetExcludedDevices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("Store.GroupGetExcludedDevices() = %v, want %v", got, tt.want)
			}
		})
	}
}
