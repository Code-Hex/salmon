package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Code-Hex/salmon/core/command"
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
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				id, msg := ToWhom(ev.Text)
				if salmon.IsMe(id) {
					salmon.Reply(msg, ev)
				}
			case *slack.PresenceChangeEvent:
				salmon.logger.Info(fmt.Sprintf("Presence Change: %v\n", ev))

			case *slack.RTMError:
				salmon.logger.Error(fmt.Sprintf("Error: %s\n", ev.Error()))
				return
			default:
				// Ignore other events..
			}
		}
	}
}

func (salmon *Salmon) IsMe(id string) bool {
	user, err := salmon.slack.GetUserInfo(id)
	if err != nil {
		salmon.logger.Error(err.Error())
		return false
	}
	return user.Name == "salmon"
}

func ToWhom(s string) (string, string) {
	var flag = false
	var index = 0
	var buf bytes.Buffer

	ch := []rune(s)
	len := len(ch)

	for i := 0; i < len; i++ {
		if ch[i] == '<' {
			if i+1 < len && ch[i+1] == '@' {
				i++
				flag = true
				continue
			}
		}

		if ch[i] == '>' {
			index = i + 1
			break
		}

		if flag {
			buf.WriteRune(ch[i])
		}
	}

	return buf.String(), strings.TrimSpace(string(ch[index:]))
}

func (salmon *Salmon) Reply(msg string, ev *slack.MessageEvent) {
	out, err := command.Execute(msg)
	if err != nil {
		salmon.RTMReply(err.Error(), ev.Msg.User, ev.Msg.Channel)
		return
	}
	salmon.RTMReply("\n"+out, ev.Msg.User, ev.Msg.Channel)
}

func (salmon *Salmon) RTMReply(msg, user, channel string) {
	salmon.rtm.SendMessage(salmon.rtm.NewOutgoingMessage(fmt.Sprintf("<@%s> %s", user, msg), channel))
}
