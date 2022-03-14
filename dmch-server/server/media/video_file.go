package media

import (
	"context"
	"dmch-server/cfs"
	"dmch-server/config"
	"fmt"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type Resolution struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type VideoInfo struct {
	Path       string        `json:"name"`
	Duration   time.Duration `json:"duration"`
	Resolution Resolution    `json:"resolution"`
}

func GenerateVideoInfo(ctx context.Context, path string) (*VideoInfo, error) {
	fullpath := cfs.Cfs.RealPath(path)
	probe, err := ffprobe.ProbeURL(ctx, fullpath)
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

	info := &VideoInfo{
		Path: path,
		Resolution: Resolution{
			Width:  videoStream.Width,
			Height: videoStream.Height,
		},
		Duration: time.Duration(duration * float64(time.Second)),
	}

	return info, nil
}

func (f VideoInfo) Name() string {
	return path.Base(f.Path)
}

func (f VideoInfo) Previews() []string {
	names := []string{}
	folder := f.previewFolder()
	dirs, err := cfs.Cfs.ReadDir(folder)
	if err != nil {
		logrus.Errorf("reading dir %s failed with error: %v", f.Path, err)
		return names
	}
	for _, dir := range dirs {
		names = append(names, path.Join(folder, dir.Name()))
	}

	return names
}

func (f VideoInfo) previewFolder() string {
	return path.Join(config.Config.Media.PreviewFolder, f.Path)
}

func (f VideoInfo) generatePreviews(ctx context.Context, videoPath string, timestamps []time.Duration) {
	for i, timestamp := range timestamps {
		exec.CommandContext(
			ctx, "ffmpeg",
			fmt.Sprintf("-ss %f", timestamp.Seconds()),
			fmt.Sprintf("-i %s", path.Join(config.Config.RootFolder, videoPath)),
			fmt.Sprintf("-vf \"scale=%d:-1\"", config.Config.Media.PreviewWidth),
			"-vframes 1",
			path.Join(f.previewFolder(), fmt.Sprintf("%d.webp", i)),
		)
	}
}
