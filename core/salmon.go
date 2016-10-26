package core

import (
	"os"
	"time"

	"github.com/nlopes/slack"
	"github.com/uber-go/zap"
)

type Salmon struct {
	slack  *slack.Client
	rtm    *slack.RTM
	logger zap.Logger
}

func (salmon *Salmon) Swim() int {
	salmon.connectRTM()
	return 0
}

func Generate(Out zap.WriteSyncer) *Salmon {
	slack := slack.New(os.Getenv("SLACK_TOKEN"))
	return &Salmon{
		slack: slack,
		rtm:   slack.NewRTM(),
		logger: zap.New(
			zap.NewTextEncoder(zap.TextTimeFormat(time.ANSIC)),
			zap.AddCaller(), // Add line number option
			zap.Output(Out),
		),
	}
}
