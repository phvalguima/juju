// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package sshclient_test

import (
	"github.com/juju/errors"
	"github.com/juju/names"
	jujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/apiserver/sshclient"
	apiservertesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/network"
	"github.com/juju/juju/state"
	"github.com/juju/juju/testing"
)

type facadeSuite struct {
	testing.BaseSuite
	backend    *mockBackend
	authorizer *apiservertesting.FakeAuthorizer
	facade     *sshclient.Facade
}

var _ = gc.Suite(&facadeSuite{})

func (s *facadeSuite) SetUpTest(c *gc.C) {
	s.backend = new(mockBackend)
	s.authorizer = new(apiservertesting.FakeAuthorizer)
	s.authorizer.Tag = names.NewUserTag("igor")
	facade, err := sshclient.New(s.backend, nil, s.authorizer)
	c.Assert(err, jc.ErrorIsNil)
	s.facade = facade
}

func (s *facadeSuite) TestFailsWhenMachine(c *gc.C) {
	s.authorizer.Tag = names.NewMachineTag("0")
	_, err := sshclient.New(s.backend, nil, s.authorizer)
	c.Assert(err, gc.Equals, common.ErrPerm)
}

func (s *facadeSuite) TestFailsWhenUnit(c *gc.C) {
	s.authorizer.Tag = names.NewUnitTag("foo/0")
	_, err := sshclient.New(s.backend, nil, s.authorizer)
	c.Assert(err, gc.Equals, common.ErrPerm)
}

func (s *facadeSuite) TestPublicAddress(c *gc.C) {
	args := params.SSHTargets{
		Targets: []params.SSHTarget{
			{Target: "0"},
			{Target: "foo/0"},
			{Target: "other/1"},
		},
	}
	results, err := s.facade.PublicAddress(args)

	c.Assert(err, jc.ErrorIsNil)
	c.Check(results, gc.DeepEquals, params.SSHAddressResults{
		Results: []params.SSHAddressResult{
			{Address: "1.1.1.1"},
			{Address: "3.3.3.3"},
			{Error: apiservertesting.ErrUnauthorized},
		},
	})
	s.backend.stub.CheckCalls(c, []jujutesting.StubCall{
		{"GetMachineForTarget", []interface{}{"0"}},
		{"GetMachineForTarget", []interface{}{"foo/0"}},
		{"GetMachineForTarget", []interface{}{"other/1"}},
	})
}

func (s *facadeSuite) TestPrivateAddress(c *gc.C) {
	args := params.SSHTargets{
		Targets: []params.SSHTarget{
			{Target: "other/1"},
			{Target: "0"},
			{Target: "foo/0"},
		},
	}
	results, err := s.facade.PrivateAddress(args)

	c.Assert(err, jc.ErrorIsNil)
	c.Check(results, gc.DeepEquals, params.SSHAddressResults{
		Results: []params.SSHAddressResult{
			{Error: apiservertesting.ErrUnauthorized},
			{Address: "2.2.2.2"},
			{Address: "4.4.4.4"},
		},
	})
	s.backend.stub.CheckCalls(c, []jujutesting.StubCall{
		{"GetMachineForTarget", []interface{}{"other/1"}},
		{"GetMachineForTarget", []interface{}{"0"}},
		{"GetMachineForTarget", []interface{}{"foo/0"}},
	})
}

func (s *facadeSuite) TestPublicKeys(c *gc.C) {
	args := params.SSHTargets{
		Targets: []params.SSHTarget{
			{Target: "0"},
			{Target: "other/1"},
			{Target: "foo/0"},
		},
	}
	results, err := s.facade.PublicKeys(args)

	c.Assert(err, jc.ErrorIsNil)
	c.Check(results, gc.DeepEquals, params.SSHPublicKeysResults{
		Results: []params.SSHPublicKeysResult{
			{PublicKeys: []string{"rsa0", "dsa0"}},
			{Error: apiservertesting.ErrUnauthorized},
			{PublicKeys: []string{"rsa1", "dsa1"}},
		},
	})
	s.backend.stub.CheckCalls(c, []jujutesting.StubCall{
		{"GetMachineForTarget", []interface{}{"0"}},
		{"GetSSHHostKeys", []interface{}{names.NewMachineTag("0")}},
		{"GetMachineForTarget", []interface{}{"other/1"}},
		{"GetMachineForTarget", []interface{}{"foo/0"}},
		{"GetSSHHostKeys", []interface{}{names.NewMachineTag("1")}},
	})
}

func (s *facadeSuite) TestProxyTrue(c *gc.C) {
	s.backend.proxySSH = true
	result, err := s.facade.Proxy()
	c.Assert(err, jc.ErrorIsNil)
	c.Check(result.UseProxy, jc.IsTrue)
	s.backend.stub.CheckCalls(c, []jujutesting.StubCall{
		{"ModelConfig", []interface{}{}},
	})
}

func (s *facadeSuite) TestProxyFalse(c *gc.C) {
	s.backend.proxySSH = false
	result, err := s.facade.Proxy()
	c.Assert(err, jc.ErrorIsNil)
	c.Check(result.UseProxy, jc.IsFalse)
	s.backend.stub.CheckCalls(c, []jujutesting.StubCall{
		{"ModelConfig", []interface{}{}},
	})
}

type mockBackend struct {
	stub     jujutesting.Stub
	proxySSH bool
}

func (backend *mockBackend) ModelConfig() (*config.Config, error) {
	backend.stub.AddCall("ModelConfig")
	attrs := testing.FakeConfig()
	attrs["proxy-ssh"] = backend.proxySSH
	conf, err := config.New(config.NoDefaults, attrs)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return conf, nil
}

func (backend *mockBackend) GetMachineForTarget(target string) (sshclient.SSHMachine, error) {
	backend.stub.AddCall("GetMachineForTarget", target)
	switch target {
	case "0":
		return &mockMachine{
			tag:            names.NewMachineTag("0"),
			publicAddress:  "1.1.1.1",
			privateAddress: "2.2.2.2",
		}, nil
	case "foo/0":
		return &mockMachine{
			tag:            names.NewMachineTag("1"),
			publicAddress:  "3.3.3.3",
			privateAddress: "4.4.4.4",
		}, nil
	}
	return nil, errors.New("unknown target")
}

func (backend *mockBackend) GetSSHHostKeys(tag names.MachineTag) (state.SSHHostKeys, error) {
	backend.stub.AddCall("GetSSHHostKeys", tag)
	switch tag {
	case names.NewMachineTag("0"):
		return state.SSHHostKeys{"rsa0", "dsa0"}, nil
	case names.NewMachineTag("1"):
		return state.SSHHostKeys{"rsa1", "dsa1"}, nil
	}
	return nil, errors.New("machine not found")
}

type mockMachine struct {
	tag            names.MachineTag
	publicAddress  string
	privateAddress string
}

func (m *mockMachine) MachineTag() names.MachineTag {
	return m.tag
}

func (m *mockMachine) PublicAddress() (network.Address, error) {
	return network.Address{
		Value: m.publicAddress,
	}, nil
}

func (m *mockMachine) PrivateAddress() (network.Address, error) {
	return network.Address{
		Value: m.privateAddress,
	}, nil
}
