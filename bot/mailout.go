package bot

import (
	"fmt"
	"musicorginizer/downloader"
)

func (t *TelegramBot) mailoutDownloads() {
	for {
		res := <-t.downloadMsgs
		if res.Err != nil {
			if res.Err == downloader.ErrDurationLimitExceeded {
				t.send(res.ChatID, fmt.Sprintf("Can't download video longer than %d minutes", t.maxDuration))
				continue
			}

			fmt.Printf("mailoutDownloads() error: %s\n", res.Err.Error())

			go t.sendError(res.ChatID)
			continue
		}

		fmt.Printf("send result: %+v\n", res)

		go t.sendAudioFile(res.ChatID, res.Filename)
	}
}
