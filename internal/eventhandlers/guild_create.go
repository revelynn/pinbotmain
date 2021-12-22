package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/sirupsen/logrus"
)

func GuildCreate(log *logrus.Entry) func(s *discordgo.Session, e *discordgo.GuildCreate) {
	return func(s *discordgo.Session, e *discordgo.GuildCreate) {
		log.Debug("Guild info received:", e.Name)

		if !config.ShouldActOnGuild(e.Guild.ID) {
			return
		}

		commandhandlers.RegisterCommandsCommandHandler(&commandhandlers.RegisterCommandsCommand{
			ApplicationID: config.ApplicationID,
			GuildID:       e.Guild.ID,
		}, s, log)
	}
}
