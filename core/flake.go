package core

import (
	"fmt"
	"io/ioutil"

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
	fileInfos, err := ioutil.ReadDir("plugin")
	if err != nil {
		return errors.Wrapf(err, "Could not open plugin directory")
	}

	for _, fi := range fileInfos {
		fmt.Println(fi.Name())
	}

	return nil
}

func (flake *Flake) grepRunInPlugins() {
	// ここに open して Run[A-Za-z]+()を探して,
	// map へ "[a-z]+":"Run[A-Za-z]+" を保存する
}
