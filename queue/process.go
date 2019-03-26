package queue

import (
	"context"
	"time"
)

func (q *DownloadQueue) startProcess(results chan *Result) {
	for {
		msg := <-q.queue
		go q.downloadAndSend(msg, results)
	}
}

func (q *DownloadQueue) downloadAndSend(m *message, results chan *Result) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(q.maxProcessTime))
	defer cancel()

	q.doneWg.Add(1)
	defer q.doneWg.Done()

	result, err := q.handler(ctx, m.url)

	results <- &Result{
		ChatID:   m.chatID,
		Filename: result,
		Err:      err,
	}
}
