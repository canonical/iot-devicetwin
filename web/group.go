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
	"encoding/json"
	"github.com/CanonicalLtd/iot-devicetwin/domain"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

// GroupCreate is the API call to create a group
func (wb Service) GroupCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	defer r.Body.Close()
	group, err := parseGroupRequest(r.Body)
	if err != nil {
		log.Printf("Error parsing the group for organization `%s`: %v", vars["orgid"], err)
		formatStandardResponse("GroupCreate", "Error creating the group", w)
		return
	}

	err = wb.Controller.GroupCreate(vars["orgid"], group.Name)
	if err != nil {
		log.Printf("Error creating the group for organization `%s`: %v", vars["orgid"], err)
		formatStandardResponse("GroupCreate", "Error creating the group", w)
		return
	}

	formatStandardResponse("", "", w)
}

// GroupList is the API call to list groups
func (wb Service) GroupList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	groups, err := wb.Controller.GroupList(vars["orgid"])
	if err != nil {
		log.Printf("Error listing the groups for organization `%s`: %v", vars["orgid"], err)
		formatStandardResponse("GroupList", "Error listing the groups", w)
		return
	}

	formatGroupsResponse(groups, w)
}

// GroupGet is the API call to retrieve a group
func (wb Service) GroupGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	group, err := wb.Controller.GroupGet(vars["orgid"], vars["name"])
	if err != nil {
		log.Printf("Error fetching the group for organization `%s`: %v", vars["orgid"], err)
		formatStandardResponse("GroupGet", "Error fetching the group", w)
		return
	}

	formatGroupResponse(group, w)
}

// GroupLinkDevice is the API call to link a device to a group
func (wb Service) GroupLinkDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := wb.Controller.GroupLinkDevice(vars["orgid"], vars["name"], vars["id"]); err != nil {
		log.Printf("Error linking the device and group `%s` - `%s`: %v", vars["id"], vars["name"], err)
		formatStandardResponse("GroupLink", "Error linking the device to the group", w)
		return
	}

	formatStandardResponse("", "", w)
}

// GroupUnlinkDevice is the API call to unlink a device from a group
func (wb Service) GroupUnlinkDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if err := wb.Controller.GroupUnlinkDevice(vars["orgid"], vars["name"], vars["id"]); err != nil {
		log.Printf("Error unlinking the device and group `%s` - `%s`: %v", vars["id"], vars["name"], err)
		formatStandardResponse("GroupUnlink", "Error unlinking the device to the group", w)
		return
	}

	formatStandardResponse("", "", w)
}

// GroupGetDevices is the API call to get the devices for a group
func (wb Service) GroupGetDevices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	devices, err := wb.Controller.GroupGetDevices(vars["orgid"], vars["name"])
	if err != nil {
		log.Printf("Error fetching the devices for group `%s`: %v", vars["name"], err)
		formatStandardResponse("GroupDevices", "Error fetching the devices for the group", w)
		return
	}

	formatDevicesResponse(devices, w)
}

func parseGroupRequest(r io.Reader) (domain.Group, error) {
	result := domain.Group{}
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}
