// Package actions contains the valid actions
package actions

const (
	// Ack is the action to add an assertion to the device
	Ack = "ack"
	// Conf is the action to get a snap's configuration
	Conf = "conf"
	// Device is the action to get a device's info
	Device = "device"
	// Disable is the action for disabling a snap
	Disable = "disable"
	// Enable is the action for enabling a snap
	Enable = "enable"
	// Info is the action to get info about a snap
	Info = "info"
	// Install is the action for installing a snap
	Install = "install"
	// List is the action for getting a list of snaps
	List = "list"
	// Refresh is the action for refreshing a snap
	Refresh = "refresh"
	// Remove is the action for removing a snap
	Remove = "remove"
	// Revert is the action for reverting a snap
	Revert = "revert"
	// Server is the action to get details of the device version
	Server = "server"
	// SetConf is the action for setting snap configuration
	SetConf = "setconf"
	// Start is the action for starting a snap or snap service
	Start = "start"
	// Stop is the action for stopping a snap or snap service
	Stop = "stop"
	// Restart is the action for restarting a snap or snap service
	Restart = "restart"
	// Unregister is the action for unregistering a device from the service
	Unregister = "unregister"
)
