package domefs

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

	"github.com/sirupsen/logrus"
	"gopkg.in/vansante/go-ffprobe.v2"
)

func (domefs *DomeFS) getPreviewsRealPath(name string) string {
	return path.Join(domefs.cacheDir, name, "previews")
}

func (domefs *DomeFS) getInfoRealPath(name string) string {
	return path.Join(domefs.cacheDir, name, "info.json")
}

func (domefs *DomeFS) getPreviews(videoPath string, timestamps []time.Duration) []string {
	names := make([]string, 0, len(timestamps))

	os.MkdirAll(domefs.getPreviewsRealPath(videoPath), os.ModePerm)
	for i, timestamp := range timestamps {
		filename := fmt.Sprintf("%d.webp", i)
		output := path.Join(domefs.getPreviewsRealPath(videoPath), filename)
		if _, err := os.Stat(output); os.IsNotExist(err) {
			body, err := exec.Command(
				"ffmpeg",
				"-y",
				"-ss", fmt.Sprintf("%f", timestamp.Seconds()),
				"-i", fmt.Sprintf("%s", domefs.realPath(videoPath)),
				"-vf", fmt.Sprintf("scale=%d:-1", config.Config.Media.PreviewWidth),
				"-vframes", "1",
				output,
			).CombinedOutput()
			if err != nil {
				logrus.Errorf("error creating preview: %s with output: %s", err.Error(), string(body))
			}
		}

		names = append(names, filename)
	}

	return names
}

func (domefs *DomeFS) getVideoInfo(fpath string) (*media.VideoInfo, error) {
	stat, err := os.Stat(domefs.realPath(fpath))
	if err != nil {
		return nil, err
	}

	probe, err := ffprobe.ProbeURL(context.TODO(), domefs.realPath(fpath))
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
		Path:    fpath,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
		Resolution: media.Resolution{
			Width:  videoStream.Width,
			Height: videoStream.Height,
		},
		Duration: time.Duration(duration * float64(time.Second)),
	}

	body, err := json.Marshal(info)
	infopath := domefs.getInfoRealPath(fpath)
	os.MkdirAll(path.Dir(infopath), os.ModePerm)
	os.WriteFile(infopath, body, os.ModePerm)

	return info, nil
}

var ErrVideoStreamNotFound = errors.New("Video stream not found")
