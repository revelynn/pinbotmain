package tests

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PinStage struct {
	t       *testing.T
	session *discordgo.Session
	require *require.Assertions
	assert  *assert.Assertions

	sendMessage         *discordgo.MessageSend
	channel             *discordgo.Channel
	expectedPinsChannel *discordgo.Channel
	message             *discordgo.Message
	messages            []*discordgo.Message
}

func NewPinStage(t *testing.T) (*PinStage, *PinStage, *PinStage) {
	s := &PinStage{
		t:       t,
		session: session,
		require: require.New(t),
		assert:  assert.New(t),
	}

	done := make(chan os.Signal, 1)

	go func() {
		bot := pinbot.New(session, logrus.New())
		s.require.NoError(bot.Run(done))
	}()

	t.Cleanup(func() {
		done <- os.Interrupt
	})

	return s, s, s
}

func (s *PinStage) and() *PinStage {
	return s
}

func (s *PinStage) a_channel_named(name string) *PinStage {
	c, err := s.session.GuildChannelCreate(testGuild, name, discordgo.ChannelTypeGuildText)
	s.require.NoError(err)

	s.t.Cleanup(func() {
		_, err = s.session.ChannelDelete(c.ID)
		s.assert.NoError(err)
	})

	if s.channel == nil {
		// register the first created channel as the "default" channel for the stage
		s.channel = c
	}
	// register the last created channel as the expected pins channel
	s.expectedPinsChannel = c

	s.session.AddHandler(s.handleMessageFor(c.ID))

	return s
}

func (s *PinStage) a_message() *PinStage {
	s.sendMessage = &discordgo.MessageSend{
		Content: "Hello, World!",
	}

	return s
}

func (s *PinStage) the_message_is_posted() *PinStage {
	if s.message == nil {
		s.a_message()
	}

	m, err := s.session.ChannelMessageSendComplex(s.channel.ID, s.sendMessage)
	s.require.NoError(err)
	s.message = m

	return s
}

func (s *PinStage) the_message_is_reacted_to() *PinStage {
	err := s.session.MessageReactionAdd(s.message.ChannelID, s.message.ID, "📌")
	s.require.NoError(err)

	return s
}

func (s *PinStage) a_pin_message_should_be_posted_in_the_last_channel() *PinStage {
	s.require.Eventually(func() bool {
		for _, m := range s.messages {
			if m.ChannelID != s.expectedPinsChannel.ID {
				continue
			}

			if len(m.Embeds) != 1 {
				continue
			}

			embed := m.Embeds[0]
			return embed.Title == "📌 New Pin" && strings.Contains(embed.Description, s.sendMessage.Content)
		}

		return false
	}, time.Second, 10*time.Millisecond)

	return s
}

func (s *PinStage) handleMessageFor(channelID string) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID == channelID {
			s.messages = append(s.messages, m.Message)
		}
	}
}
