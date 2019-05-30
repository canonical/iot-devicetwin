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

const createActionTableSQL = `
CREATE TABLE IF NOT EXISTS action (
   id             serial primary key,
   created        timestamp default current_timestamp,
   modified       timestamp default current_timestamp,
   org_id         varchar(200) not null,
   device_id      varchar(200) not null,
   action_id      varchar(200) not null,
   action         varchar(200) not null,
   status         varchar(200) default '',
   message        text default ''
)
`

const createActionSQL = `
insert into action (org_id, device_id, action_id, action, status, message)
values ($1,$2,$3,$4,$5,$6) RETURNING id`

const updateActionSQL = `
update action
set status=$2, message=$3, modified=current_timestamp
where action_id=$1`

const listActionSQL = `
select id, created, modified, org_id, device_id, action_id, action, status, message
from action
where device_id=$1
order by created`
