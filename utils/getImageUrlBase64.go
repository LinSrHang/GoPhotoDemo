package utils

import (
	"encoding/base64"
	"logger"
	"os"
	"strings"
)

func GetImageUrlBase64(imagePath string) (baseImage string, ok bool) {
	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		logger.Logger.Println(err)
		return "", false
	}

	tmp := strings.Split(imagePath, ".")
	switch strings.ToLower(tmp[len(tmp)-1]) {
	case "png":
		baseImage = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imgBytes)
	case "jpg":
		baseImage = "data:image/jpg;base64," + base64.StdEncoding.EncodeToString(imgBytes)
	case "jpeg":
		baseImage = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imgBytes)
	default:
		logger.Logger.Println("unsupported file format")
		return "", false
	}

	return baseImage, true
}
