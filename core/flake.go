package core

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Flake struct {
	plugin  bool
	command *cobra.Command
}

func FlakeNew() *Flake {
	flake := &Flake{
		command: &cobra.Command{
			Use:   "flake",
			Short: "run with cli mode",
			Long:  "run with cli mode",
		},
	}

	flake.command.RunE = flake.FlakeCmdRun

	// Register flags on flake sub command
	flake.command.Flags().BoolVarP(&flake.plugin, "plugin-register", "", false, "register plugin")

	return flake
}

func (flake *Flake) FlakeCmdRun(cmd *cobra.Command, args []string) error {
	if flake.plugin {
		return flake.pluginRegister()
	}
	return nil
}

func (flake *Flake) pluginRegister() error {

	dir, err := os.Open("plugin")
	if err != nil {
		return errors.Wrapf(err, "Could not open plugin directory")
	}
	defer dir.Close()

	for {
		entries, err := dir.Readdir(256)
		if err == io.EOF {
			break
		}

		for _, fi := range entries {
			fmt.Println(fi.Name())
		}
	}

	return nil
}
