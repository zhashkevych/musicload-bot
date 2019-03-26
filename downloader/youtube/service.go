package youtube

import (
	"context"
	"fmt"
	"musicorginizer/downloader"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rylio/ytdl"

	"os/exec"
)

const (
	expression = "^(http(s)?:\\/\\/)?((w){3}.)?(music\\.)?youtu(be|.be)?(\\.com)?\\/.+"
)

type Downloader struct {
	maxVideoDuration time.Duration

	r *regexp.Regexp
}

func NewDownloader(maxVideoDuration int64) (*Downloader, error) {
	r, err := regexp.Compile(expression)
	if err != nil {
		return nil, err
	}

	return &Downloader{
		maxVideoDuration: time.Minute * time.Duration(maxVideoDuration),
		r:                r,
	}, nil
}

func (d *Downloader) Download(ctx context.Context, url string) (string, error) {
	info, err := ytdl.GetVideoInfo(url)
	if err != nil {
		return "", errors.Wrap(err, "error getting video info")
	}

	if info.Duration > d.maxVideoDuration {
		return "", downloader.ErrDurationLimitExceeded
	}

	filename := info.Title
	strings.Replace(filename, " ", "\\ ", 0)

	// this command downloads video and extracts mp3
	cmd := exec.CommandContext(ctx, "youtube-dl", "-x", "--audio-format", "mp3", url, "-o", filename+".%(ext)s")
	data, err := cmd.CombinedOutput()
	if err != nil {
		os.Remove(filename)
		return "", errors.Wrap(err, fmt.Sprintf("error from CombinedOutput, data: %s", string(data)))
	}

	if strings.Contains(string(data), "ERROR") {
		os.Remove(filename)
		return "", errors.New(fmt.Sprintf("error downloading video with youtube-dl, output: %s", string(data)))
	}

	return filename + ".mp3", nil
}

func (d *Downloader) IsValidURL(url string) bool {
	return d.r.MatchString(url)
}
