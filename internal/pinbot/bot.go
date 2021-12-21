package pinbot

import (
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	session         *discordgo.Session
	log             *logrus.Logger
	applicationID   string
	testGuildID     string
	healthCheckAddr *string
}

func New(id string, s *discordgo.Session, l *logrus.Logger) *Bot {
	return &Bot{session: s, log: l, applicationID: id}
}

func (bot *Bot) WithTestGuildID(id string) *Bot {
	bot.testGuildID = id

	return bot
}

func (bot *Bot) WithHealthCheck(addr string) *Bot {
	bot.healthCheckAddr = &addr

	return bot
}

// Run runs the bot without tampering with the session
// This is useful in scenarios where the session is managed externally
func (bot *Bot) Run(notify chan os.Signal) error {
	cleanup := bot.registerHandlers()
	defer cleanup()

	if bot.healthCheckAddr != nil {
		go bot.httpListen()
	}

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

func (bot *Bot) httpListen() {
	http.HandleFunc("/v1/health", func(w http.ResponseWriter, req *http.Request) {
		latency := bot.session.HeartbeatLatency()
		// fail if we have not received a heartbeat response in more than 5 minutes
		if latency > 5*time.Minute {
			w.WriteHeader(500)
		}

		if _, err := w.Write([]byte(latency.String())); err != nil {
			bot.log.WithError(err).Error("Could not write health check response")
		}
	})

	err := http.ListenAndServe(*bot.healthCheckAddr, nil)
	if err != nil {
		bot.log.WithError(err).Error("Could not serve health check endpoint")
		return
	}
}
