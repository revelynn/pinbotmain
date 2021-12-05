package eventhandlers

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
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

		log.Info("Pinning message")

		m, err := s.ChannelMessage(e.ChannelID, e.MessageID)
		if err != nil {
			log.WithError(err).Error("Could not get channel message")
			return
		}

		pinChannel, err := getTargetChannel(e.GuildID, e.ChannelID)
		if err != nil {
			log.WithError(err).Error("Could not get target channel")
			return
		}

		_, err = s.ChannelMessageSendEmbed(pinChannel, &discordgo.MessageEmbed{
			Type:        discordgo.EmbedTypeRich,
			Title:       e.Emoji.Name + " New Pin",
			Description: fmt.Sprintf("%s said: %s", m.Author.Mention(), m.Content),
			URL:         fmt.Sprintf("https://discord.com/channels/%s/%s/%s", e.GuildID, m.ChannelID, m.ID),
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Pinned by %s", e.Member.User.String()),
			},
		})
		if err != nil {
			log.WithError(err).Error("Could not send message")
		}
	}
}

// getTargetChannel returns the target pin channel for a given channel #channel in the following order:
// #channel-pins (a specific pin channel)
// #pins (a generic pin channel)
// #channel (the channel itself)
func getTargetChannel(guildID, channelID string) (string, error) {
	k, _ := storage.Guilds.LoadOrStore(guildID, &storage.GuildChannels{})
	gc, ok := k.(*storage.GuildChannels)
	if !ok {
		return "", errors.New("map did not contain type *storage.GuildChannels")
	}

	// get the channel
	var channel *discordgo.Channel

	for _, c := range gc.Channels {
		if c.ID == channelID {
			channel = c
			break
		}
	}

	if channel == nil {
		return "", errors.New("missing channel from map")
	}

	// check for #channel-pins
	for _, c := range gc.Channels {
		if c.Name == channel.Name+"-pins" {
			return c.ID, nil
		}
	}

	for _, c := range gc.Channels {
		if c.Name == "pins" {
			return c.ID, nil
		}
	}

	return channel.ID, nil
}
