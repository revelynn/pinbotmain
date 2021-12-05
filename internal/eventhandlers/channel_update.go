package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func ChannelUpdate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.ChannelUpdate) {
	return func(s *discordgo.Session, e *discordgo.ChannelUpdate) {
		commandhandlers.SaveChannelCommandHandler(&commandhandlers.SaveChannelCommand{
			GuildID: e.GuildID,
			Channel: e.Channel,
		}, s, log)
	}
}
