package commandhandlers

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/storage"
	"github.com/sirupsen/logrus"
)

const (
	emojiSeen = "👀"
	emojiErr  = "💩"
	emojiDone = "✅"
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

	l.Info("Pinning message")

	if !config.SelfPinEnabled && m.Author.ID == s.State.User.ID {
		l.Info("Ignoring self pin")
	}

	pinned, err := isAlreadyPinned(s, m)
	if err != nil {
		l.WithError(err).Error("Could not determine if message already pinned")
	}
	if pinned {
		l.Info("Message already pinned")
		return
	}

	// acknowledge the message
	l.Debug("Acknowledging message")
	if err = s.MessageReactionAdd(m.ChannelID, m.ID, emojiSeen); err != nil {
		l.WithError(err).Error("Could not acknowledge the message")
		return
	}

	// determine the target pin channel for the message
	pinChannel, err := getTargetChannel(c.GuildID, m.ChannelID)
	if err != nil {
		l.WithError(err).Error("Could not get target channel")
		err = s.MessageReactionAdd(m.ChannelID, m.ID, emojiErr)
		if err != nil {
			l.WithError(err).Error("Could not mark the message as failed")
		}
		return
	}

	l = l.WithField("target_channel_id", pinChannel)

	// send the pin message
	embed := &discordgo.MessageEmbed{
		Title:       "📌 New Pin",
		Description: fmt.Sprintf("%s said: %s", m.Author.Mention(), m.Content),
		URL:         fmt.Sprintf("https://discord.com/channels/%s/%s/%s", c.GuildID, m.ChannelID, m.ID),
	}

	if c.PinnedBy != nil {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Pinned by %s", c.PinnedBy.String()),
		}
	}

	// If there's only one attachment then add it to the message embed
	if len(m.Attachments) == 1 {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: m.Attachments[0].URL,
		}
	}

	pinMessage := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	// If there are multiple attachments then add them to separate embeds
	if len(m.Attachments) > 1 {
		for _, a := range m.Attachments {
			pinMessage.Embeds = append(pinMessage.Embeds, &discordgo.MessageEmbed{
				Type: discordgo.EmbedTypeImage,
				Image: &discordgo.MessageEmbedImage{
					URL:      a.URL,
					ProxyURL: a.ProxyURL,
					Width:    a.Width,
					Height:   a.Height,
				},
			})
		}
	}

	_, err = s.ChannelMessageSendComplex(pinChannel, pinMessage)
	if err != nil {
		l.WithError(err).Error("Could not send message")
		err = s.MessageReactionAdd(m.ChannelID, m.ID, emojiErr)
		if err != nil {
			l.WithError(err).Error("Could not mark the message as failed")
		}
		return
	}

	// mark the message as done
	l.Debug("Marking message as done")
	err = s.MessageReactionAdd(m.ChannelID, m.ID, emojiDone)
	if err != nil {
		l.WithError(err).Error("Could not mark the message as done")

		return
	}
}

func isAlreadyPinned(s *discordgo.Session, m *discordgo.Message) (bool, error) {
	acks, err := s.MessageReactions(m.ChannelID, m.ID, emojiDone, 0, "", "")
	if err != nil {
		return false, err
	}

	for _, ack := range acks {
		if ack.ID == s.State.User.ID {
			return true, nil
		}
	}

	return false, nil
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
