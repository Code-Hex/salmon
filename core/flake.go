package core

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (salmon *Salmon) FlakeCmdNew() *cobra.Command {
	return &cobra.Command{
		Use:   "flake",
		Short: "run with cli mode",
		Long:  "run with cli mode",
		Run:   salmon.FlakeCmdRun,
	}
}

func (salmon *Salmon) FlakeCmdRun(cmd *cobra.Command, args []string) {
	fmt.Println("flake")
}
