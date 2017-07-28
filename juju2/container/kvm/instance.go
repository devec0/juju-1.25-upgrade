// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package kvm

import (
	"fmt"

	"github.com/juju/1.25-upgrade/juju2/instance"
	"github.com/juju/1.25-upgrade/juju2/network"
	"github.com/juju/1.25-upgrade/juju2/status"
)

type kvmInstance struct {
	container Container
	id        string
}

var _ instance.Instance = (*kvmInstance)(nil)

// Id implements instance.Instance.Id.
func (kvm *kvmInstance) Id() instance.Id {
	return instance.Id(kvm.id)
}

// Status implements instance.Instance.Status.
func (kvm *kvmInstance) Status() instance.InstanceStatus {
	if kvm.container.IsRunning() {
		return instance.InstanceStatus{
			Status:  status.Running,
			Message: "running",
		}
	}
	return instance.InstanceStatus{
		Status:  status.Stopped,
		Message: "stopped",
	}
}

func (*kvmInstance) Refresh() error {
	return nil
}

func (kvm *kvmInstance) Addresses() ([]network.Address, error) {
	logger.Errorf("kvmInstance.Addresses not implemented")
	return nil, nil
}

// OpenPorts implements instance.Instance.OpenPorts.
func (kvm *kvmInstance) OpenPorts(machineId string, rules []network.IngressRule) error {
	return fmt.Errorf("not implemented")
}

// ClosePorts implements instance.Instance.ClosePorts.
func (kvm *kvmInstance) ClosePorts(machineId string, rules []network.IngressRule) error {
	return fmt.Errorf("not implemented")
}

// IngressRules implements instance.Instance.IngressRules.
func (kvm *kvmInstance) IngressRules(machineId string) ([]network.IngressRule, error) {
	return nil, fmt.Errorf("not implemented")
}

// Add a string representation of the id.
func (kvm *kvmInstance) String() string {
	return fmt.Sprintf("kvm:%s", kvm.id)
}
