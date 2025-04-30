package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

func compressImage(inputPath, outputPath string) (int64, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	newImg := resize.Resize(0, uint(img.Bounds().Dy()/2), img, resize.Lanczos3)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return 0, err
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, newImg, &jpeg.Options{Quality: 70})
	case "png":
		err = png.Encode(outFile, newImg)
	default:
		return 0, fmt.Errorf("unsupported image format: %s", format)
	}

	if err != nil {
		return 0, err
	}

	stat, err := outFile.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}
