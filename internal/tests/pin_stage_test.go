package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/commandhandlers"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PinStage struct {
	t       *testing.T
	session *discordgo.Session
	require *require.Assertions
	assert  *assert.Assertions

	log     *logrus.Logger
	logHook *test.Hook

	sendMessage         *discordgo.MessageSend
	channel             *discordgo.Channel
	expectedPinsChannel *discordgo.Channel
	message             *discordgo.Message
	messages            []*discordgo.Message
}

func NewPinStage(t *testing.T) (*PinStage, *PinStage, *PinStage) {
	log := logrus.New()

	s := &PinStage{
		t:       t,
		log:     log,
		session: session,
		require: require.New(t),
		assert:  assert.New(t),
		logHook: test.NewLocal(log),
	}

	done := make(chan os.Signal, 1)

	go func() {
		bot := pinbot.New(config.ApplicationID, session, log)
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
	c, err := s.session.GuildChannelCreate(config.TestGuildID, name, discordgo.ChannelTypeGuildText)
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

func (s *PinStage) the_message_is_reacted_to_with(emoji string) *PinStage {
	err := s.session.MessageReactionAdd(s.message.ChannelID, s.message.ID, emoji)
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
	}, 5*time.Second, 10*time.Millisecond)

	return s
}

func (s *PinStage) the_bot_should_add_the_emoji(emoji string) *PinStage {
	s.require.Eventually(func() bool {
		reactions, err := s.session.MessageReactions(s.channel.ID, s.message.ID, emoji, 0, "", "")
		if err != nil {
			return false
		}

		for _, r := range reactions {
			if r.ID == s.session.State.User.ID {
				return true
			}
		}

		return false
	}, 5*time.Second, 100*time.Millisecond)

	return s
}

func (s *PinStage) handleMessageFor(channelID string) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID == channelID {
			s.messages = append(s.messages, m.Message)
		}
	}
}

type MockRoundTripper func(request *http.Request) (*http.Response, error)

func (m MockRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return m(request)
}

func (s *PinStage) the_message_is_already_pinned() {
	s.session.Client.Transport = MockRoundTripper(func(request *http.Request) (*http.Response, error) {
		expectedPath := discordgo.EndpointMessageReactions(s.message.ChannelID, s.message.ID, url.PathEscape("✅"))
		if request.Method == http.MethodGet && request.URL.String() == expectedPath {
			bs, err := json.Marshal([]*discordgo.User{
				{ID: s.session.State.User.ID},
			})
			s.require.NoError(err)

			s.session.Client.Transport = nil

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(bs)),
			}, nil
		}

		return http.DefaultTransport.RoundTrip(request)
	})
}

func (s *PinStage) the_bot_should_log_the_message_as_already_pinned() {
	s.require.Eventually(func() bool {
		for _, e := range s.logHook.AllEntries() {
			if e.Message == "Message already pinned" {
				return true
			}
		}

		return false
	}, 1*time.Second, 10*time.Millisecond)
}

func (s *PinStage) self_pin_is_disabled() *PinStage {
	c := config.SelfPinEnabled
	config.SelfPinEnabled = false

	s.t.Cleanup(func() {
		config.SelfPinEnabled = c
	})

	return s
}

func (s *PinStage) the_bot_should_log_the_message_as_an_avoided_self_pin() {
	s.require.Eventually(func() bool {
		for _, e := range s.logHook.AllEntries() {
			if e.Message == "Ignoring self pin" {
				return true
			}
		}

		return false
	}, 1*time.Second, 10*time.Millisecond)
}

func (s *PinStage) the_message_is_pinned() {
	s.require.NoError(s.session.ChannelMessagePin(s.channel.ID, s.message.ID))
}

func (s *PinStage) an_import_is_triggered() {
	commandhandlers.ImportChannelCommandHandler(&commandhandlers.ImportChannelCommand{
		GuildID:   config.TestGuildID,
		ChannelID: s.channel.ID,
	}, s.session, s.log.WithField("test", true))
}
