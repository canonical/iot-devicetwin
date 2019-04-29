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

// GroupCreate is the API call to create a device
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

func parseGroupRequest(r io.Reader) (domain.Group, error) {
	result := domain.Group{}
	err := json.NewDecoder(r).Decode(&result)
	return result, err
}
