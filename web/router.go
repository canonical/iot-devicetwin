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
	"time"

	"github.com/gorilla/mux"
)

// Router returns the application router
func (wb Service) Router() *mux.Router {
	// Start the web service router
	router := mux.NewRouter()

	// Actions on a device twin
	router.Handle("/v1/device/{orgid}/{id}/snaps", Middleware(http.HandlerFunc(wb.SnapList))).Methods("GET")
	router.Handle("/v1/device/{orgid}", Middleware(http.HandlerFunc(wb.DeviceList))).Methods("GET")
	router.Handle("/v1/device/{orgid}/{id}", Middleware(http.HandlerFunc(wb.DeviceGet))).Methods("GET")
	router.Handle("/v1/device/{orgid}/{id}/actions", Middleware(http.HandlerFunc(wb.ActionList))).Methods("GET")

	// Actions on a device
	router.Handle("/v1/device/{orgid}/{id}/snaps/list", Middleware(http.HandlerFunc(wb.SnapListPublish))).Methods("POST")
	router.Handle("/v1/device/{orgid}/{id}/snaps/{snap}", Middleware(http.HandlerFunc(wb.SnapInstall))).Methods("POST")
	router.Handle("/v1/device/{orgid}/{id}/snaps/{snap}", Middleware(http.HandlerFunc(wb.SnapRemove))).Methods("DELETE")
	router.Handle("/v1/device/{orgid}/{id}/snaps/{snap}/settings", Middleware(http.HandlerFunc(wb.SnapUpdateConf))).Methods("PUT")
	router.Handle("/v1/device/{orgid}/{id}/snaps/{snap}/{action}", Middleware(http.HandlerFunc(wb.SnapUpdateAction))).Methods("PUT")

	// Actions on a group
	router.Handle("/v1/group/{orgid}", Middleware(http.HandlerFunc(wb.GroupCreate))).Methods("POST")
	router.Handle("/v1/group/{orgid}", Middleware(http.HandlerFunc(wb.GroupList))).Methods("GET")
	router.Handle("/v1/group/{orgid}/{name}", Middleware(http.HandlerFunc(wb.GroupGet))).Methods("GET")
	router.Handle("/v1/group/{orgid}/{name}/{id}", Middleware(http.HandlerFunc(wb.GroupLinkDevice))).Methods("POST")
	router.Handle("/v1/group/{orgid}/{name}/{id}", Middleware(http.HandlerFunc(wb.GroupUnlinkDevice))).Methods("DELETE")
	router.Handle("/v1/group/{orgid}/{name}/devices", Middleware(http.HandlerFunc(wb.GroupGetDevices))).Methods("GET")
	router.Handle("/v1/group/{orgid}/{name}/devices/excluded", Middleware(http.HandlerFunc(wb.GroupGetExcludedDevices))).Methods("GET")

	return router
}

// Logger Handle logging for the web service
func Logger(start time.Time, r *http.Request) {
	log.Printf(
		"%s\t%s\t%s",
		r.Method,
		r.RequestURI,
		time.Since(start),
	)
}

// Middleware to pre-process web service requests
func Middleware(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request
		Logger(start, r)

		inner.ServeHTTP(w, r)
	})
}
