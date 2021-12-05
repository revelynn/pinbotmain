package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func ChannelCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.ChannelCreate) {
	return func(s *discordgo.Session, e *discordgo.ChannelCreate) {
		commandhandlers.SaveChannelCommandHandler(&commandhandlers.SaveChannelCommand{
			GuildID: e.GuildID,
			Channel: e.Channel,
		}, s, log)
	}
}
