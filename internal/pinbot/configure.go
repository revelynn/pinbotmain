package pinbot

import (
	"github.com/elliotwms/pinbot/internal/eventhandlers"
)

const logFieldHandler = "handler"

func (bot *Bot) configure() *Bot {
	bot.Session.AddHandler(eventhandlers.Ready(bot.Log.WithField(logFieldHandler, "Ready")))
	bot.Session.AddHandler(eventhandlers.MessageReactionAdd(bot.Log.WithField(logFieldHandler, "MessageReactionAdd"), bot.TestGuildID))
	bot.Session.AddHandler(eventhandlers.GuildCreate(bot.Log.WithField(logFieldHandler, "GuildCreate")))
	bot.Session.AddHandler(eventhandlers.ChannelCreate(bot.Log.WithField(logFieldHandler, "ChannelCreate")))
	bot.Session.AddHandler(eventhandlers.ChannelUpdate(bot.Log.WithField(logFieldHandler, "ChannelUpdate")))
	bot.Session.AddHandler(eventhandlers.ChannelDelete(bot.Log.WithField(logFieldHandler, "ChannelDelete")))

	return bot
}
