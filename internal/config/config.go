package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	Token            string
	ApplicationID    string
	TestGuildID      string
	HealthCheckAddr  string
	LogLevel         logrus.Level
	SelfPinEnabled   bool
	ExcludedChannels []string
)

var once sync.Once

func Configure() {
	once.Do(func() {
		Token = mustGetEnv("TOKEN")
		ApplicationID = mustGetEnv("APPLICATION_ID")
		TestGuildID = os.Getenv("TEST_GUILD_ID")
		SelfPinEnabled = strings.ToLower(os.Getenv("SELF_PIN_ENABLED")) == "true"
		HealthCheckAddr = os.Getenv("HEALTH_CHECK_ADDR")

		if s, ok := os.LookupEnv("EXCLUDED_CHANNELS"); ok && s != "" {
			ExcludedChannels = strings.Split(s, ",")
		}

		l, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
		if err != nil {
			LogLevel = logrus.InfoLevel
		} else {
			LogLevel = l
		}
	})
}

func Output(showSensitive bool) logrus.Fields {
	fields := logrus.Fields{
		"APPLICATION_ID":    ApplicationID,
		"TEST_GUILD_ID":     TestGuildID,
		"HEALTH_CHECK_ADDR": HealthCheckAddr,
		"LOG_LEVEL":         LogLevel,
		"SELF_PIN_ENABLED":  SelfPinEnabled,
		"EXCLUDED_CHANNELS": ExcludedChannels,
	}

	if showSensitive {
		fields["TOKEN"] = Token
	}

	return fields
}

func mustGetEnv(s string) string {
	token := os.Getenv(s)
	if token == "" {
		panic(fmt.Sprintf("Missing '%s'", s))
	}
	return token
}

func IsExcludedChannel(id string) bool {
	for _, c := range ExcludedChannels {
		if c == id {
			return true
		}
	}

	return false
}

func ShouldActOnGuild(id string) bool {
	return TestGuildID == "" || TestGuildID == id
}
