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
		{"valid", args{datastore.Action{ActionID: "a1", Action: "device", DeviceID: "a111"}}, 1, false},
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

			actions, err := mem.ActionListForDevice(tt.args.act.DeviceID)
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
