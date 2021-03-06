// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package deployer

import (
	"fmt"

	"gopkg.in/juju/names.v3"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/state"
)

// DeployerAPI provides access to the Deployer API facade.
type DeployerAPI struct {
	*common.Remover
	*common.PasswordChanger
	*common.LifeGetter
	*common.APIAddresser
	*common.UnitsWatcher
	*common.StatusSetter

	st         *state.State
	resources  facade.Resources
	authorizer facade.Authorizer
}

// NewDeployerAPI creates a new server-side DeployerAPI facade.
func NewDeployerAPI(
	st *state.State,
	resources facade.Resources,
	authorizer facade.Authorizer,
) (*DeployerAPI, error) {
	if !authorizer.AuthMachineAgent() {
		return nil, common.ErrPerm
	}
	getAuthFunc := func() (common.AuthFunc, error) {
		// Get all units of the machine and cache them.
		thisMachineTag := authorizer.GetAuthTag()
		units, err := getAllUnits(st, thisMachineTag)
		if err != nil {
			return nil, err
		}
		// Then we just check if the unit is already known.
		return func(tag names.Tag) bool {
			for _, unit := range units {
				// TODO (thumper): remove the names.Tag conversion when gccgo
				// implements concrete-type-to-interface comparison correctly.
				if names.Tag(names.NewUnitTag(unit)) == tag {
					return true
				}
			}
			return false
		}, nil
	}
	getCanWatch := func() (common.AuthFunc, error) {
		return authorizer.AuthOwner, nil
	}
	return &DeployerAPI{
		Remover:         common.NewRemover(st, true, getAuthFunc),
		PasswordChanger: common.NewPasswordChanger(st, getAuthFunc),
		LifeGetter:      common.NewLifeGetter(st, getAuthFunc),
		APIAddresser:    common.NewAPIAddresser(st, resources),
		UnitsWatcher:    common.NewUnitsWatcher(st, resources, getCanWatch),
		StatusSetter:    common.NewStatusSetter(st, getAuthFunc),
		st:              st,
		resources:       resources,
		authorizer:      authorizer,
	}, nil
}

// ConnectionInfo returns all the address information that the
// deployer task needs in one call.
func (d *DeployerAPI) ConnectionInfo() (result params.DeployerConnectionValues, err error) {
	apiAddrs, err := d.APIAddresses()
	if err != nil {
		return result, err
	}
	result = params.DeployerConnectionValues{
		APIAddresses: apiAddrs.Result,
	}
	return result, err
}

// SetStatus sets the status of the specified entities.
func (d *DeployerAPI) SetStatus(args params.SetStatus) (params.ErrorResults, error) {
	return d.StatusSetter.SetStatus(args)
}

// getAllUnits returns a list of all principal and subordinate units
// assigned to the given machine.
func getAllUnits(st *state.State, tag names.Tag) ([]string, error) {
	machine, err := st.Machine(tag.Id())
	if err != nil {
		return nil, err
	}
	// Start a watcher on machine's units, read the initial event and stop it.
	watch := machine.WatchUnits()
	defer watch.Stop()
	if units, ok := <-watch.Changes(); ok {
		return units, nil
	}
	return nil, fmt.Errorf("cannot obtain units of machine %q: %v", tag, watch.Err())
}
