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

func (wb Service) SnapList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	installed, err := wb.Controller.DeviceSnaps(vars["id"])
	if err != nil {
		log.Println("Error fetching snaps for a device:", err)
		formatStandardResponse("SnapList", "Error fetching snaps for the device", w)
		return
	}

	formatSnapsResponse(installed, w)
}
