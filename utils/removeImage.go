package utils

import (
	"logger"
	"os"

	"github.com/spf13/viper"
)

func RemoveImage(filenames ...string) {
	for _, filename := range filenames {
		if err := os.Remove(viper.GetString("thumbnailPrefix") + filename); err != nil {
			logger.Logger.Println(err)
		}
		if err := os.Remove(viper.GetString("originalPrefix") + filename); err != nil {
			logger.Logger.Println(err)
		}
	}
}

func RemoveOriginal(filenames ...string) {
	for _, filename := range filenames {
		if err := os.Remove(viper.GetString("originalPrefix") + filename); err != nil {
			logger.Logger.Println(err)
		}
	}
}

func RemoveThumbnail(filenames ...string) {
	for _, filename := range filenames {
		if err := os.Remove(viper.GetString("thumbnailPrefix") + filename); err != nil {
			logger.Logger.Println(err)
		}
	}
}
