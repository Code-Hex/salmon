package core

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/nlopes/slack"
	"github.com/spf13/cobra"
)

func (salmon *Salmon) SlackCmdNew() *cobra.Command {
	return &cobra.Command{
		Use:   "slack",
		Short: "run with slack bot mode",
		Long:  "run with slack bot mode",
		Run:   salmon.SlackCmdRun,
	}
}

func (salmon *Salmon) SlackCmdRun(cmd *cobra.Command, args []string) {
	salmon.connectRTM()
}

func (salmon *Salmon) connectRTM() {
	go salmon.rtm.ManageConnection()

	for {
		select {
		case msg := <-salmon.rtm.IncomingEvents:
			salmon.logger.Info("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				// fmt.Println("Infos:", zap.String(ev.Info))
				// salmon.logger.Info("Connection counter:", ev.ConnectionCount)
				// Replace #general with your Channel ID
				salmon.rtm.SendMessage(salmon.rtm.NewOutgoingMessage("Hello world", "#general"))

			case *slack.MessageEvent:
				//Parse(ev.Text)
				pp.Print(ev.Text)
				salmon.Reply("Hi!!", ev.Msg.User, ev.Msg.Channel)
				salmon.logger.Info(fmt.Sprintf("Message: %v\n", ev))

			case *slack.PresenceChangeEvent:
				salmon.logger.Info(fmt.Sprintf("Presence Change: %v\n", ev))

			case *slack.LatencyReport:
				salmon.logger.Info(fmt.Sprintf("Current latency: %v\n", ev.Value))

			case *slack.RTMError:
				salmon.logger.Info(fmt.Sprintf("Error: %s\n", ev.Error()))

			case *slack.FileSharedEvent:
				salmon.logger.Info(fmt.Sprintf("File: %v\n", ev))

			case *slack.InvalidAuthEvent:
				salmon.logger.Info("Invalid credentials")
				return
			default:

				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}

func (salmon *Salmon) Reply(msg, user, channel string) {
	salmon.rtm.SendMessage(salmon.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> %s", user, msg), channel))
}
