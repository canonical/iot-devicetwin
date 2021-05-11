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

func TestService_ActionList(t *testing.T) {
	type args struct {
		orgID    string
		deviceID string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"valid", args{"abc", "c333"}, 2, false},
		{"valid", args{"abc", "a111"}, 0, false},
	}
	for _, tt := range tests {
		localtt := tt
		t.Run(localtt.name, func(t *testing.T) {
			srv := NewService(config.TestConfig(), memory.NewStore())
			got, err := srv.ActionList(localtt.args.orgID, localtt.args.deviceID)
			if (err != nil) != localtt.wantErr {
				t.Errorf("ActionList() error = %v, wantErr %v", err, localtt.wantErr)
				return
			}
			if len(got) != localtt.want {
				t.Errorf("ActionList() got = %v, want %v", len(got), localtt.want)
			}
		})
	}
}
