package commands

import (
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"github.com/codegangsta/cli"
	"strings"
)

type ShowSpace struct {
	ui        terminal.UI
	config    *configuration.Configuration
	spaceReq requirements.SpaceRequirement
}

func NewShowSpace(ui terminal.UI, config *configuration.Configuration) (cmd ShowSpace) {
	cmd.ui = ui
	cmd.config = config
	return
}

func (cmd ShowSpace) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	cmd.spaceReq = reqFactory.NewSpaceRequirement(cmd.config.Space.Name)

	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		cmd.spaceReq,
	}
	return
}

func (cmd ShowSpace) Run(c *cli.Context) {
	space := cmd.spaceReq.GetSpace()
	cmd.ui.Say("Getting info for space %s...", terminal.EntityNameColor(space.Name))
	cmd.ui.Ok()
	cmd.ui.Say("%s:", terminal.EntityNameColor(space.Name))
	cmd.ui.Say("organization: %s", terminal.EntityNameColor(space.Organization.Name))

	domains := []string{}
	for _, domain := range space.Domains {
		domains = append(domains, domain.Name)
	}
	cmd.ui.Say("  domains: %s", terminal.EntityNameColor(strings.Join(domains, ", ")))
}
