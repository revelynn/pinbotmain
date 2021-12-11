package pinbot

import (
	"github.com/elliotwms/pinbot/internal/eventhandlers"
)

const logFieldHandler = "handler"

func (bot *Bot) registerHandlers() func() {
	closers := []func(){
		bot.session.AddHandler(eventhandlers.Ready(bot.log.WithField(logFieldHandler, "Ready"))),
		bot.session.AddHandler(eventhandlers.MessageReactionAdd(bot.log.WithField(logFieldHandler, "MessageReactionAdd"), bot.testGuildID)),
		bot.session.AddHandler(eventhandlers.GuildCreate(bot.log.WithField(logFieldHandler, "GuildCreate"))),
		bot.session.AddHandler(eventhandlers.ChannelCreate(bot.log.WithField(logFieldHandler, "ChannelCreate"))),
		bot.session.AddHandler(eventhandlers.ChannelUpdate(bot.log.WithField(logFieldHandler, "ChannelUpdate"))),
		bot.session.AddHandler(eventhandlers.ChannelDelete(bot.log.WithField(logFieldHandler, "ChannelDelete"))),
	}

	return func() {
		bot.log.Debugf("Deregistering handlers (count: %d)", len(closers))
		for _, closer := range closers {
			closer()
		}
	}
}
