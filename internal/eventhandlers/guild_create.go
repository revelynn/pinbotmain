package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func GuildCreate(log *logrus.Entry, applicationID string, guildID string) func(s *discordgo.Session, e *discordgo.GuildCreate) {
	return func(s *discordgo.Session, e *discordgo.GuildCreate) {
		log.Info("Guild info received:", e.Name)

		if guildID != "" && e.Guild.ID != guildID {
			log.Debugf("Ignoring non-test guild")
			return
		}

		commandhandlers.RegisterCommandsCommandHandler(&commandhandlers.RegisterCommandsCommand{
			ApplicationID: applicationID,
			GuildID:       guildID,
		}, s, log)
	}
}
