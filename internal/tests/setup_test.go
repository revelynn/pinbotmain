package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var (
	botToken  string
	testGuild string
	session   *discordgo.Session
)

func TestMain(m *testing.M) {
	configure()
	openSession()
	defer closeSession()

	m.Run()
	os.Exit(0)
}

func configure() {
	botToken = mustGetEnv("TOKEN")
	testGuild = mustGetEnv("TEST_GUILD_ID")
}

func mustGetEnv(e string) string {
	s := os.Getenv(e)
	if s == "" {
		panic(fmt.Errorf("missing env '%s'", e))
	}

	return s
}

func openSession() {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", botToken))
	if err != nil {
		panic(err)
	}

	session.LogLevel = discordgo.LogDebug

	if err := session.Open(); err != nil {
		panic(err)
	}
}

func closeSession() {
	if err := session.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
