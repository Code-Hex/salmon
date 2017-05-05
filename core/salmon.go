package core

import (
	"fmt"
	"os"

	slack "github.com/lestrrat/go-slack"
	rtm "github.com/lestrrat/go-slack/rtm"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	version = "0.0.2"
	msg     = "salmon v" + version + ", salmon is .* bot\n"
)

type Salmon struct {
	trace   bool
	version bool
	args    []string
	swim    *cobra.Command
	slack   *slack.Client
	flake   *Flake
	rtm     *rtm.Client
	logger  *zap.Logger
}

func (salmon *Salmon) RootCmdNew() *cobra.Command {
	return &cobra.Command{
		Use:           "salmon",
		Short:         msg,
		Long:          msg,
		RunE:          salmon.Dive,
		SilenceErrors: true,
	}
}

func Generate() *Salmon {
	slack := slack.New(os.Getenv("SLACK_TOKEN"))
	salmon := &Salmon{
		slack: slack,
		flake: FlakeNew(),
		rtm:   rtm.New(slack),
	}

	switch os.Getenv("STAGE") {
	case "production":
		logger, err := zap.NewProduction(
			zap.AddCaller(),
		)
		if err != nil {
			panic("Failed to set logger")
		}
		salmon.logger = logger
	case "development":
		logger, err := zap.NewDevelopment(
			zap.AddCaller(),
		)
		if err != nil {
			panic("Failed to set logger")
		}
		salmon.logger = logger
	default:
		fmt.Fprintln(os.Stderr, `Do not specified "STAGE" environment variable`)
		os.Exit(1)
	}

	salmon.swim = salmon.RootCmdNew()

	// Register sub command
	salmon.swim.AddCommand(salmon.flake.command)
	salmon.swim.AddCommand(salmon.SlackCmdNew())

	// Register flags on root command
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

func (salmon *Salmon) Dive(cmd *cobra.Command, args []string) error {
	if salmon.version {
		os.Stdout.Write([]byte(msg))
	}

	return cmd.Usage()
}
