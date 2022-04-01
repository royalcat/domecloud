package ffmpegwrap

import (
	"dmch-server/src/config"
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

func GeneratePreviews(input string, outputDir string, timestamps []time.Duration) error {
	for i, timestamp := range timestamps {
		filename := fmt.Sprintf("%d.webp", i)
		output := path.Join(outputDir, filename)
		if _, err := os.Stat(output); os.IsNotExist(err) {
			body, err := exec.Command(
				"ffmpeg",
				"-y",
				"-ss", fmt.Sprintf("%f", timestamp.Seconds()),
				"-i", fmt.Sprintf("%s", input),
				"-vf", fmt.Sprintf("scale=%d:-1", config.Config.Media.PreviewWidth),
				"-vframes", "1",
				output,
			).CombinedOutput()
			if err != nil {
				logrus.Errorf("error creating preview: %s with output: %s", err.Error(), string(body))
				return err
			}
		}
	}
	return nil
}

func GeneratePreviewsBackgroud(input string, outputDir string, timestamps []time.Duration) error {
	for i, timestamp := range timestamps {
		filename := fmt.Sprintf("%d.webp", i)
		output := path.Join(outputDir, filename)
		if _, err := os.Stat(output); os.IsNotExist(err) {
			body, err := exec.Command(
				"ionice", "-c", "3",
				"ffmpeg",
				"-y",
				"-ss", fmt.Sprintf("%f", timestamp.Seconds()),
				"-i", fmt.Sprintf("%s", input),
				"-vf", fmt.Sprintf("scale=%d:-1", config.Config.Media.PreviewWidth),
				"-vframes", "1",
				output,
			).CombinedOutput()
			if err != nil {
				logrus.Errorf("error creating preview: %s with output: %s", err.Error(), string(body))
				return err
			}
		}
	}
	return nil
}
