package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadImage(imageUrl string) string {

	// Make HTTP GET request to fetch the image
	response, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("Error fetching image:", err)
		return ""
	}
	defer response.Body.Close()

	// Create the file to which the image will be saved
	fileName := filepath.Base(imageUrl)
	filePath := fmt.Sprintf("/%s/%s", "tmp", fileName)
	outFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return ""
	}
	defer outFile.Close()

	// Copy the response body (image data) to the file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return ""
	}
	fmt.Println("Image downloaded successfully!")
	return filePath
}
