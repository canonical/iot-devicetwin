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
	"io"
	"strings"
	"testing"

	"github.com/CanonicalLtd/iot-devicetwin/config"
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

func TestService_GroupList(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc", "GET", nil, 200, ""},
		{"valid-org", "/v1/group/invalid", "GET", nil, 400, "GroupList"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupList() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupList() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupList() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_GroupGet(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc/workshop", "GET", nil, 200, ""},
		{"invalid-org", "/v1/group/invalid/workshop", "GET", nil, 400, "GroupGet"},
		{"invalid-name", "/v1/group/abc/invalid", "GET", nil, 400, "GroupGet"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupGet() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupGet() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupGet() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_GroupLinkDevice(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc/workshop/c333", "POST", nil, 200, ""},
		{"invalid-org", "/v1/group/invalid/workshop/c333", "POST", nil, 400, "GroupLink"},
		{"invalid-name", "/v1/group/abc/invalid/c333", "POST", nil, 400, "GroupLink"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupLinkDevice() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupLinkDevice() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupLinkDevice() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_GroupUnlinkDevice(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc/workshop/c333", "DELETE", nil, 200, ""},
		{"invalid-org", "/v1/group/invalid/workshop/c333", "DELETE", nil, 400, "GroupUnlink"},
		{"invalid-name", "/v1/group/abc/invalid/c333", "DELETE", nil, 400, "GroupUnlink"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupUnlinkDevice() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupUnlinkDevice() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupUnlinkDevice() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_GroupGetDevices(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc/workshop/devices", "GET", nil, 200, ""},
		{"invalid-org", "/v1/group/invalid/workshop/devices", "GET", nil, 400, "GroupDevices"},
		{"invalid-name", "/v1/group/abc/invalid/devices", "GET", nil, 400, "GroupDevices"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupGetDevices() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupGetDevices() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupGetDevices() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}

func TestService_GroupGetExcludedDevices(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		method string
		data   io.Reader
		code   int
		result string
	}{
		{"valid", "/v1/group/abc/workshop/devices/excluded", "GET", nil, 200, ""},
		{"invalid-org", "/v1/group/invalid/workshop/devices/excluded", "GET", nil, 400, "GroupDevices"},
		{"invalid-name", "/v1/group/abc/invalid/devices/excluded", "GET", nil, 400, "GroupDevices"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := NewService(config.TestConfig(), testController())
			w := sendRequest(tt.method, tt.url, tt.data, wb)
			if w.Code != tt.code {
				t.Errorf("Web.GroupGetExcludedDevices() got = %v, want %v", w.Code, tt.code)
			}
			resp, err := parseStandardResponse(w.Body)
			if err != nil {
				t.Errorf("Web.GroupGetExcludedDevices() got = %v", err)
			}
			if resp.Code != tt.result {
				t.Errorf("Web.GroupGetExcludedDevices() got = %v, want %v", resp.Code, tt.result)
			}
		})
	}
}
