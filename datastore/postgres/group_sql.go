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

const createOrgGroupTableSQL = `
CREATE TABLE IF NOT EXISTS org_group (
   id             serial primary key,
   created        timestamp default current_timestamp,
   modified       timestamp default current_timestamp,
   org_id         varchar(200) not null,
   name           varchar(200) not null
)
`

const createGroupDeviceLinkTableSQL = `
CREATE TABLE IF NOT EXISTS group_device_link (
   id             serial primary key,
   created        timestamp default current_timestamp,
   org_id         varchar(200) not null,
   group_id       int references org_group not null,
   device_id      int references device not null
)
`

const createOrgGroupIndexSQL = "CREATE INDEX IF NOT EXISTS org_group_idx ON org_group (org_id, name)"

const createOrgGroupSQL = `
insert into org_group (org_id, name)
values ($1,$2) RETURNING id`

const listOrgGroupSQL = `
select id, created, modified, org_id, name
from org_group
where org_id=$1
order by name`

const getOrgGroupSQL = `
select id, created, modified, org_id, name
from org_group
where org_id=$1 and name=$2`

const createGroupDeviceLinkSQL = `
insert into group_device_link (org_id, group_id, device_id)
values ($1,$2,$3)`

const deleteGroupDeviceLinkSQL = `delete from group_device_link where group_id=$1 and device_id=$2`

const listGroupDeviceLinkSQL = `
select d.id, d.created, d.lastrefresh, d.org_id, d.device_id, d.brand, d.model, d.serial, d.store_id, d.device_key, d.active
from device d
inner join group_device_link lnk on lnk.device_id=d.id
where lnk.org_id=$1 and lnk.group_id=$2
order by d.brand, d.model, d.serial
`

const listGroupDeviceExcludedLinkSQL = `
select d.id, d.created, d.lastrefresh, d.org_id, d.device_id, d.brand, d.model, d.serial, d.store_id, d.device_key, d.active
from device d
where not exists (
   select device_id from group_device_link
   where device_id = d.id
     and org_id=$1 and group_id=$2
 )
order by d.brand, d.model, d.serial
`
