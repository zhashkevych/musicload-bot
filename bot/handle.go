package bot

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func (t *TelegramBot) handleUpdates(update tgbotapi.Update) {
	if m := update.Message; m != nil {
		if m.IsCommand() && m.Command() == "start" {
			t.send(m.Chat.ID, "Hi there! Send me a link to a video you want extract music from.")
			return
		}

		if t.downloadService.IsValidURL(m.Text) {
			t.queue.Enqueue(m)
			return
		}

		t.send(m.Chat.ID, "Invalid message text. I'm waiting for youtube link.")
	}
}
