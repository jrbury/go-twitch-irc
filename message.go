package twitch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MessageType different message types possible to receive via IRC
type MessageType int

const (
	// UNSET default type
	UNSET MessageType = -1
	// WHISPER private messages
	WHISPER MessageType = 0
	// PRIVMSG standard chat message
	PRIVMSG MessageType = 1
	// CLEARCHAT timeout messages
	CLEARCHAT MessageType = 2
	// ROOMSTATE changes like sub mode
	ROOMSTATE MessageType = 3
	// USERNOTICE messages like subs, resubs, raids, etc
	USERNOTICE MessageType = 4
	// USERSTATE messages
	USERSTATE MessageType = 5
	// NOTICE messages like sub mode, host on
	NOTICE MessageType = 6
)

type message struct {
	Type        MessageType
	Time        time.Time
	Channel     string
	ChannelID   string
	UserID      string
	Username    string
	DisplayName string
	UserType    string
	Color       string
	Action      bool
	Badges      map[string]int
	Emotes      []*Emote
	Tags        map[string]string
	Text        string
	Raw         string
}

// Emote twitch emotes
type Emote struct {
	Name  string
	ID    string
	Count int
}

func parseMessage(line string) *message {
	if !strings.HasPrefix(line, "@") {
		return &message{
			Text: line,
			Raw:  line,
			Type: UNSET,
		}
	}

	tags, middle, text := splitLine(line)

	action := false
	if strings.HasPrefix(text, "\u0001ACTION ") {
		action = true
		text = text[8 : len(text)-1]
	}
	msg := &message{
		Text:   text,
		Tags:   map[string]string{},
		Action: action,
		Type:   UNSET,
	}
	msg.Username, msg.Type, msg.Channel = parseMiddle(middle)
	parseTags(msg, tags[1:])
	if msg.Type == CLEARCHAT {
		targetUser := msg.Text
		msg.Username = targetUser

		msg.Text = fmt.Sprintf("%s was timed out for %s: %s", targetUser, msg.Tags["ban-duration"], msg.Tags["ban-reason"])
	}
	msg.Raw = line
	return msg
}

func splitLine(line string) (string, string, string) {
	spl := strings.SplitN(line, " :", 3)

	switch len(spl) {
	case 2:
		return spl[0], spl[1], ""
	case 3:
		return spl[0], spl[1], spl[2]
	default:
		return "", "", line // we don't have what we expect, return the line as text
	}
}

// The main reason for using regex vs splitting is to reduce the edge cases
// that needed to be accounted for
func parseMiddle(middle string) (string, MessageType, string) {
	var username string
	var msgType MessageType
	var channel string

	userRe := regexp.MustCompile(`@([a-z0-9_]+)\.tmi\.twitch\.tv`)
	userMatch := userRe.FindStringSubmatch(middle)
	if len(userMatch) > 1 {
		username = userMatch[1]
	}

	typeRe := regexp.MustCompile(`\s([A-Z]+)\s`)
	typeMatch := typeRe.FindStringSubmatch(middle)
	if len(typeMatch) > 1 {
		switch typeMatch[1] {
		case "PRIVMSG":
			msgType = PRIVMSG
		case "WHISPER":
			msgType = WHISPER
		case "CLEARCHAT":
			msgType = CLEARCHAT
		case "NOTICE":
			msgType = NOTICE
		case "ROOMSTATE":
			msgType = ROOMSTATE
		case "USERSTATE":
			msgType = USERSTATE
		case "USERNOTICE":
			msgType = USERNOTICE
		default:
			msgType = UNSET
		}
	}

	channelRe := regexp.MustCompile(`#([a-zA-Z0-9_]+)$`)
	channelMatch := channelRe.FindStringSubmatch(middle)
	if len(channelMatch) > 1 {
		channel = channelMatch[1]
	}

	return username, msgType, channel
}

func parseTags(msg *message, tagsRaw string) {
	tags := strings.Split(tagsRaw, ";")

	for _, tag := range tags {
		spl := strings.SplitN(tag, "=", 2)
		if len(spl) != 2 {
			continue
		}

		value := strings.Replace(spl[1], "\\:", ";", -1)
		value = strings.Replace(value, "\\s", " ", -1)
		value = strings.Replace(value, "\\\\", "\\", -1)
		switch spl[0] {
		case "badges":
			msg.Badges = parseBadges(value)
		case "color":
			msg.Color = value
		case "display-name":
			msg.DisplayName = value
		case "emotes":
			msg.Emotes = parseTwitchEmotes(value, msg.Text)
		case "user-type":
			msg.UserType = value
		case "tmi-sent-ts":
			i, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				msg.Time = time.Unix(0, int64(i*1e6))
			}
		case "room-id":
			msg.ChannelID = value
		case "target-user-id":
			msg.UserID = value
		case "user-id":
			msg.UserID = value
		}
		msg.Tags[spl[0]] = value
	}
}

func parseBadges(badges string) map[string]int {
	m := map[string]int{}
	spl := strings.Split(badges, ",")
	for _, badge := range spl {
		s := strings.SplitN(badge, "/", 2)
		if len(s) < 2 {
			continue
		}
		n, _ := strconv.Atoi(s[1])
		m[s[0]] = n
	}
	return m
}

func parseTwitchEmotes(emoteTag, text string) []*Emote {
	emotes := []*Emote{}

	if emoteTag == "" {
		return emotes
	}

	runes := []rune(text)

	emoteSlice := strings.Split(emoteTag, "/")
	for i := range emoteSlice {
		spl := strings.Split(emoteSlice[i], ":")
		pos := strings.Split(spl[1], ",")
		sp := strings.Split(pos[0], "-")
		start, _ := strconv.Atoi(sp[0])
		end, _ := strconv.Atoi(sp[1])
		id := spl[0]
		e := &Emote{
			ID:    id,
			Count: strings.Count(emoteSlice[i], "-"),
			Name:  string(runes[start : end+1]),
		}

		emotes = append(emotes, e)
	}
	return emotes
}

func parseJoinPart(text string) (string, string) {
	username := strings.Split(text, "!")
	channel := strings.Split(username[1], "#")
	return strings.Trim(channel[1], " "), strings.Trim(username[0], " :")
}

func parseNames(text string) (string, []string) {
	lines := strings.Split(text, ":")
	channelDirty := strings.Split(lines[1], "#")
	channel := strings.Trim(channelDirty[1], " ")
	users := strings.Split(lines[2], " ")

	return channel, users
}
