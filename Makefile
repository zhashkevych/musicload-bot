build:
	GOOS=linux go build -o app ./cmd/
	docker build -t downloader/telegrambot .
	rm -f app

run:
	docker run --rm downloader/telegrambot 