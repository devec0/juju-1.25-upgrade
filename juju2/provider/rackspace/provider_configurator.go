// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package rackspace

import (
	"github.com/juju/errors"
	"github.com/juju/schema"
	"gopkg.in/goose.v2/nova"

	"github.com/juju/1.25-upgrade/juju2/cloudconfig/cloudinit"
	"github.com/juju/1.25-upgrade/juju2/environs"
)

type rackspaceConfigurator struct {
}

// ModifyRunServerOptions implements ProviderConfigurator interface.
func (c *rackspaceConfigurator) ModifyRunServerOptions(options *nova.RunServerOpts) {
	// More on how ConfigDrive option is used on rackspace:
	// http://docs.rackspace.com/servers/api/v2/cs-devguide/content/config_drive_ext.html
	options.ConfigDrive = true
}

// GetCloudConfig implements ProviderConfigurator interface.
func (c *rackspaceConfigurator) GetCloudConfig(args environs.StartInstanceParams) (cloudinit.CloudConfig, error) {
	cloudcfg, err := cloudinit.New(args.Tools.OneSeries())
	if err != nil {
		return nil, errors.Trace(err)
	}
	// Additional package required for sshInstanceConfigurator, to save
	// iptables state between restarts.
	cloudcfg.AddPackage("iptables-persistent")
	return cloudcfg, nil
}

// GetConfigDefaults implements ProviderConfigurator interface.
func (c *rackspaceConfigurator) GetConfigDefaults() schema.Defaults {
	return schema.Defaults{
		"use-floating-ip":      false,
		"use-default-secgroup": false,
		"network":              "",
		"external-network":     "",
	}
}
