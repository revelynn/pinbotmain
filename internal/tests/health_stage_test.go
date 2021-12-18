package tests

import (
	"github.com/bwmarrin/discordgo"
	"github.com/elliotwms/pinbot/internal/pinbot"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
)

type HealthStage struct {
	t       *testing.T
	session *discordgo.Session
	require *require.Assertions
	assert  *assert.Assertions
	logHook *test.Hook
	res     *http.Response
}

func NewHealthStage(t *testing.T) (*HealthStage, *HealthStage, *HealthStage) {
	log := logrus.New()

	s := &HealthStage{
		t:       t,
		session: session,
		require: require.New(t),
		assert:  assert.New(t),
		logHook: test.NewLocal(log),
	}

	done := make(chan os.Signal, 1)

	go func() {
		bot := pinbot.New(session, log)
		s.require.NoError(bot.WithHealthCheck(":8080").Run(done))
	}()

	t.Cleanup(func() {
		done <- os.Interrupt
	})

	return s, s, s
}

func (s *HealthStage) the_bot_is_running() *HealthStage {
	return s // no-op
}

func (s *HealthStage) a_health_check_request_is_sent() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/v1/health", nil)
	s.require.NoError(err)
	s.res, err = http.DefaultClient.Do(req)
	s.require.NoError(err)
}

func (s *HealthStage) a_response_should_be_received_with_status_code(code int) {
	s.require.Equal(code, s.res.StatusCode)
}
