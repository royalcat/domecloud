package entrycache

import (
	"context"
	"dmch-server/src/config"
	"dmch-server/src/domefs/entrymodel"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func (mw *EntryIndex) generateCache(ctx context.Context, realpath, virtpath string) error {
	entry, err := mw.genEntryInfo(ctx, realpath, virtpath)
	if err != nil {
		mw.log.Errorf("Error generating video info: %s", err.Error())
		return err
	}
	err = mw.genPreviews(ctx, realpath, virtpath, getTimestamps(entry.MediaInfo.VideoInfo.Duration))
	if err != nil {
		return err
	}

	return nil
}

func (mw *EntryIndex) genPreviews(ctx context.Context, realpath, virtpath string, timestamps []time.Duration) error {
	previewsDir := mw.getPreviewsDirPath(virtpath)

	os.MkdirAll(previewsDir, os.ModePerm)
	for i, timestamp := range timestamps {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		filename := fmt.Sprintf("%d.webp", i)
		output := path.Join(previewsDir, filename)
		if _, err := os.Stat(output); os.IsNotExist(err) {
			mw.ffPool.Acquire(nil, 1)
			body, err := exec.CommandContext(ctx,
				"ffmpeg",
				"-y",
				"-ss", fmt.Sprintf("%f", timestamp.Seconds()),
				"-i", realpath,
				"-vf", fmt.Sprintf("scale=%d:-1", config.Config.Media.PreviewWidth),
				"-vframes", "1",
				output,
			).CombinedOutput()
			mw.ffPool.Release(1)
			if err != nil {
				mw.log.Errorf("error creating preview: %s with output: %s", err.Error(), string(body))
				return err
			}
		}
	}
	return nil
}

func (mw *EntryIndex) genEntryInfo(ctx context.Context, realpath, virtpath string) (*entrymodel.EntryInfo, error) {
	stat, err := os.Stat(realpath)
	if err != nil {
		return nil, err
	}

	probe, err := ffprobe.ProbeURL(ctx, realpath)
	if err != nil {
		return nil, err
	}
	var videoStream *ffprobe.Stream
	for _, stream := range probe.Streams {
		if stream.CodecType == "video" {
			videoStream = stream
		}
	}
	if videoStream == nil {
		return nil, ErrVideoStreamNotFound
	}

	duration, err := strconv.ParseFloat(videoStream.Duration, 10)
	if videoStream == nil {
		return nil, fmt.Errorf("Cant parse duration with error: %s", err.Error())
	}

	info := &entrymodel.EntryInfo{
		Path:    virtpath,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
		MediaInfo: &entrymodel.MediaInfo{
			VideoInfo: &entrymodel.VideoInfo{
				Duration: time.Duration(duration * float64(time.Second)),
				Resolution: entrymodel.Resolution{
					Width:  uint64(videoStream.Width),
					Height: uint64(videoStream.Height),
				},
			},
		},
	}

	mw.index.VideoInfo.Set(ctx, *info)

	body, err := json.Marshal(info)
	infopath := mw.getInfoPath(virtpath)
	os.MkdirAll(path.Dir(infopath), os.ModePerm)
	os.WriteFile(infopath, body, os.ModePerm)

	return info, nil
}

var ErrVideoStreamNotFound = errors.New("Video stream not found")
