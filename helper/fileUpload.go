// Package helper provides utility functions and structures to assist with
// various operations.
package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// FileUploadStruct represents the structure for handling file uploads.
// It contains information about the file's storage path and public access path.
type FileUploadStruct struct {
	StorePath           string                // The base directory where files are stored.
	File                *multipart.FileHeader // The uploaded file header.
	StoreFilePath       string                // The full path where the file is stored on the server.
	StoreFilePublicPath string                // The public path to access the file.
}

// NewFileUpload creates a new FileUploadStruct instance and saves the uploaded file
// to the specified directory.
//
// Parameters:
//   - dir: The directory where the file will be stored.
//   - file: The uploaded file header.
//
// Returns:
//   - A pointer to the FileUploadStruct instance containing file details.
//   - An error if the file upload or storage operation fails.
func NewFileUpload(dir string, file *multipart.FileHeader) (*FileUploadStruct, error) {
	fileName := uuid.New().String()
	fileExt := filepath.Ext(file.Filename)

	storageFilePath := fmt.Sprintf("files/%s/%s%s", dir, fileName, fileExt)
	storageFilePublicPath := fmt.Sprintf("static/files/%s/%s%s", dir, fileName, fileExt)

	err := os.MkdirAll(filepath.Dir(storageFilePath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(storageFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	fileUpload := &FileUploadStruct{
		File:                file,
		StorePath:           "files",
		StoreFilePath:       storageFilePath,
		StoreFilePublicPath: storageFilePublicPath,
	}

	return fileUpload, nil
}

// RemoveFile deletes the stored file from the server.
//
// Returns:
//   - An error if the file removal operation fails.
func (fileUpload *FileUploadStruct) RemoveFile() error {
	return os.Remove(fileUpload.StoreFilePath)
}
