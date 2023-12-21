package filesystem

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"math/rand"
)

func getRandomFile(files []fs.FileInfo) (string, error) {
	randomIndex := rand.Intn(len(files))
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

func MoveFileToPostedDirectory(directory string, fileName string) (error) {
	currentLocation := createFullFilePath(directory, fileName)
	postedDirectory := directory + "/posted"
	if _, err := os.Stat(postedDirectory); os.IsNotExist(err) {
		err := os.Mkdir(postedDirectory, 0777)
		if err != nil {
			return err
		}
	}

	err := os.Rename(currentLocation, postedDirectory + "/" + fileName)
	if err != nil {
		return err
	}

	return nil
}