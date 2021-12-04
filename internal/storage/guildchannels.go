package storage

import "github.com/bwmarrin/discordgo"

type GuildChannels struct {
	Channels []*discordgo.Channel
}

func (gc *GuildChannels) Add(c *discordgo.Channel) (bool, error) {
	if c.Type != discordgo.ChannelTypeGuildText {
		return false, nil
	}

	// skip adding if this channel already exists
	for _, channel := range gc.Channels {
		if channel.ID == c.ID {
			return false, nil
		}
	}

	gc.Channels = append(gc.Channels, c)

	return true, nil
}

func (gc *GuildChannels) Delete(id string) (bool, error) {
	for i, channel := range gc.Channels {
		if channel.ID == id {
			gc.Channels = append(gc.Channels[:i], gc.Channels[i+1:]...)
			return true, nil
		}
	}

	return false, nil
}
