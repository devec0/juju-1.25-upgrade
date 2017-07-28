// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package controller

import (
	"time"

	"github.com/juju/cmd"
	"github.com/juju/utils/clock"

	"github.com/juju/1.25-upgrade/juju2/api"
	"github.com/juju/1.25-upgrade/juju2/api/base"
	"github.com/juju/1.25-upgrade/juju2/cmd/modelcmd"
	"github.com/juju/1.25-upgrade/juju2/environs"
	"github.com/juju/1.25-upgrade/juju2/jujuclient"
)

// NewListControllersCommandForTest returns a listControllersCommand with the clientstore provided
// as specified.
func NewListControllersCommandForTest(testStore jujuclient.ClientStore, api func(string) ControllerAccessAPI) *listControllersCommand {
	return &listControllersCommand{
		store: testStore,
		api:   api,
	}
}

// NewShowControllerCommandForTest returns a showControllerCommand with the clientstore provided
// as specified.
func NewShowControllerCommandForTest(testStore jujuclient.ClientStore, api func(string) ControllerAccessAPI) *showControllerCommand {
	return &showControllerCommand{
		store: testStore,
		api:   api,
	}
}

type AddModelCommand struct {
	*addModelCommand
}

// NewAddModelCommandForTest returns a AddModelCommand with
// the api provided as specified.
func NewAddModelCommandForTest(
	apiRoot api.Connection,
	api AddModelAPI,
	cloudAPI CloudAPI,
	store jujuclient.ClientStore,
	providerRegistry environs.ProviderRegistry,
) (cmd.Command, *AddModelCommand) {
	c := &addModelCommand{
		apiRoot: apiRoot,
		newAddModelAPI: func(caller base.APICallCloser) AddModelAPI {
			return api
		},
		newCloudAPI: func(base.APICallCloser) CloudAPI {
			return cloudAPI
		},
		providerRegistry: providerRegistry,
	}
	c.SetClientStore(store)
	return modelcmd.WrapController(c), &AddModelCommand{c}
}

// NewListModelsCommandForTest returns a ListModelsCommand with the API
// and userCreds provided as specified.
func NewListModelsCommandForTest(modelAPI ModelManagerAPI, sysAPI ModelsSysAPI, store jujuclient.ClientStore) cmd.Command {
	c := &modelsCommand{
		modelAPI: modelAPI,
		sysAPI:   sysAPI,
	}
	c.SetClientStore(store)
	return modelcmd.WrapController(c)
}

// NewRegisterCommandForTest returns a RegisterCommand with the function used
// to open the API connection mocked out.
func NewRegisterCommandForTest(apiOpen api.OpenFunc, listModels func(jujuclient.ClientStore, string, string) ([]base.UserModel, error), store jujuclient.ClientStore) modelcmd.Command {
	return modelcmd.WrapBase(&registerCommand{
		apiOpen:        apiOpen,
		listModelsFunc: listModels,
		store:          store,
	})
}

// NewEnableDestroyControllerCommandForTest returns a enableDestroyController with the
// function used to open the API connection mocked out.
func NewEnableDestroyControllerCommandForTest(api removeBlocksAPI, store jujuclient.ClientStore) cmd.Command {
	c := &enableDestroyController{
		api: api,
	}
	c.SetClientStore(store)
	return modelcmd.WrapController(c)
}

// NewDestroyCommandForTest returns a DestroyCommand with the controller and
// client endpoints mocked out.
func NewDestroyCommandForTest(
	api destroyControllerAPI,
	clientapi destroyClientAPI,
	store jujuclient.ClientStore,
	apierr error,
) cmd.Command {
	cmd := &destroyCommand{
		destroyCommandBase: destroyCommandBase{
			api:       api,
			clientapi: clientapi,
			apierr:    apierr,
		},
	}
	cmd.SetClientStore(store)
	return modelcmd.WrapController(
		cmd,
		modelcmd.WrapControllerSkipControllerFlags,
		modelcmd.WrapControllerSkipDefaultController,
	)
}

// NewKillCommandForTest returns a killCommand with the controller and client
// endpoints mocked out.
func NewKillCommandForTest(
	api destroyControllerAPI,
	clientapi destroyClientAPI,
	store jujuclient.ClientStore,
	apierr error,
	clock clock.Clock,
	apiOpen api.OpenFunc,
) cmd.Command {
	kill := &killCommand{
		destroyCommandBase: destroyCommandBase{
			api:       api,
			clientapi: clientapi,
			apierr:    apierr,
		},
		clock: clock,
	}
	kill.SetClientStore(store)
	wrapped := wrapKillCommand(kill)
	if apiOpen != nil {
		wrapped.SetAPIOpen(apiOpen)
	}
	return wrapped
}

// KillTimeout returns the internal timeout duration of the kill command.
func KillTimeout(command cmd.Command) time.Duration {
	return modelcmd.InnerCommand(command).(*killCommand).timeout
}

// KillWaitForModels calls the WaitForModels method of the kill command.
func KillWaitForModels(command cmd.Command, ctx *cmd.Context, api destroyControllerAPI, uuid string) error {
	return modelcmd.InnerCommand(command).(*killCommand).WaitForModels(ctx, api, uuid)
}

// NewGetConfigCommandCommandForTest returns a GetConfigCommandCommand with
// the api provided as specified.
func NewGetConfigCommandForTest(api controllerAPI, store jujuclient.ClientStore) cmd.Command {
	c := &getConfigCommand{api: api}
	c.SetClientStore(store)
	return modelcmd.WrapController(c)
}

type CtrData ctrData
type ModelData modelData

func FmtCtrStatus(data CtrData) string {
	return fmtCtrStatus(ctrData(data))
}

func FmtModelStatus(data ModelData) string {
	return fmtModelStatus(modelData(data))
}

func NewData(api destroyControllerAPI, ctrUUID string) (ctrData, []modelData, error) {
	return newData(api, ctrUUID)
}
