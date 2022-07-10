package previewbuilder

import (
	"context"
	"dmch-server/src/config"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type PreviewBuilder struct {
	ffPool semaphore.Weighted

	log *logrus.Entry
}

func (b *PreviewBuilder) GenVideoPreviews(ctx context.Context, in, previewsDir string) error {
	dur, err := b.getDuration(ctx, in)
	if err != nil {
		return err
	}

	os.MkdirAll(previewsDir, os.ModePerm)
	for i, timestamp := range b.getTimestamps(dur) {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		filename := fmt.Sprintf("%d.webp", i)
		output := path.Join(previewsDir, filename)
		if _, err := os.Stat(output); os.IsNotExist(err) {
			err := b.ffPool.Acquire(ctx, 1)
			if err != nil {
				return err
			}
			body, err := exec.CommandContext(ctx,
				"ffmpeg",
				"-y",
				"-ss", fmt.Sprintf("%f", timestamp.Seconds()),
				"-i", in,
				"-vf", fmt.Sprintf("scale=%d:-1", config.Config.Media.PreviewWidth),
				"-vframes", "1",
				output,
			).CombinedOutput()
			b.ffPool.Release(1)
			if err != nil {
				b.log.Errorf("error creating preview: %s with output: %s", err.Error(), string(body))
				return err
			}
		}
	}
	return nil
}

// MAYBE тут может быть умная функция
func (b *PreviewBuilder) getTimestamps(duration time.Duration) []time.Duration {
	stamps := make([]time.Duration, 0, 1)
	stampDuration := duration / 10
	for i := 0; i < int(duration); i += int(stampDuration) {
		stamps = append(stamps, time.Duration(i))
	}
	return stamps
}

func (b *PreviewBuilder) getDuration(ctx context.Context, in string) (time.Duration, error) {
	probe, err := ffprobe.ProbeURL(ctx, in)
	if err != nil {
		return 0, err
	}
	return probe.Format.Duration(), nil

}
