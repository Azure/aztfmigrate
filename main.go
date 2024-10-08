package main

import (
	"os"

	"github.com/Azure/aztfmigrate/cmd"
	"github.com/mitchellh/cli"
)

func main() {
	c := &cli.CLI{
		Name:       "aztfmigrate",
		Version:    version,
		Args:       os.Args[1:],
		HelpWriter: os.Stdout,
	}

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			Reader:      os.Stdin,
			ErrorWriter: os.Stderr,
		},
	}

	c.Commands = map[string]cli.CommandFactory{
		"migrate": func() (cli.Command, error) {
			return &cmd.MigrateCommand{
				Ui: ui,
			}, nil
		},
		"plan": func() (cli.Command, error) {
			return &cmd.PlanCommand{
				Ui: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &cmd.VersionCommand{
				Ui:      ui,
				Version: version,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error("Error: " + err.Error())
	}

	os.Exit(exitStatus)
}
