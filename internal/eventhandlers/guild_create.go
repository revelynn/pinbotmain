package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func GuildCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.GuildCreate) {
	return func(s *discordgo.Session, e *discordgo.GuildCreate) {
		log.Info("Guild info received:", e.Name)

		for _, c := range e.Channels {
			commandhandlers.SaveChannelCommandHandler(&commandhandlers.SaveChannelCommand{
				GuildID: e.Guild.ID,
				Channel: c,
			}, s, log)
		}
	}
}
