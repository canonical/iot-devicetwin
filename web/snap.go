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
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// SnapList is the API call to list snaps for a device
func (wb Service) SnapList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	installed, err := wb.Controller.DeviceSnaps(vars["orgid"], vars["id"])
	if err != nil {
		log.Println("Error fetching snaps for a device:", err)
		formatStandardResponse("SnapList", "Error fetching snaps for the device", w)
		return
	}

	formatSnapsResponse(installed, w)
}

// SnapInstall is the API call to install a snap for a device
func (wb Service) SnapInstall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := wb.Controller.DeviceSnapInstall(vars["orgid"], vars["id"], vars["snap"])
	if err != nil {
		log.Println("Error requesting snap install for the device:", err)
		formatStandardResponse("SnapInstall", "Error requesting snap install for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapRemove is the API call to uninstall a snap for a device
func (wb Service) SnapRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := wb.Controller.DeviceSnapRemove(vars["orgid"], vars["id"], vars["snap"])
	if err != nil {
		log.Println("Error requesting snap remove for the device:", err)
		formatStandardResponse("SnapRemove", "Error requesting snap remove for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}

// SnapUpdateAction is the API call to update a snap for a device (enable, disable, refresh)
func (wb Service) SnapUpdateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := wb.Controller.DeviceSnapUpdate(vars["orgid"], vars["id"], vars["snap"], vars["action"])
	if err != nil {
		log.Println("Error requesting snap update for the device:", err)
		formatStandardResponse("SnapUpdate", "Error requesting snap update for the device", w)
		return
	}

	formatStandardResponse("", "", w)
}
