package commandhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

type SaveChannelCommand struct {
	GuildID string
	Channel *discordgo.Channel
}

func SaveChannelCommandHandler(c *SaveChannelCommand, _ *discordgo.Session, log *logrus.Entry) {
	gc, _ := storage.Guilds.LoadOrStore(c.GuildID, &storage.GuildChannels{})
	_, err := gc.(*storage.GuildChannels).Add(c.Channel)
	if err != nil {
		log.WithError(err).Error("Could not add channel")
		return
	}
}
