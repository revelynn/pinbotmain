# Pinbot

Whenever you react to a message with a pushpin 📌 emoji, Pinbot posts the message to a channel.

For a channel with name `#channel` it will look for the following channels:
* `#channel-pins` (a specific pin channel)
* `#pins` (a general pin channel)
* `#channel` (falls back to the channel itself)

You should set up a channel _that is only writeable by the bot_ as the pin channel

## Testing
`/tests` contains a suite of integration tests which run against the real Discord API in a test guild. It will create and destroy a guild for the test run.

In order to run these test yourself you will need to:
* [Create a new bot](https://discord.com/developers/applications), obtaining the bot Client ID and Token
* Run the tests with the following environment variables:
  * `TOKEN`: the bot token

### Debugging
If you would rather view the bot activity, then it's also possible to use an existing guild instead of creating and destroying one for each test run.
 
* Invite your bot to your guild, giving it the usual bot permissions as well as:
  * Manage Channels (to create channels during tests)
  * Read Messages (to assert on message creation)
  * `https://discord.com/oauth2/authorize?client_id={bot_client_id}&permissions=68688&redirect_uri=http%3A%2F%2Flocalhost&scope=bot`
* Set the `TEST_GUILD_ID` environment variable to the test guild's ID