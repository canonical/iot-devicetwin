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

package mqtt

import (
	"testing"

	"github.com/canonical/iot-devicetwin/config"
)

func TestGetConnection(t *testing.T) {
	p1 := `{"id":"a1", "action":"device", "success":true, "message":"", "result": {"orgId":"abc", "deviceId":"d444", "brand":"example", "model":"drone-1000", "serial":"d444"}}`
	client = &MockClient{}

	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := GetConnection(config.TestConfig())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if conn == nil {
				t.Error("GetConnection() error = connection is nil")
			}

			err = conn.Publish("device/sub/a111", p1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = conn.Subscribe("device/pub/a111", nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			conn.Close()
			if conn.client.IsConnectionOpen() {
				t.Error("Close() error = connection not closed")
			}
		})
	}
}

func Test_newTLSConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"valid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTLSConfig(config.TestConfig())
			if (err != nil) != tt.wantErr {
				t.Errorf("newTLSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("newTLSConfig() error = certificate is nil")
			}
		})
	}
}
