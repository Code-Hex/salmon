package core

import (
	"fmt"
	"os"
	"time"

	"github.com/nlopes/slack"
	"github.com/spf13/cobra"
	"github.com/uber-go/zap"
)

const (
	version = "0.0.1"
	msg     = "salmon v" + version + ", salmon is .* bot\n"
)

type Salmon struct {
	trace   bool
	version bool
	args    []string
	swim    *cobra.Command
	slack   *slack.Client
	rtm     *slack.RTM
	logger  zap.Logger
}

func (salmon *Salmon) RootCmdNew() *cobra.Command {
	return &cobra.Command{
		Use:   "salmon",
		Short: msg,
		Long:  msg,
		Run:   salmon.Dive,
	}
}

func Generate(Out zap.WriteSyncer) *Salmon {
	slack := slack.New(os.Getenv("SLACK_TOKEN"))
	salmon := &Salmon{
		slack: slack,
		rtm:   slack.NewRTM(),
		logger: zap.New(
			zap.NewTextEncoder(zap.TextTimeFormat(time.ANSIC)),
			zap.AddCaller(), // Add line number option
			zap.Output(Out),
		),
	}

	salmon.swim = salmon.RootCmdNew()

	// Register sub command
	salmon.swim.AddCommand(salmon.FlakeCmdNew())
	salmon.swim.AddCommand(salmon.SlackCmdNew())

	salmon.swim.Flags().BoolVarP(&salmon.trace, "trace", "t", false, "display detail error messages")
	salmon.swim.Flags().BoolVarP(&salmon.version, "version", "v", false, "display the version of salmon and exit")
	return salmon
}

func (salmon *Salmon) Swim() int {
	if err := salmon.swim.Execute(); err != nil {
		if salmon.trace {
			fmt.Fprintf(os.Stderr, "Error:\n%+v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Error:\n  %v\n", err)
		}
		return 1
	}
	return 0
}

func (salmon *Salmon) Dive(cmd *cobra.Command, args []string) {
	if salmon.version {
		os.Stdout.Write([]byte(msg))
	} else {
		fmt.Println("salmon")
	}

}
