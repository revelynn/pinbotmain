package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/config"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	config.Configure()
	log.Infof("Excluding channels: %s", config.ExcludedChannels)

	s, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		s.LogLevel = discordgo.LogDebug
	}

	bot := pinbot.New(config.ApplicationID, s, log)

	if id := os.Getenv("TEST_GUILD_ID"); id != "" {
		bot.WithTestGuildID(id)
	}

	if addr := os.Getenv("HEALTH_CHECK_ADDR"); addr != "" {
		bot.WithHealthCheck(addr)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	if err := bot.StartSession(sc); err != nil {
		os.Exit(1)
	}
}
