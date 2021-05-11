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
		{"invalid", args{"abc", "invalid"}, 0, true},
		{"invalid-orgid", args{"invalid", "a111"}, 0, true},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			got, err := srv.DeviceSnaps(localtt.args.orgID, localtt.args.clientID)
			if (err != nil) != localtt.wantErr {
				t.Errorf("Service.DeviceSnaps() error = %v, wantErr %v", err, localtt.wantErr)
				return
			}
			if localtt.wantErr {
				return
			}

			if len(got) != localtt.want {
				t.Errorf("Service.DeviceSnaps() = %v, want %v", len(got), localtt.want)
			}
		})
	}
}
