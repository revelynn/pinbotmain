package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func ChannelDelete(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.ChannelDelete) {
	return func(s *discordgo.Session, e *discordgo.ChannelDelete) {
		commandhandlers.DeleteChannelCommandHandler(&commandhandlers.DeleteChannelCommand{
			GuildID:   e.GuildID,
			ChannelID: e.Channel.ID,
		}, s, log)
	}
}
