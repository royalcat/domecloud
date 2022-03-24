package config

import (
	"fmt"
	"mime"
	"os/exec"
	"path"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

var Config = struct {
	RootFolder  string `required:"true"`
	CacheFolder string

	Media struct {
		PreviewWidth int `default:"300"`
	}
}{}

func Load() {
	err := configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadCallback: Reload,
	}).Load(&Config, "config.yml")
	if err != nil {
		logrus.Fatalf("error loading config: %s", err.Error())
	}

	Reload(nil)
}

func Reload(_ interface{}) {
	if Config.CacheFolder == "" {
		Config.CacheFolder = path.Join(Config.RootFolder, ".cache")
	}
}

func ffmpegAvalibleExtensions() []string {
	out, err := exec.Command("bash", "-c", "ffmpeg -demuxers -hide_banner | tail -n +5 | cut -d' ' -f4 | xargs -i{} ffmpeg -hide_banner -h demuxer={} | grep 'Common extensions' | cut -d' ' -f7 | tr ',' $'\n' | tr -d '.'").Output()
	if err != nil {
		logrus.Panicf("Cant get avalible extensions from ffmpeg with error: %s", err.Error())
	}

	extensions := strings.Split(string(out), "\n")

	for _, ext := range extensions {
		mime.AddExtensionType("."+ext, fmt.Sprint("video/", ext))
	}

	return extensions
}
