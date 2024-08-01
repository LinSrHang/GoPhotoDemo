package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"logger"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Decoder func(io.Reader) (image.Image, error)

func ImageCompress(inputExt string, inputBytes []byte, sum256 string) (ok bool) {
	var decoder Decoder

	switch strings.ToLower(inputExt) {
	case "png":
		decoder = png.Decode
	case "jpg", "jpeg":
		decoder = jpeg.Decode
	default:
		logger.Logger.Println("unsupported file format")
		return false
	}

	img, err := decoder(bytes.NewReader(inputBytes))
	if err != nil {
		logger.Logger.Println(err)
		return false
	}

	bufferHigh := bytes.NewBuffer([]byte{})
	if err := jpeg.Encode(bufferHigh, img, &jpeg.Options{Quality: viper.GetInt("imageQualityHigh")}); err != nil {
		logger.Logger.Println(err)
		return false
	}
	bufferLow := bytes.NewBuffer([]byte{})
	if err := jpeg.Encode(bufferLow, img, &jpeg.Options{Quality: viper.GetInt("imageQualityLow")}); err != nil {
		logger.Logger.Println(err)
		return false
	}

	if err := os.WriteFile(viper.GetString("originalPrefix")+sum256+".jpg", bufferHigh.Bytes(), 0777); err != nil {
		logger.Logger.Println(err)
		os.RemoveAll(viper.GetString("originalPrefix") + sum256 + ".jpg")
		return false
	}
	if err := os.WriteFile(viper.GetString("thumbnailPrefix")+sum256+".jpg", bufferLow.Bytes(), 0777); err != nil {
		logger.Logger.Println(err)
		os.RemoveAll(viper.GetString("originalPrefix") + sum256 + ".jpg")
		os.RemoveAll(viper.GetString("thumbnailPrefix") + sum256 + ".jpg")
		return false
	}

	return true
}
