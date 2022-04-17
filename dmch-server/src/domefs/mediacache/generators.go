package mediacache

import (
	"context"
	"dmch-server/src/config"
	"dmch-server/src/domefs/media"
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

func (mw *MediaCache) generateCache(ctx context.Context, virtpath string) error {
	realpath := mw.dfs.RealPath(virtpath)
	info, err := mw.genVideoInfo(ctx, virtpath, realpath)
	if err != nil {
		mw.log.Errorf("Eror generating video info: %w", err)
		return err
	}
	mw.genPreviews(ctx, virtpath, realpath, getTimestamps(info.Duration))

	return nil
}

func (mw *MediaCache) genPreviews(ctx context.Context, virtpath, realpath string, timestamps []time.Duration) error {
	previewsDir := mw.getPreviewsDirPath(virtpath)

	os.MkdirAll(previewsDir, os.ModePerm)
	for i, timestamp := range timestamps {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		filename := fmt.Sprintf("%d.webp", i)
		output := path.Join(previewsDir, filename)
		if _, err := os.Stat(output); os.IsNotExist(err) {
			mw.fflock.Lock()
			body, err := exec.CommandContext(ctx,
				"ffmpeg",
				"-y",
				"-ss", fmt.Sprintf("%f", timestamp.Seconds()),
				"-i", realpath,
				"-vf", fmt.Sprintf("scale=%d:-1", config.Config.Media.PreviewWidth),
				"-vframes", "1",
				output,
			).CombinedOutput()
			mw.fflock.Unlock()
			if err != nil {
				mw.log.Errorf("error creating preview: %s with output: %s", err.Error(), string(body))
				return err
			}
		}
	}
	return nil
}

func (mw *MediaCache) genVideoInfo(ctx context.Context, virtpath, realpath string) (*media.VideoInfo, error) {
	stat, err := os.Stat(realpath)
	if err != nil {
		return nil, err
	}

	mw.fflock.Lock()
	probe, err := ffprobe.ProbeURL(ctx, realpath)
	mw.fflock.Unlock()
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

	info := &media.VideoInfo{
		Path:    virtpath,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
		Resolution: media.Resolution{
			Width:  videoStream.Width,
			Height: videoStream.Height,
		},
		Duration: time.Duration(duration * float64(time.Second)),
	}

	body, err := json.Marshal(info)
	infopath := mw.getInfoPath(virtpath)
	os.MkdirAll(path.Dir(infopath), os.ModePerm)
	os.WriteFile(infopath, body, os.ModePerm)

	return info, nil
}

var ErrVideoStreamNotFound = errors.New("Video stream not found")
