# Pinbot

[Install 📌](https://discord.com/oauth2/authorize?client_id=921554139740254209&permissions=3136&redirect_uri=https%3A%2F%2Fgithub.com%2Felliotwms%2Fpinbot&scope=applications.commands%20bot)

Whenever you react to a message with a pushpin 📌 emoji, Pinbot posts the message to a channel.

For a channel with name `#channel` it will look for the following channels:

* `#channel-pins` (a specific pin channel)
* `#pins` (a general pin channel)
* `#channel` (falls back to the channel itself)

You should set up a channel _that is only writeable by the bot_ as the pin channel

⚠️ Note that this bot is currently in _alpha_ and is used at your own risk! ⚠️

### Permissions

Pinbot is designed to be run with as few permissions as possible. In order for it to work, it needs to be able to see
channels and messages, and write messages to whichever channel you have set up as its pin output channel.

## Run

```shell
export TOKEN {bot_token}
export APPLICATION_ID {bot_application_id}
docker run -e TOKEN -e APPLICATION_ID ghcr.io/elliotwms/pinbot:{version}
```

## Configuration

| Variable            | Description                                                            | Required |
|---------------------|------------------------------------------------------------------------|----------|
| `TOKEN`             | Bot token ID                                                           | `true`   |
| `APPLICATION_ID`    | Bot application ID                                                     | `true`   |
| `TEST_GUILD_ID`     | When specified, the bot should only respond to pins in this test guild | `false`  |
| `HEALTH_CHECK_ADDR` | Address to serve the `/v1/health/` endpoint on (e.g. `:8080`)          | `false`  |
| `EXCLUDED_CHANNELS` | Comma-separated list of excluded channel IDs                           | `false`  |
| `DEBUG`             | Output debug Discord client logs                                       | `false`  |

## Testing

`/tests` contains a suite of integration tests which run against the real Discord API in a test guild. It will create
and destroy a guild for the test run.

In order to run these test yourself you will need to:

* [Create a new bot](https://discord.com/developers/applications), obtaining the bot token
* Run the tests with the following environment variables:
    * `TOKEN`: the bot token

### Debugging

If you would rather view the bot activity, then it's also possible to use an existing guild instead of creating and
destroying one for each test run.

* Invite your bot to your guild, giving it the usual bot permissions as well as:
    * Manage Channels (to create channels during tests)
    * Read Messages (to assert on message creation)
    * Manage messages (to pin messages)
    * `https://discord.com/oauth2/authorize?client_id={bot_client_id}&permissions=76880&redirect_uri=http%3A%2F%2Flocalhost&scope=applications.commands%20bot`
* Set the `TEST_GUILD_ID` environment variable to the test guild's ID
