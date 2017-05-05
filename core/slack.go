package core

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"syscall"

	"github.com/Code-Hex/salmon/core/command"
	"github.com/Code-Hex/sigctx"
	"github.com/lestrrat/go-slack/rtm"
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
	ctx := sigctx.WithCancelSignals(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
	)
	go func() {
		if err := salmon.rtm.Run(ctx); err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case e := <-salmon.rtm.Events():
			switch typ := e.Type(); typ {
			case rtm.MessageType:
				msgev := e.Data().(*rtm.MessageEvent)
				id, msg := ToWhom(msgev.Text)
				if salmon.IsMe(ctx, id) {
					salmon.Reply(ctx, msg, msgev)
				}
			case rtm.PresenceChangeType:
				pchev := e.Data().(*rtm.PresenceChangeEvent)
				salmon.logger.Info(fmt.Sprintf("Presence Change: %v\n", pchev))
			case rtm.InvalidEventType:
				salmon.logger.Info(fmt.Sprintf("Invalid EventType\n"))
			default:
				// Ignore other events..
			}
		}
	}
}

func (salmon *Salmon) IsMe(ctx context.Context, id string) bool {
	user, err := salmon.slack.Users().Info(id).Do(ctx)
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

func (salmon *Salmon) Reply(ctx context.Context, msg string, ev *rtm.MessageEvent) {
	out, err := command.Execute(msg)
	if err != nil {
		salmon.RTMReply(ctx, err.Error(), ev.User, ev.Channel)
		return
	}
	salmon.RTMReply(ctx, "\n"+out, ev.User, ev.Channel)
}

func (salmon *Salmon) RTMReply(ctx context.Context, msg, user, channel string) {
	salmon.slack.Chat().PostMessage(channel).Username(user).Text(msg).Do(ctx)
}
