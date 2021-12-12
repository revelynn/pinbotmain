package tests

import "testing"

func TestPin(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to()

	then.
		a_pin_message_should_be_posted_in_the_last_channel().and().
		the_bot_should_add_the_emoji("👀").and().
		the_bot_should_add_the_emoji("✅")
}

func TestPinGeneralPinsChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_channel_named("pins").and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to()

	then.
		a_pin_message_should_be_posted_in_the_last_channel()
}

func TestPinSpecificPinsChannel(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		a_channel_named("pins").and().
		a_channel_named("test-pins").and().
		the_message_is_posted()

	when.
		the_message_is_reacted_to()

	then.
		a_pin_message_should_be_posted_in_the_last_channel()
}

func TestPinDuplicate(t *testing.T) {
	given, when, then := NewPinStage(t)

	given.
		a_channel_named("test").and().
		the_message_is_posted().and().
		the_message_is_already_pinned()

	when.
		the_message_is_reacted_to()

	then.
		the_bot_should_log_the_message_as_already_pinned()
}
