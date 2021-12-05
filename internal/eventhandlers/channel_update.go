package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func ChannelUpdate(log *logrus.Entry) func(_ *discordgo.Session, e *discordgo.ChannelUpdate) {
	return func(_ *discordgo.Session, e *discordgo.ChannelUpdate) {
		if err := addChannel(e.Channel); err != nil {
			log.WithError(err).Error("Could not add channel")
		}
	}
}
