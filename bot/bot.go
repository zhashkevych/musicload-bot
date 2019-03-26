package bot

import (
	"fmt"
	"musicorginizer/downloader"
	"musicorginizer/downloader/youtube"
	"musicorginizer/queue"
	"os"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI

	downloadService downloader.Service
	queue           *queue.DownloadQueue
	downloadMsgs    chan *queue.Result

	username    string
	maxDuration int64
}

func NewTelegramBot(token string, maxDownloadTime, maxVideoDuration int64, username string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	downloadService, err := youtube.NewDownloader(maxVideoDuration)
	if err != nil {
		return nil, err
	}

	downloadQueue := queue.NewDownloadQueue(downloadService.Download, maxDownloadTime)

	return &TelegramBot{
		bot:             bot,
		downloadService: downloadService,
		queue:           downloadQueue,
		downloadMsgs:    make(chan *queue.Result),
		username:        username,
		maxDuration:     maxVideoDuration,
	}, nil
}

func (t *TelegramBot) Run(debug bool) error {
	t.bot.Debug = debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	t.queue.Start(t.downloadMsgs)

	go t.mailoutDownloads()

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		t.handleUpdates(update)
	}

	return nil
}

func (t *TelegramBot) Stop() {
	t.queue.Stop()
	close(t.downloadMsgs)
}

func (t *TelegramBot) send(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	t.bot.Send(msg)
}

func (t *TelegramBot) sendError(chatID int64) {
	t.send(chatID, "Error occured. Try again.")
}

func (t *TelegramBot) sendAudioFile(chatID int64, filename string) {
	path := "./" + filename

	defer os.Remove(path)

	audioCfg := tgbotapi.NewAudioUpload(chatID, path)
	audioCfg.Caption = "Downloaded via @" + t.username

	_, err := t.bot.Send(audioCfg)
	if err != nil {
		fmt.Printf("error sending message: %s", err.Error())
		t.sendError(chatID)
	}
}
