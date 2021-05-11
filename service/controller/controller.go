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

package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/everactive/iot-devicetwin/pkg/actions"

	"github.com/everactive/iot-devicetwin/pkg/messages"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/everactive/iot-devicetwin/config"
	"github.com/everactive/iot-devicetwin/domain"
	"github.com/everactive/iot-devicetwin/service/devicetwin"
	"github.com/everactive/iot-devicetwin/service/mqtt"
	"github.com/segmentio/ksuid"
)

// Controller interface for the service
type Controller interface {
	// MQTT handlers
	HealthHandler(client MQTT.Client, msg MQTT.Message)
	ActionHandler(client MQTT.Client, msg MQTT.Message)

	// Passthrough to the device twin service
	DeviceSnaps(orgID, clientID string) ([]messages.DeviceSnap, error)
	DeviceList(orgID string) ([]messages.Device, error)
	DeviceDelete(deviceID string) error
	DeviceGet(orgID, clientID string) (messages.Device, error)
	DeviceUnregister(orgID, clientID string) error
	GroupCreate(orgID, name string) error
	GroupList(orgID string) ([]domain.Group, error)
	GroupGet(orgID, name string) (domain.Group, error)
	GroupLinkDevice(orgID, name, clientID string) error
	GroupUnlinkDevice(orgID, name, clientID string) error
	GroupGetDevices(orgID, name string) ([]messages.Device, error)
	GroupGetExcludedDevices(orgID, name string) ([]messages.Device, error)

	// Actions on a device
	DeviceSnapList(orgID, clientID string) error
	DeviceSnapInstall(orgID, clientID, snap string) error
	DeviceSnapServiceAction(orgID, clientID, snap, action string, services *messages.SnapServiceMessage) error
	DeviceSnapRemove(orgID, clientID, snap string) error
	DeviceSnapUpdate(orgID, clientID, snap, action string) error
	DeviceSnapConf(orgID, clientID, snap, settings string) error
	ActionList(orgID, clientID string) ([]domain.Action, error)
}

const (
	clientIDMQTTTopicPartsCount = 4
)

// Service implementation of the devicetwin service use cases
type Service struct {
	Settings   *config.Settings
	MQTT       mqtt.Connect
	DeviceTwin devicetwin.DeviceTwin
}

// NewService creates an implementation of the devicetwin use cases
func NewService(settings *config.Settings, m mqtt.Connect, twin devicetwin.DeviceTwin) *Service {
	srv := &Service{
		Settings:   settings,
		MQTT:       m,
		DeviceTwin: twin,
	}

	// Setup the MQTT client and handle pub/sub from here... as the MQTT and DeviceTwin services are mutually dependent
	// This service plugs them together
	if err := srv.SubscribeToActions(); err != nil {
		log.Printf("Error subscribing to actions: %v", err)
	}
	return srv
}

// SubscribeToActions subscribes to the published topics from the devices
func (srv *Service) SubscribeToActions() error {
	const (
		topicActions = "devices/pub/+"
		topicHealth  = "devices/health/+"
	)

	// Subscribe to the device health messages
	if err := srv.MQTT.Subscribe(topicHealth, srv.HealthHandler); err != nil {
		log.Printf("Error subscribing to topic `%s`: %v", topicHealth, err)
		return err
	}

	// Subscribe to the device action responses
	if err := srv.MQTT.Subscribe(topicActions, srv.ActionHandler); err != nil {
		log.Printf("Error subscribing to topic `%s`: %v", topicActions, err)
		return err
	}

	return nil
}

// ActionHandler is the handler for the main subscription topic
func (srv *Service) ActionHandler(client MQTT.Client, msg MQTT.Message) {
	clientID := getClientID(msg)
	log.Printf("Action response from %s", clientID)

	// Parse the body
	a := messages.PublishResponse{}
	if err := json.Unmarshal(msg.Payload(), &a); err != nil {
		log.Printf("Error in action message: %v", err)
		return
	}

	// Check if there is an error and handle it
	if !a.Success {
		log.Printf("Error in action `%s`: (%s) %s", a.Action, a.Id, a.Message)
		return
	}

	// Handle the action
	if err := srv.DeviceTwin.ActionResponse(clientID, a.Id, a.Action, msg.Payload()); err != nil {
		log.Printf("Error with action `%s`: %v", a.Action, err)
	}
}

// HealthHandler is the handler for the devices health messages
func (srv *Service) HealthHandler(client MQTT.Client, msg MQTT.Message) {
	clientID := getClientID(msg)
	log.Printf("Health update from %s", clientID)

	// Parse the body
	h := messages.Health{}
	if err := json.Unmarshal(msg.Payload(), &h); err != nil {
		log.Printf("Error in health message: %v", err)
		return
	}

	// Check that the client ID matches
	if clientID != h.DeviceId {
		log.Printf("Client/device ID mismatch: %s and %s", clientID, h.DeviceId)
		return
	}

	// Update the device record
	if err := srv.DeviceTwin.HealthHandler(h); err == nil {
		// Exit if successful
		return
	}

	// We don't have the device details, so request them from the device
	act := messages.SubscribeAction{
		Action: actions.Device,
	}
	if err := srv.triggerActionOnDevice(h.OrgId, h.DeviceId, act); err != nil {
		log.Printf("Triggering action: %v", err)
	}

	// Get the snaps from the device
	act = messages.SubscribeAction{
		Action: actions.List,
	}
	if err := srv.triggerActionOnDevice(h.OrgId, h.DeviceId, act); err != nil {
		log.Printf("Triggering action: %v", err)
	}
}

// getClientID sets the client ID from the topic
func getClientID(msg MQTT.Message) string {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) != clientIDMQTTTopicPartsCount-1 {
		log.Printf("Error in health message: expected 4 parts, got %d", len(parts))
		return ""
	}
	return parts[2]
}

// triggerActionOnDevice triggers an action on the device via MQTT
func (srv *Service) triggerActionOnDevice(orgID, deviceID string, act messages.SubscribeAction) error {
	// Generate a request ID
	id := ksuid.New()
	act.Id = id.String()

	// Serialize the action
	data, err := serializePayload(act)
	if err != nil {
		log.Printf("Error in action serialization: %v", err)
		return err
	}

	// Publish the request
	t := fmt.Sprintf("devices/sub/%s", deviceID)
	err = srv.MQTT.Publish(t, string(data))
	if err != nil {
		log.Printf("Error in publish: %v", err)
		return fmt.Errorf("error in publish: %v", err)
	}

	// Log the request
	return srv.DeviceTwin.ActionCreate(orgID, deviceID, act)
}

func serializePayload(act messages.SubscribeAction) ([]byte, error) {
	return json.Marshal(act)
}
