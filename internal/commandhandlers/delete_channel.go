package commandhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

type DeleteChannelCommand struct {
	GuildID, ChannelID string
}

func DeleteChannelCommandHandler(c *DeleteChannelCommand, _ *discordgo.Session, log *logrus.Entry) {
	gc, _ := storage.Guilds.LoadOrStore(c.GuildID, &storage.GuildChannels{})

	_, err := gc.(*storage.GuildChannels).Delete(c.ChannelID)
	if err != nil {
		log.WithError(err).Error("Could not add channel")
		return
	}

}
