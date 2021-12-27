package commandhandlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	emojiSeen    = "👀"
	emojiDone    = "✅"
	emojiErr     = "💩"
	emojiSelfPin = "🔄"
	emojiNo      = "🚫"
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

	// acknowledge the message
	l.Debug("Acknowledging message")
	react(s, m, emojiSeen, l)

	if !config.SelfPinEnabled && m.Author.ID == s.State.User.ID {
		l.Debug("Ignoring self pin")
		react(s, m, emojiSelfPin, l)
		return
	}

	if config.IsExcludedChannel(m.ChannelID) {
		l.Debug("Skipping excluded channel")
		react(s, m, emojiNo, l)
		return
	}

	pinned, err := isAlreadyPinned(s, m)
	if err != nil {
		l.WithError(err).Error("Could not determine if message already pinned")
	}
	if pinned {
		l.Info("Message already pinned")
		return
	}

	// determine the target pin channel for the message
	pinChannel, err := getTargetChannel(s, c.GuildID, m.ChannelID)
	if err != nil {
		l.WithError(err).Error("Could not get target channel")
		react(s, m, emojiErr, l)
		return
	}

	l = l.WithField("target_channel_id", pinChannel)

	// build the rich embed pin message
	pinMessage := buildPinMessage(c, m)

	// send the pin message
	_, err = s.ChannelMessageSendComplex(pinChannel, pinMessage)
	if err != nil {
		l.WithError(err).Error("Could not send message")
		react(s, m, emojiErr, l)
		return
	}

	// mark the message as done
	l.Debug("Marking message as done")
	react(s, m, emojiDone, l)
}

func react(s *discordgo.Session, m *discordgo.Message, emoji string, l *logrus.Entry) {
	if err := s.MessageReactionAdd(m.ChannelID, m.ID, emoji); err != nil {
		l.WithError(err).Error("Could not react to message")
	}
}

func buildPinMessage(c *PinMessageCommand, m *discordgo.Message) *discordgo.MessageSend {
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

	pinMessage := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}

	// If there are multiple attachments then add them to separate embeds
	for i, a := range m.Attachments {
		if a.Width == 0 || a.Height == 0 {
			// only embed images
			continue
		}
		e := &discordgo.MessageEmbedImage{URL: a.URL}

		if i == 0 {
			// add the first image to the existing embed
			pinMessage.Embeds[0].Image = e
		} else {
			// add any other images to their own embed
			pinMessage.Embeds = append(pinMessage.Embeds, &discordgo.MessageEmbed{
				Type:  discordgo.EmbedTypeImage,
				Image: e,
			})
		}
	}
	return pinMessage
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
func getTargetChannel(s *discordgo.Session, guildID, channelID string) (string, error) {
	origin, err := s.State.GuildChannel(guildID, channelID)
	if err != nil {
		return "", err
	}

	guild, err := s.State.Guild(guildID)
	if err != nil {
		return "", err
	}

	// use the same channel by default
	channel := origin

	// check for #channel-pins
	for _, c := range guild.Channels {
		if c.Name == channel.Name+"-pins" {
			return c.ID, nil
		}
	}

	for _, c := range guild.Channels {
		if c.Name == "pins" {
			return c.ID, nil
		}
	}

	return channel.ID, nil
}
