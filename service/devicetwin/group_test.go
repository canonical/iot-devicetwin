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
	"github.com/canonical/iot-devicetwin/config"
	"github.com/canonical/iot-devicetwin/datastore/memory"
	"testing"
)

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
		{"invalid", args{"abc", "does-not-exist"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
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
