package utils

import (
	"fmt"
	"image"
	"os"
)

func LoadImage(imgPath string) image.Image {

	file, err := os.Open(imgPath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer file.Close()

	// Decode the image file
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return img
}
