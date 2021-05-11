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
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// DeviceGet is the API call to get a device
func (wb Service) DeviceGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	device, err := wb.Controller.DeviceGet(vars["orgid"], vars["id"])
	if err != nil {
		log.Printf("Error fetching the device `%s`: %v", vars["id"], err)
		formatStandardResponse("DeviceGet", "Error fetching the device", w)
		return
	}

	formatDeviceResponse(device, w)
}

// DeviceList is the API call to get devices
func (wb Service) DeviceList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	devices, err := wb.Controller.DeviceList(vars["orgid"])
	if err != nil {
		log.Printf("Error fetching the device list for `%s`: %v", vars["orgid"], err)
		formatStandardResponse("DeviceList", "Error fetching devices", w)
		return
	}

	formatDevicesResponse(devices, w)
}

// DeviceDelete is the API call to delete devices
func (wb Service) DeviceDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := wb.Controller.DeviceDelete(vars["id"])
	if err != nil {
		log.Printf("Error deleting the device `%s`: %v", vars["id"], err)
		formatStandardResponse("DeviceDelete", "Error deleting device", w)
		return
	}

	w.Header().Set("Content-Type", JSONHeader)
	response := StandardResponse{
		Message: "device deleted",
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}

// DeviceUnregister is the API call to send unregister event to the device
func (wb Service) DeviceUnregister(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := wb.Controller.DeviceUnregister(vars["orgid"], vars["id"])
	if err != nil {
		log.Printf("Error unregistering the device `%s`: %v", vars["id"], err)
		formatStandardResponse("DeviceUnregister", "Error unregistering device", w)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", JSONHeader)
	response := StandardResponse{
		Message: "device unregister sent",
	}

	// Encode the response as JSON
	encodeResponse(w, response)
}
