package storage

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

func TestGuildChannels_Add(t *testing.T) {
	gc := GuildChannels{}

	add, err := gc.Add(&discordgo.Channel{ID: "1", Name: "foo"})
	require.NoError(t, err)
	require.True(t, add)
	require.Len(t, gc.Channels, 1)

	add, err = gc.Add(&discordgo.Channel{ID: "1", Name: "foo"})
	require.NoError(t, err)
	require.False(t, add)
	require.Len(t, gc.Channels, 1)
}

func TestGuildChannels_AddMultiple(t *testing.T) {
	gc := GuildChannels{
		Channels: []*discordgo.Channel{
			{ID: "1", Name: "foo"},
		},
	}

	add, err := gc.Add(&discordgo.Channel{ID: "2", Name: "bar"})
	require.NoError(t, err)
	require.True(t, add)
	require.Len(t, gc.Channels, 2)
}

func TestGuildChannels_Delete(t *testing.T) {
	gc := GuildChannels{
		Channels: []*discordgo.Channel{
			{ID: "1", Name: "foo"},
		},
	}

	b, err := gc.Delete("1")
	require.NoError(t, err)
	require.True(t, b)
}

func TestGuildChannels_DeleteNotFound(t *testing.T) {
	gc := GuildChannels{
		Channels: []*discordgo.Channel{
			{ID: "1", Name: "foo"},
		},
	}

	b, err := gc.Delete("2")
	require.NoError(t, err)
	require.False(t, b)
}

func TestGuildChannels_DeleteFromMultiple(t *testing.T) {
	tests := map[string]string{
		"First":  "1",
		"Middle": "2",
		"Last":   "3",
	}

	for name, id := range tests {
		t.Run(name, func(t *testing.T) {
			gc := GuildChannels{
				Channels: []*discordgo.Channel{
					{ID: "1", Name: "foo"},
					{ID: "2", Name: "bar"},
					{ID: "3", Name: "bar"},
				},
			}

			b, err := gc.Delete(id)
			require.NoError(t, err)
			require.True(t, b)
			require.Len(t, gc.Channels, 2)

			for _, channel := range gc.Channels {
				require.NotEqual(t, channel.ID, id)
			}
		})
	}
}
