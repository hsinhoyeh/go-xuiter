package goxuiter

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func SaveFile(prefix, folder, filename string, url string) error {
	fullFolderPath := filepath.Join(prefix, folder)
	if err := os.MkdirAll(fullFolderPath, os.ModePerm); err != nil {
		return err
	}
	fullFileName := filepath.Join(fullFolderPath, filename)

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	file, err := os.Create(fullFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
