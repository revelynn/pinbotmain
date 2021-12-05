package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

func ChannelDelete(log *logrus.Entry) func(_ *discordgo.Session, e *discordgo.ChannelDelete) {
	return func(_ *discordgo.Session, e *discordgo.ChannelDelete) {
		gc, _ := storage.Guilds.LoadOrStore(e.GuildID, &storage.GuildChannels{})

		_, err := gc.(*storage.GuildChannels).Delete(e.Channel.ID)
		if err != nil {
			log.WithError(err).Error("Could not add channel")
			return
		}
	}
}
