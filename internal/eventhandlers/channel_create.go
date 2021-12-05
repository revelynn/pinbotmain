package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

func ChannelCreate(log *logrus.Entry) func(_ *discordgo.Session, e *discordgo.ChannelCreate) {
	return func(_ *discordgo.Session, e *discordgo.ChannelCreate) {
		if err := addChannel(e.Channel); err != nil {
			log.WithError(err).Error("Could not add channel")
		}
	}
}

func addChannel(e *discordgo.Channel) error {
	gc, _ := storage.Guilds.LoadOrStore(e.GuildID, &storage.GuildChannels{})

	_, err := gc.(*storage.GuildChannels).Add(e)

	return err
}
