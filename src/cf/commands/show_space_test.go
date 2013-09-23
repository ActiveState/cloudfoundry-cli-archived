package commands_test

import (
	"cf"
	. "cf/commands"
	"cf/configuration"
	"github.com/stretchr/testify/assert"
	"testhelpers"
	"testing"
)

func TestShowSpaceRequirements(t *testing.T) {
	config := &configuration.Configuration{
		Space:        cf.Space{},
		Organization: cf.Organization{},
	}

	reqFactory := &testhelpers.FakeReqFactory{LoginSuccess: true}
	callShowSpace([]string{}, reqFactory, config)
	assert.True(t, testhelpers.CommandDidPassRequirements)

	reqFactory = &testhelpers.FakeReqFactory{LoginSuccess: false}
	callShowSpace([]string{}, reqFactory, config)
	assert.False(t, testhelpers.CommandDidPassRequirements)
}

func TestShowSpaceInfoSuccess(t *testing.T) {
	space := cf.Space{Name: "space1"}
	org := cf.Organization{Name: "org1"}
	config := &configuration.Configuration{
		Space:        space,
		Organization: org,
	}

	reqFactory := &testhelpers.FakeReqFactory{LoginSuccess: true}
	ui := callShowSpace([]string{}, reqFactory, config)
	assert.Contains(t, ui.Outputs[0], "Getting info for space")
	assert.Contains(t, ui.Outputs[0], "space1")
	assert.Contains(t, ui.Outputs[1], "OK")
	assert.Contains(t, ui.Outputs[2], "space1:")
	assert.Contains(t, ui.Outputs[3], "organization")
	assert.Contains(t, ui.Outputs[3], "org1")
	assert.Contains(t, ui.Outputs[3], "apps")
	assert.Contains(t, ui.Outputs[3], "app1")
	assert.Contains(t, ui.Outputs[3], "domains")
	assert.Contains(t, ui.Outputs[3], "domain1")
	assert.Contains(t, ui.Outputs[3], "services")
	assert.Contains(t, ui.Outputs[3], "service1")
}

func callShowSpace(args []string, reqFactory *testhelpers.FakeReqFactory, config *configuration.Configuration) (ui *testhelpers.FakeUI) {
	ui = new(testhelpers.FakeUI)
	ctxt := testhelpers.NewContext("space", args)

	cmd := NewShowSpace(ui, config)
	testhelpers.RunCommand(cmd, ctxt, reqFactory)
	return
}
