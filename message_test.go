package twitch

import (
	"testing"
)

func TestCanParseMessage(t *testing.T) {
	testMessage := "@badges=subscriber/6,premium/1;color=#FF0000;display-name=Redflamingo13;emotes=;id=2a31a9df-d6ff-4840-b211-a2547c7e656e;mod=0;room-id=11148817;subscriber=1;tmi-sent-ts=1490382457309;turbo=0;user-id=78424343;user-type= :redflamingo13!redflamingo13@redflamingo13.tmi.twitch.tv PRIVMSG #pajlada :Thrashh5, FeelsWayTooAmazingMan kinda"
	message := parseMessage(testMessage)

	assertStringsEqual(t, "pajlada", message.Channel)
	assertIntsEqual(t, 6, message.Badges["subscriber"])
	assertStringsEqual(t, "#FF0000", message.Color)
	assertStringsEqual(t, "Redflamingo13", message.DisplayName)
	assertIntsEqual(t, 0, len(message.Emotes))
	assertStringsEqual(t, "0", message.Tags["mod"])
	assertStringsEqual(t, "Thrashh5, FeelsWayTooAmazingMan kinda", message.Text)
	if message.Type != PRIVMSG {
		t.Error("parsing message type failed")
	}
	assertStringsEqual(t, "redflamingo13", message.Username)
	assertStringsEqual(t, "", message.UserType)
	assertFalse(t, message.Action, "parsing action failed")
}

func TestCanParseActionMessage(t *testing.T) {
	testMessage := "@badges=subscriber/6,premium/1;color=#FF0000;display-name=Redflamingo13;emotes=;id=2a31a9df-d6ff-4840-b211-a2547c7e656e;mod=0;room-id=11148817;subscriber=1;tmi-sent-ts=1490382457309;turbo=0;user-id=78424343;user-type= :redflamingo13!redflamingo13@redflamingo13.tmi.twitch.tv PRIVMSG #pajlada :\u0001ACTION Thrashh5, FeelsWayTooAmazingMan kinda"
	message := parseMessage(testMessage)

	assertStringsEqual(t, "pajlada", message.Channel)
	assertIntsEqual(t, 6, message.Badges["subscriber"])
	assertStringsEqual(t, "#FF0000", message.Color)
	assertStringsEqual(t, "Redflamingo13", message.DisplayName)
	assertIntsEqual(t, 0, len(message.Emotes))
	assertStringsEqual(t, "0", message.Tags["mod"])
	assertStringsEqual(t, "Thrashh5, FeelsWayTooAmazingMan kinda", message.Text)
	if message.Type != PRIVMSG {
		t.Error("parsing message type failed")
	}
	assertStringsEqual(t, "redflamingo13", message.Username)
	assertStringsEqual(t, "", message.UserType)
	assertTrue(t, message.Action, "parsing action failed")
}

func TestCanParseWhisper(t *testing.T) {
	testMessage := "@badges=;color=#00FF7F;display-name=Danielps1;emotes=;message-id=20;thread-id=32591953_77829817;turbo=0;user-id=32591953;user-type= :danielps1!danielps1@danielps1.tmi.twitch.tv WHISPER gempir :i like memes"
	message := parseMessage(testMessage)

	assertIntsEqual(t, 0, message.Badges["subscriber"])
	assertStringsEqual(t, "#00FF7F", message.Color)
	assertStringsEqual(t, "Danielps1", message.DisplayName)
	assertIntsEqual(t, 0, len(message.Emotes))
	assertStringsEqual(t, "", message.Tags["mod"])
	assertStringsEqual(t, "i like memes", message.Text)
	if message.Type != WHISPER {
		t.Error("parsing message type failed")
	}
	assertStringsEqual(t, "danielps1", message.Username)
	assertFalse(t, message.Action, "parsing action failed")
}

func TestCantParseNoTagsMessage(t *testing.T) {
	testMessage := "my test message"

	message := parseMessage(testMessage)

	assertStringsEqual(t, testMessage, message.Text)
}

func TestCantParseInvalidMessage(t *testing.T) {
	testMessage := "@my :test message"

	message := parseMessage(testMessage)

	assertStringsEqual(t, "", message.Text)
}

func TestCanParseClearChatMessage(t *testing.T) {
	testMessage := `@ban-duration=1;ban-reason=testing\sxd;room-id=11148817;target-user-id=40910607 :tmi.twitch.tv CLEARCHAT #pajlada :ampzyh`

	message := parseMessage(testMessage)

	if message.Type != CLEARCHAT {
		t.Error("parsing CLEARCHAT message failed")
	}
}

func TestCanParseEmoteMessage(t *testing.T) {
	testMessage := "@badges=;color=#008000;display-name=Zugren;emotes=120232:0-6,13-19,26-32,39-45,52-58;id=51c290e9-1b50-497c-bb03-1667e1afe6e4;mod=0;room-id=11148817;sent-ts=1490382458685;subscriber=0;tmi-sent-ts=1490382456776;turbo=0;user-id=65897106;user-type= :zugren!zugren@zugren.tmi.twitch.tv PRIVMSG #pajlada :TriHard Clap TriHard Clap TriHard Clap TriHard Clap TriHard Clap"

	message := parseMessage(testMessage)

	assertIntsEqual(t, 1, len(message.Emotes))
}

func TestParseUsernameMiddleRegex(t *testing.T) {
	testMessage := "thexin1!thexin1@thexin1.tmi.twitch.tv PRIVMSG #n1nja"
	username, mType, channel := parseMiddle(testMessage)

	assertStringsEqual(t, "thexin1", username)
	assertIntsEqual(t, int(PRIVMSG), int(mType))
	assertStringsEqual(t, "n1nja", channel)
}

func TestParseNoUserMiddleRegex(t *testing.T) {
	testMessage := "tmi.twitch.tv ROOMSTATE #dallas"
	username, mType, channel := parseMiddle(testMessage)

	assertStringsEqual(t, "", username)
	assertIntsEqual(t, int(ROOMSTATE), int(mType))
	assertStringsEqual(t, "dallas", channel)
}

func TestCanParseUsernoticeResubMessage(t *testing.T) {
	testMessage := `@badges=staff/1,broadcaster/1,turbo/1;color=#008000;display-name=ronni;emotes=;id=db25007f-7a18-43eb-9379-80131e44d633;login=ronni;mod=0;msg-id=resub;msg-param-months=6;msg-param-sub-plan=Prime;msg-param-sub-plan-name=Prime;room-id=1337;subscriber=1;system-msg=ronni\shas\ssubscribed\sfor\s6\smonths!;tmi-sent-ts=1507246572675;turbo=1;user-id=1337;user-type=staff :tmi.twitch.tv USERNOTICE #dallas :Great stream -- keep it up!`

	message := parseMessage(testMessage)

	assertIntsEqual(t, int(USERNOTICE), int(message.Type))
	assertStringsEqual(t, "dallas", message.Channel)
	assertStringsEqual(t, "ronni", message.Tags["login"])
	assertStringsEqual(t, "resub", message.Tags["msg-id"])
	assertStringsEqual(t, "Prime", message.Tags["msg-param-sub-plan"])
	assertStringsEqual(t, "6", message.Tags["msg-param-months"])
	assertStringsEqual(t, "Great stream -- keep it up!", message.Text)
}

func TestCanParseUsernoticeGiftSubMessage(t *testing.T) {
	testMessage := `@badges=subscriber/24,bits/25000;color=#2E8B57;display-name=TheXin1;emotes=;id=2dd9310c-1bcb-494f-929c-d0d222e245d3;login=thexin1;mod=0;msg-id=subgift;msg-param-months=1;msg-param-recipient-display-name=Fuse404;msg-param-recipient-id=36547385;msg-param-recipient-user-name=fuse404;msg-param-sub-plan-name=Channel\sSubscription\s(theattack);msg-param-sub-plan=1000;room-id=41226075;subscriber=1;system-msg=TheXin1\sgifted\sa\s$4.99\ssub\sto\sFuse404!;tmi-sent-ts=1519844687512;turbo=0;user-id=30403955;user-type= :tmi.twitch.tv USERNOTICE #theattack`

	message := parseMessage(testMessage)

	assertIntsEqual(t, int(USERNOTICE), int(message.Type))
	assertStringsEqual(t, "theattack", message.Channel)
	assertStringsEqual(t, "thexin1", message.Tags["login"])
	assertStringsEqual(t, "subgift", message.Tags["msg-id"])
	assertStringsEqual(t, "1000", message.Tags["msg-param-sub-plan"])
	assertStringsEqual(t, "1", message.Tags["msg-param-months"])
	assertStringsEqual(t, "fuse404", message.Tags["msg-param-recipient-user-name"])
	assertStringsEqual(t, "", message.Text)
}

func TestCanParseRoomstateMessage(t *testing.T) {
	testMessage := `@broadcaster-lang=<broadcaster-lang>;r9k=<r9k>;slow=<slow>;subs-only=<subs-only> :tmi.twitch.tv ROOMSTATE #nothing`

	message := parseMessage(testMessage)

	if message.Type != ROOMSTATE {
		t.Error("parsing ROOMSTATE message failed")
	}

	assertStringsEqual(t, message.Channel, "nothing")
}

func TestCanParseUserNoticeRaidMessage(t *testing.T) {
	testMessage := `@badges=turbo/1;color=#9ACD32;display-name=TestChannel;emotes=;id=3d830f12-795c-447d-af3c-ea05e40fbddb;login=testchannel;mod=0;msg-id=raid;msg-param-displayName=TestChannel;msg-param-login=testchannel;msg-param-viewerCount=15;room-id=56379257;subscriber=0;system-msg=15\sraiders\sfrom\sTestChannel\shave\sjoined\n!;tmi-sent-ts=1507246572675;tmi-sent-ts=1507246572675;turbo=1;user-id=123456;user-type= :tmi.twitch.tv USERNOTICE #othertestchannel`
	message := parseMessage(testMessage)

	if message.Type != USERNOTICE {
		t.Error("parsing USERNOTICE message failed")
	}

	assertStringsEqual(t, message.Tags["msg-id"], "raid")
	assertStringsEqual(t, message.Channel, "othertestchannel")
}
