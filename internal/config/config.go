package config

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

var (
	Token         string
	ApplicationID string

	TestGuildID string

	SelfPinEnabled bool
)

var once sync.Once

func Configure() {
	once.Do(func() {
		Token = mustGetEnv("TOKEN")
		ApplicationID = mustGetEnv("APPLICATION_ID")
		TestGuildID = os.Getenv("TEST_GUILD_ID")
		SelfPinEnabled = strings.ToLower(os.Getenv("SELF_PIN_ENABLED")) == "true"
	})
}

func mustGetEnv(s string) string {
	token := os.Getenv(s)
	if token == "" {
		panic(fmt.Sprintf("Missing '%s'", s))
	}
	return token
}
