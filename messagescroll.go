package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func messageGetTopID() int64 {
	if len(messageStore) < 1 {
		return 0
	}

	for i := 0; i < len(messageStore); i++ {
		if index := getIDfromindex(i); index != 0 {
			return index
		}
	}

	return 0
}

func getIDfromindex(i int) int64 {
	if len(messageStore[i]) < 23 {
		return 0
	}

	switch {
	case messageStore[i][1] != '[':
		return 0
	case messageStore[i][2] != '"':
		return 0
	}

	var (
		idRune string
		msg    = messageStore[i]
	)

	for i := 3; i < len(msg); i++ {
		if msg[i] == '"' {
			break
		}

		idRune += string(msg[i])
	}

	r, _ := strconv.ParseInt(idRune, 10, 64)
	return r
}

var loading bool

func loadMore() {
	if loading {
		return
	}

	loading = true
	input.SetPlaceholder("Loading more...")
	beforeID := messageGetTopID()

	msgs, err := d.ChannelMessages(ChannelID, 35, beforeID, 0, 0)
	if err != nil {
		return
	}

	if len(msgs) < 1 {
		// Drop out early if no messages
		return
	}

	var reversed []string

	for i := len(msgs) - 1; i >= 0; i-- {
		m := msgs[i]

		//wg.Add(1)
		//go func(m *discordgo.Message, i int) {
		//defer wg.Done()

		if rstore.Check(m.Author, RelationshipBlocked) && HideBlocked {
			continue
		}

		if !isRegularMessage(m) {
			continue
		}

		sentTime, err := m.Timestamp.Parse()
		if err != nil {
			sentTime = time.Now()
		}

		if i < len(msgs)-1 && msgs[i+1].Author.ID != m.Author.ID {
			username, color := us.DiscordThis(m)

			reversed = append(reversed, fmt.Sprintf(
				authorFormat,
				color, username,
				sentTime.Format(time.Stamp),
			))
		}

		reversed = append(reversed, fmt.Sprintf(
			messageFormat,
			m.ID, fmtMessage(m),
		))

		//}(m, i)
	}

	//wg.Wait()

	messageStore = append(reversed, messageStore...)

	messagesView.SetText(strings.Join(messageStore, ""))

	input.SetPlaceholder("Done.")
	app.Draw()

	messagesView.Highlight(strconv.FormatInt(beforeID, 10))
	messagesView.ScrollToHighlight()

	time.Sleep(time.Second * 5)

	input.SetPlaceholder(DefaultStatus)
	loading = false
}