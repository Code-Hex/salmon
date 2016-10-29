package core

import (
	"github.com/Code-Hex/salmon/core/flake"
	"github.com/spf13/cobra"
)

type Flake struct {
	plugin  bool
	command *cobra.Command
}

func FlakeNew() *Flake {
	f := &Flake{
		command: &cobra.Command{
			Use:   "flake",
			Short: "run with cli mode",
			Long:  "run with cli mode",
		},
	}

	f.command.RunE = f.FlakeCmdRun

	// Register flags on flake sub command
	f.command.Flags().BoolVarP(&f.plugin, "plugin-register", "p", false, "register plugin")

	return f
}

func (f *Flake) FlakeCmdRun(cmd *cobra.Command, args []string) error {
	if f.plugin {
		return flake.PluginRegister()
	}
	return nil
}
