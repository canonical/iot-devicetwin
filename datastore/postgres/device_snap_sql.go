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

const createDeviceSnapTableSQL = `
CREATE TABLE IF NOT EXISTS device_snap (
   id             serial primary key,
   created        timestamp default current_timestamp,
   modified       timestamp default current_timestamp,
   device_id      int references device not null,
   name           varchar(200) not null,
   installed_size int default 0,
   installed_date timestamp default current_timestamp,
   status         varchar(200) default '',
   channel        varchar(200) default '',
   confinement    varchar(200) default '',
   version        varchar(200) default '',
   revision       int default 0,
   devmode        bool default false,
   config         varchar(200) default ''
)
`

const createDeviceSnapIndexSQL = "CREATE UNIQUE INDEX IF NOT EXISTS device_snap_idx ON device_snap (device_id, name)"

const upsertDeviceSnapSQL = `
INSERT INTO device_snap(device_id, name, installed_size, installed_date, status, channel, confinement, version, revision, devmode, config)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT (device_id, name)
DO
  UPDATE
  SET installed_size = EXCLUDED.installed_size, 
      installed_date = EXCLUDED.installed_date,
      status = EXCLUDED.status,
      channel = EXCLUDED.channel,
      confinement = EXCLUDED.confinement,
      version = EXCLUDED.version,
      revision = EXCLUDED.revision,
      devmode = EXCLUDED.devmode,
      config = EXCLUDED.config,
      modified = current_timestamp
  RETURNING id;
`

const listDeviceSnapSQL = `
select id, created, modified, device_id, name, installed_size, installed_date, status, channel, confinement, version, revision, devmode, config
from device_snap
where device_id=$1
order by name`

const deleteDeviceSnapSQL = `
delete from device_snap where device_id=$1`
