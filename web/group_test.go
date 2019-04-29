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

package web

import (
	"github.com/CanonicalLtd/iot-devicetwin/config"
	"io"
	"strings"
	"testing"
)

func TestService_GroupCreate(t *testing.T) {
	a1 := strings.NewReader(`{"organizationId":"abc", "name":"test-group"}`)
	a2 := strings.NewReader(`{"organizationId":"invalid", "name":"test-group"}`)
	a3 := strings.NewReader(`\u1000`)

	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc", "POST", a1, 200, ""},
		{"invalid-org", "/v1/group/invalid", "POST", a2, 400, "GroupCreate"},
		{"invalid-body", "/v1/group/abc", "POST", a3, 400, "GroupCreate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupCreate() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupCreate() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupCreate() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}
