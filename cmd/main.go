package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"musicorginizer/bot"

	"github.com/kylelemons/go-gypsy/yaml"
)

func main() {
	config, err := yaml.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Error occured while reading config file: %s", err.Error())
		return
	}

	token, err := config.Get("token")
	if err != nil {
		log.Printf("Error parsing token value: %s\n", err.Error())
		return
	}

	maxDownloadTime, err := config.GetInt("max_download_time")
	if err != nil {
		log.Printf("Error parsing max_download_time value: %s\n", err.Error())
		return
	}

	maxVideoDuration, err := config.GetInt("max_video_duration")
	if err != nil {
		log.Printf("Error parsing max_video_duration value: %s\n", err.Error())
		return
	}

	debug, err := config.GetBool("debug")
	if err != nil {
		log.Printf("Error parsing debug value: %s\n", err.Error())
		return
	}

	botUsername, err := config.Get("bot_username")
	if err != nil {
		log.Printf("Error parsing logging value: %s\n", err.Error())
		return
	}

	b, err := bot.NewTelegramBot(token, maxDownloadTime, maxVideoDuration, botUsername)
	if err != nil {
		log.Printf("Error initializing bot: %s\n", err.Error())
		return
	}

	// graceful shutdown
	sigint := make(chan os.Signal)

	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)

	go func() {
		<-sigint

		fmt.Println("Gracefully stopping application")
		b.Stop()

		os.Exit(1)
	}()

	if err := b.Run(debug); err != nil {
		fmt.Printf("Error occured while running main event loop: %s", err.Error())
	}
}
