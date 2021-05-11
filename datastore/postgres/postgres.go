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

// Package postgres is the Datastore implementation for Postgres
package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // postgresql driver
)

// DataStore is the postgreSQL implementation of a data store
type DataStore struct {
	driver string
	*sql.DB
}

var pgStore *DataStore

// OpenDataStore returns an open database connection
func OpenDataStore(driver, dataSource string) *DataStore {
	if pgStore != nil {
		return pgStore
	}

	// Open the database
	pgStore = openDatabase(driver, dataSource)

	// Create the tables, if needed
	pgStore.createTables()

	return pgStore
}

// openDatabase return an open database connection for a postgreSQL database
func openDatabase(driver, dataSource string) *DataStore {
	// Open the database connection
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		log.Fatalf("Error opening the database: %v\n", err)
	}

	// Check that we have a valid database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error accessing the database: %v\n", err)
	}

	return &DataStore{driver, db}
}

func (db *DataStore) createTables() {
	_ = db.createDeviceTable()
	_ = db.createActionTable()
	_ = db.createDeviceSnapTable()
	_ = db.createDeviceVersionTable()
	_ = db.createOrgGroupTable()
}
