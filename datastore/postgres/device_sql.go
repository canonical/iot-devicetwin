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

package postgres

const createDeviceTableSQL = `
CREATE TABLE IF NOT EXISTS device (
   id             serial primary key,
   created        timestamp default current_timestamp,
   lastrefresh    timestamp default current_timestamp,
   org_id         varchar(200) not null,
   device_id      varchar(200) unique not null,
   brand          varchar(200) not null,
   model          varchar(200) not null,
   serial         varchar(200) not null,
   store_id       varchar(200) not null,
   device_key     text,
   active         bool default true
)
`

const createDeviceSQL = `
insert into device (org_id, device_id, brand, model, serial, store_id, device_key)
values ($1,$2,$3,$4,$5,$6,$7) RETURNING id`

const getDeviceSQL = `
select id, created, lastrefresh, org_id, device_id, brand, model, serial, store_id, device_key, active
from device
where device_id=$1`

const listDeviceSQL = `
select id, created, lastrefresh, org_id, device_id, brand, model, serial, store_id, device_key, active
from device
where org_id=$1
order by brand, model, serial`

const pingDeviceSQL = `
update device
set lastrefresh=$2
where device_id=$1`

const deleteDeviceSQL = `
delete from device where device_id=$1
`

const deleteDeviceSnapsSQL = `
delete from device_snap where device_id=$1
`
