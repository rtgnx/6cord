package main

import "strconv"

func reactMessage(text []string) {
	if len(text) != 3 {
		Message("Invalid arguments! Refer to description.")
		return
	}

	messageID, err := strconv.ParseInt(text[1], 10, 64)
	if err != nil {
		Message("Failed to find the message.")
		return
	}

	message, err := d.State.Message(ChannelID, messageID)
	if err != nil {
		Message("Failed to find the message.")
		return
	}

	var (
		emoji   string
		reacted bool
	)

	regres := EmojiRegex.FindAllStringSubmatch(text[2], -1)
	if len(regres) > 0 && len(regres[0]) == 4 {
		emoji = regres[0][2] + ":" + regres[0][3]

		for _, r := range message.Reactions {
			if r.Emoji == nil {
				continue
			}

			if strconv.FormatInt(r.Emoji.ID, 10) == regres[0][3] {
				reacted = r.Me
				break
			}
		}
	} else {
		emoji = text[2]

		for _, r := range message.Reactions {
			if r.Emoji == nil {
				continue
			}

			if r.Emoji.Name == text[2] {
				reacted = r.Me
				break
			}
		}
	}

	if reacted {
		err = d.MessageReactionRemoveMe(
			ChannelID,
			message.ID,
			emoji,
		)
	} else {
		err = d.MessageReactionAdd(
			ChannelID,
			message.ID,
			emoji,
		)
	}

	if err != nil {
		Warn(err.Error())
	}

}
