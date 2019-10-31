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

// ActionList is the API call to list actions for a device
func (wb Service) ActionList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	actions, err := wb.Controller.ActionList(vars["orgid"], vars["id"])
	if err != nil {
		log.Printf("Error fetching the actions for `%s`: %v", vars["id"], err)
		formatStandardResponse("ActionList", "Error fetching the actions", w)
		return
	}

	formatActionsResponse(actions, w)
}
