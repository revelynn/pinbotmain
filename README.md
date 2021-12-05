# Pinbot

Whenever you react to a message with a pushpin 📌 emoji, Pinbot posts the message to a channel.

For a channel with name `#channel` it will look for the following channels:
* `#channel-pins` (a specific pin channel)
* `#pins` (a general pin channel)
* `#channel` (falls back to the channel itself)

You should set up a channel _that is only writeable by the bot_ as the pin channel