package pinbot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	session     *discordgo.Session
	log         *logrus.Logger
	testGuildID *string
}

func New(s *discordgo.Session, l *logrus.Logger) *Bot {
	return &Bot{session: s, log: l}
}

func (bot *Bot) WithTestGuildID(id string) *Bot {
	bot.testGuildID = &id

	return bot
}

// Run runs the bot without tampering with the session
// This is useful in scenarios where the session is managed externally
func (bot *Bot) Run(notify chan os.Signal) error {
	cleanup := bot.registerHandlers()
	defer cleanup()
	<-notify

	return nil
}

// StartSession starts the session, calls Run (blocking until notify is received), then ends the session
func (bot *Bot) StartSession(notify chan os.Signal) error {
	bot.log.Info("Starting bot...")
	if err := bot.session.Open(); err != nil {
		bot.log.WithError(err).Error("Could not open session")
		return err
	}

	if err := bot.Run(notify); err != nil {
		bot.log.WithError(err).Error("Could not run bot")
		return err
	}

	bot.log.Info("Stopping bot...")
	if err := bot.session.Close(); err != nil {
		bot.log.WithError(err).Error("Could not close session")
		return err
	}

	return nil
}
