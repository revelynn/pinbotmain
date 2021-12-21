package eventhandlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/sirupsen/logrus"
)

func MessageReactionAdd(log *logrus.Entry, testGuildID string) func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		log.WithField("emoji", e.Emoji.Name).Info("Received reaction")

		if e.Emoji.Name != "📌" {
			return
		}

		if testGuildID != "" && testGuildID != e.GuildID {
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

		m, err := s.ChannelMessage(e.ChannelID, e.MessageID)
		if err != nil {
			log.WithError(err).Error("Could not get channel message")
			return
		}

		commandhandlers.PinMessageCommandHandler(&commandhandlers.PinMessageCommand{
			GuildID:  e.GuildID,
			Message:  m,
			PinnedBy: e.Member.User,
		}, s, log)
	}
}
