package commandhandlers

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

type PinMessageCommand struct {
	GuildID  string
	Message  *discordgo.Message
	PinnedBy *discordgo.User
}

func PinMessageCommandHandler(c *PinMessageCommand, s *discordgo.Session, log *logrus.Entry) {
	m := c.Message
	l := log.WithFields(map[string]interface{}{
		"guild_id":   c.GuildID,
		"channel_id": m.ChannelID,
		"message_id": m.ID,
	})

	// acknowledge the message
	l.Debug("Acknowledging message")
	err := s.MessageReactionAdd(m.ChannelID, m.ID, "👀")
	if err != nil {
		l.WithError(err).Error("Could not acknowledge the message")
		return
	}

	// determine the target pin channel for the message
	pinChannel, err := getTargetChannel(c.GuildID, m.ChannelID)
	if err != nil {
		l.WithError(err).Error("Could not get target channel")
		err = s.MessageReactionAdd(m.ChannelID, m.ID, "🤔")
		if err != nil {
			l.WithError(err).Error("Could not mark the message as failed")
		}
		return
	}

	l = l.WithField("target_channel_id", pinChannel)

	// send the pin message
	pinMessage := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       "📌 New Pin",
		Description: fmt.Sprintf("%s said: %s", m.Author.Mention(), m.Content),
		URL:         fmt.Sprintf("https://discord.com/channels/%s/%s/%s", c.GuildID, m.ChannelID, m.ID),
	}

	if c.PinnedBy != nil {
		pinMessage.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Pinned by %s", c.PinnedBy.String()),
		}
	}

	_, err = s.ChannelMessageSendEmbed(pinChannel, pinMessage)
	if err != nil {
		l.WithError(err).Error("Could not send message")
		err = s.MessageReactionAdd(m.ChannelID, m.ID, "💩")
		if err != nil {
			l.WithError(err).Error("Could not mark the message as failed")
		}
		return
	}

	// mark the message as done
	l.Debug("Marking message as done")
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		l.WithError(err).Error("Could not mark the message as done")

		return
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
