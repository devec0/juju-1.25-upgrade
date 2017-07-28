// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package imagemanager

import (
	names "gopkg.in/juju/names.v2"

	"github.com/juju/1.25-upgrade/juju2/state"
	"github.com/juju/1.25-upgrade/juju2/state/imagestorage"
)

type stateInterface interface {
	ImageStorage() imagestorage.Storage
	ControllerTag() names.ControllerTag
}

type stateShim struct {
	*state.State
}

func (s stateShim) ImageStorage() imagestorage.Storage {
	return s.State.ImageStorage()
}
