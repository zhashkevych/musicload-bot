FROM alpine

RUN apk update && apk add curl && apk add python && apk add ffmpeg

RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

ADD app .
ADD config.yaml .

RUN chmod +x ./app

ENTRYPOINT ["/app"]