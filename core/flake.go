package core

import (
	"github.com/Code-Hex/salmon/core/flake"
	"github.com/spf13/cobra"
)

type Flake struct {
	ed      bool   `ExecDebug`
	pr      bool   `PluginRegistor`
	pm      string `PluginMaker`
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
	f.command.Flags().BoolVarP(&f.ed, "exec-debug", "e", false, "command exec on cli")
	f.command.Flags().BoolVarP(&f.pr, "plugin-register", "r", false, "register plugin")
	f.command.Flags().StringVarP(&f.pm, "plugin-maker", "m", "", "plugin maker")

	return f
}

func (f *Flake) FlakeCmdRun(cmd *cobra.Command, args []string) error {

	if f.ed {
		return flake.ExecDebug(args)
	}

	if f.pm != "" {
		return flake.PluginMaker(f.pm)
	}

	if f.pr {
		return flake.PluginRegister()
	}

	return cmd.Usage()
}
