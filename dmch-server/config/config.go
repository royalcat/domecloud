package config

import (
	"path"

	"github.com/jinzhu/configor"
)

var Config = struct {
	RootFolder string `required:"true"`

	Media struct {
		PreviewFolder string
		PreviewWidth  int `default:"300"`
	}
}{}

func Load() {
	configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadCallback: Reload,
	}).Load(&Config, "config.yml")
}

func Reload(_ interface{}) {
	if Config.Media.PreviewFolder == "" {
		Config.Media.PreviewFolder = path.Join(Config.RootFolder, ".previews")
	}

}
