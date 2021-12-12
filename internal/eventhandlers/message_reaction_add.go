package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func MessageReactionAdd(log *logrus.Entry, testGuildID *string) func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		log.WithField("emoji", e.Emoji.Name).Info("Received reaction")

		if e.Emoji.Name != "📌" {
			return
		}

		if testGuildID != nil && *testGuildID != e.GuildID {
			log.Info("Skipping non-test guild")
			return
		}

		reactions, err := s.MessageReactions(e.ChannelID, e.MessageID, e.Emoji.APIName(), 0, "", "")
		if err != nil {
			log.WithError(err).Error("Could not get message reactions")
			return
		}

		if len(reactions) > 1 {
			log.WithField("reactions", len(reactions)).Info("Message already pinned")
			return
		}

		log.Info("Pinning message")

		commandhandlers.PinMessageCommandHandler(&commandhandlers.PinMessageCommand{Event: e}, s, log)
	}
}
