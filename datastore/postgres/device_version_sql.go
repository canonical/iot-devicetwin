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

const createDeviceVersionTableSQL = `
CREATE TABLE IF NOT EXISTS device_version (
   id             serial primary key,
   device_id      int references device not null unique,
   version        varchar(200) not null,
   series         varchar(200) default '',
   os_id          varchar(200) default '',
   os_version_id  varchar(200) default '',
   on_classic     bool default false,
   kernel_version varchar(200) default ''
)
`

const getDeviceVersionSQL = `
select id, device_id, version, series, os_id, os_version_id, on_classic, kernel_version
from device_version
where device_id=$1`

const upsertDeviceVersionSQL = `
INSERT INTO device_version (device_id, version, series, os_id, os_version_id, on_classic, kernel_version)
VALUES($1,$2,$3,$4,$5,$6,$7)
ON CONFLICT (device_id)
DO
  UPDATE
  SET version = EXCLUDED.version, 
      series = EXCLUDED.series,
      os_id = EXCLUDED.os_id,
      os_version_id = EXCLUDED.os_version_id,
      on_classic = EXCLUDED.on_classic,
      kernel_version = EXCLUDED.kernel_version
  RETURNING id;
`

const deleteDeviceVersionSQL = `delete from device_version where id=$1`
