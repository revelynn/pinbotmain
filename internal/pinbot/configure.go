package pinbot

import (
	"github.com/elliotwms/pinbot/internal/eventhandlers"
)

const logFieldHandler = "handler"

func (bot *Bot) registerHandlers() func() {
	closers := []func(){
		bot.session.AddHandler(eventhandlers.Ready(bot.log.WithField(logFieldHandler, "Ready"))),
		bot.session.AddHandler(eventhandlers.MessageReactionAdd(bot.log.WithField(logFieldHandler, "MessageReactionAdd"))),
		bot.session.AddHandler(eventhandlers.GuildCreate(bot.log.WithField(logFieldHandler, "GuildCreate"))),
		bot.session.AddHandler(eventhandlers.InteractionCreate(bot.log.WithField(logFieldHandler, "InteractionCreate"))),
		bot.session.AddHandler(eventhandlers.ChannelPinsUpdate(bot.log.WithField(logFieldHandler, "ChannelPinsUpdate"))),
	}

	return func() {
		bot.log.Debugf("Deregistering handlers (count: %d)", len(closers))
		for _, closer := range closers {
			closer()
		}
	}
}
