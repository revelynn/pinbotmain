package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	config.Configure()
	log.SetLevel(config.LogLevel)

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	if log.IsLevelEnabled(logrus.TraceLevel) {
		s.LogLevel = discordgo.LogDebug
	}

	bot := pinbot.New(config.ApplicationID, s, log)

	if config.TestGuildID != "" {
		bot.WithTestGuildID(config.TestGuildID)
	}

	if config.HealthCheckAddr != "" {
		bot.WithHealthCheck(config.HealthCheckAddr)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := bot.StartSession(sc); err != nil {
		os.Exit(1)
	}
}
