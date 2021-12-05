package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

func GuildCreate(log *logrus.Entry) func(_ *discordgo.Session, e *discordgo.GuildCreate) {
	return func(_ *discordgo.Session, e *discordgo.GuildCreate) {
		log.Info("Guild info received:", e.Name)

		gc, _ := storage.Guilds.LoadOrStore(e.Guild.ID, &storage.GuildChannels{})

		for _, c := range e.Channels {
			_, err := gc.(*storage.GuildChannels).Add(c)
			if err != nil {
				log.WithError(err).Error("Could not add channel")
				return
			}
		}
	}
}
