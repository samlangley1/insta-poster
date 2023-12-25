package filesystem

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
)

func getRandomFile(files []fs.FileInfo) (string, error) {
	// If a non-directory file is found, set validFilesExist to true.
	validFilesExist := false
	for _, v := range files {
		if !v.IsDir() {
			validFilesExist = true
			break
		}
	}

	// If no valid files were found, return error
	if !validFilesExist {
		return "", fmt.Errorf("no valid image files found in image directory")
	}

	// Select a random file until you get one that isn't a directory
	randomIndex := rand.Intn(len(files))
	for files[randomIndex].IsDir() {
		randomIndex = rand.Intn(len(files))
	}

	img := files[randomIndex].Name()
	return img, nil
}

func GetRandomContent(imageDirectory string) (io.Reader, string, error) {
	// Get all files within the imageFilepath directory
	files, err := ioutil.ReadDir(imageDirectory)
	if err != nil {
		return nil, "", err
	}

	// Get random file from the directory
	fileName, err := getRandomFile(files)
	if err != nil {
		return nil, "", err
	}

	// Generate full file path to the randomly selected file & convert to correct type with os.Open
	fullFilePath := createFullFilePath(imageDirectory, fileName)
	postContent, err := os.Open(fullFilePath)
	if err != nil {
		return nil, "", err
	}
	return postContent, fileName, nil
}

func createFullFilePath(directory string, fileName string) string {
	fullFilePath := fmt.Sprintf("%s/%s", directory, fileName)
	return fullFilePath
}

func MoveFileToPostedDirectory(directory string, fileName string) error {
	currentLocation := createFullFilePath(directory, fileName)
	postedDirectory := directory + "/posted"

	// Create a /posted sub directory if it doesn't already exist
	if _, err := os.Stat(postedDirectory); os.IsNotExist(err) {
		err := os.Mkdir(postedDirectory, 0777)
		if err != nil {
			return err
		}
	}

	// Move image file from it's current location to the /posted sub directory
	err := os.Rename(currentLocation, postedDirectory+"/"+fileName)
	if err != nil {
		return err
	}

	return nil
}
