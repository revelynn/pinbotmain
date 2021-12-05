package pinbot

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	Session     *discordgo.Session
	Log         *logrus.Logger
	TestGuildID *string
}

func New(s *discordgo.Session, l *logrus.Logger) *Bot {
	bot := &Bot{Session: s, Log: l}

	return bot.configure()
}

func (bot *Bot) WithTestGuildID(id string) *Bot {
	bot.TestGuildID = &id

	return bot
}

func (bot *Bot) Run(notify chan os.Signal) error {
	bot.Log.Info("Starting bot...")
	if err := bot.Session.Open(); err != nil {
		return err
	}

	<-notify

	bot.Log.Info("Stopping bot...")
	if err := bot.Session.Close(); err != nil {
		return err
	}

	return nil
}
