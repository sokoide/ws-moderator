package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func createDirectoryIfNotExist(path string) error {
	// Create the directory if it doesn't exist
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func downloadFile(targetUrl string, base string, email string) (string, error) {
	dir := downloadFileDir(base, email)
	err := createDirectoryIfNotExist(dir)
	if err != nil {
		return "", err
	}

	filepath := filepath.Join(dir, uuid.New().String()+".png")
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(targetUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Write the body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Printf("file downloaded to %s", filepath)
	return path.Base(filepath), nil
}

func downloadFileDir(base string, email string) string {
	return filepath.Join(base, convertEmailToPath(email))
}

func convertEmailToPath(email string) string {
	// These are are allowed in email, but may cause a problem in URL
	targetRunes := `!#$%&'*+/=?^{|}~@[]();:,`
	var sb strings.Builder

	for _, char := range email {
		if strings.ContainsRune(targetRunes, char) {
			sb.WriteRune('_')
		} else {
			sb.WriteRune(char)
		}
	}
	return sb.String()
}
