package config

import (
	"os"
	"path"
)

func AppName() string {
	appPath := os.Args[0]
	return path.Base(appPath)
}

func Version() string {
	return "v0.9"
}
